package rv

import (
	"github.com/temnok/rv/bi"
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/trap"
)

func (cpu *CPU) translateSv39(virtAddr int, physAddr *int, access int) {
	if upper := virtAddr >> 38; upper != 0 && upper != -1 {
		trap.Enter(&cpu.State, ExceptionPageFault+access, virtAddr)
		return
	}

	// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_memory_privilege_in_mstatus_register
	epriv := cpu.Priv
	if bi.T(cpu.CSR.Mstatus, csr.MstatusMPRV) == 1 && access != AccessExecute {
		epriv = bi.Ts(cpu.CSR.Mstatus, csr.MstatusMPP, 2)
	}

	// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#satp-mode
	if bi.Ts(cpu.CSR.Satp, csr.SatpMODE64, 4) == 0 || epriv == PrivM {
		*physAddr = virtAddr
		return
	}

	pte, shift := cpu.TLB.lookup(virtAddr)
	if pte == 0 {
		if cpu.loadPTEsv39(virtAddr, &pte, &shift); trap.IsEntered(&cpu.State) {
			return
		}

		if pte != 0 {
			cpu.TLB.append(virtAddr, shift, pte)
		}
	}

	sum, mxr := bi.T(cpu.CSR.Mstatus, csr.MstatusSUM), bi.T(cpu.CSR.Mstatus, csr.MstatusMXR)

	if pte == 0 ||
		epriv == PrivU && bi.T(pte, PteU) == 0 ||
		epriv == PrivS && bi.T(pte, PteU) == 1 && !(sum == 1 && access != AccessExecute) ||
		access == AccessExecute && bi.T(pte, PteX) == 0 ||
		access == AccessRead && bi.T(pte, PteR) == 0 && !(mxr == 1 && bi.T(pte, PteX) == 1) ||
		access == AccessWrite && !(bi.T(pte, PteW) == 1 && bi.T(pte, PteD) == 1) ||
		bi.T(pte, PteA) == 0 {

		trap.Enter(&cpu.State, ExceptionPageFault+access, virtAddr)
		return
	}

	*physAddr = bi.Ts(pte, 10, 44)<<12 | bi.Ts(virtAddr, 0, shift)
}

// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#sv32algorithm
func (cpu *CPU) loadPTEsv39(virtAddr int, targetPTE, shift *int) {
	*targetPTE = 0
	var pte int

	pteAddr := bi.Ts(cpu.CSR.Satp, 0, 44)<<12 | bi.Ts(virtAddr, 30, 9)<<3

	//panic(fmt.Sprintf("*** oops: virtAddr:%x, pteAddr:%x, pte:%x", uint(virtAddr), uint(pteAddr), uint(pte)))

	if !cpu.Bus.Read(pteAddr, &pte, 8) {
		trap.Enter(&cpu.State, ExceptionLoadAccessFault, virtAddr)
		return
	}

	isLeaf := bi.T(pte, PteR) == 1 || bi.T(pte, PteX) == 1

	if bi.T(pte, PteV) == 0 || // valid bit not set
		bi.T(pte, PteR) == 0 && bi.T(pte, PteW) == 1 || // reserved
		isLeaf && bi.Ts(pte, 10, 18) != 0 { // misaligned gigapage
		return
	}

	if isLeaf {
		*targetPTE = pte
		*shift = 30
		return
	}

	pteAddr = bi.Ts(pte, 10, 44)<<12 | bi.Ts(virtAddr, 21, 9)<<3
	if !cpu.Bus.Read(pteAddr, &pte, 8) {
		trap.Enter(&cpu.State, ExceptionLoadAccessFault, virtAddr)
		return
	}

	isLeaf = bi.T(pte, PteR) == 1 || bi.T(pte, PteX) == 1

	if bi.T(pte, PteV) == 0 || // valid bit not set
		bi.T(pte, PteR) == 0 && bi.T(pte, PteW) == 1 || // reserved
		isLeaf && bi.Ts(pte, 10, 9) != 0 { // misaligned megapage
		return
	}

	if isLeaf {
		*targetPTE = pte
		*shift = 21
		return
	}

	pteAddr = bi.Ts(pte, 10, 44)<<12 | bi.Ts(virtAddr, 12, 9)<<3
	if !cpu.Bus.Read(pteAddr, &pte, 8) {
		trap.Enter(&cpu.State, ExceptionLoadAccessFault, virtAddr)
		return
	}

	if bi.T(pte, PteV) == 0 ||
		bi.T(pte, PteR) == 0 && !(bi.T(pte, PteW) == 0 && bi.T(pte, PteX) == 1) {
		return
	}

	*targetPTE = pte
	*shift = 12
}
