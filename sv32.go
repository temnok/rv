package rv

import (
	"github.com/temnok/rv/bi"
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/trap"
)

func (cpu *CPU) translateSv32(virtAddr int, physAddr *int, access int) {
	// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_memory_privilege_in_mstatus_register
	epriv := cpu.Priv
	if bi.T(cpu.CSR.Mstatus, csr.MstatusMPRV) == 1 && access != AccessExecute {
		epriv = bi.Ts(cpu.CSR.Mstatus, csr.MstatusMPP, 2)
	}

	// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#satp-mode
	if bi.T(cpu.CSR.Satp, csr.SatpMODE32) == 0 || epriv == PrivM {
		*physAddr = virtAddr
		return
	}

	pte, shift := cpu.TLB.lookup(virtAddr)
	if pte == 0 {
		if cpu.loadPTEsv32(virtAddr, &pte, &shift); trap.IsEntered(&cpu.State) {
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

	*physAddr = cpu.Xint(bi.Ts(pte, 10, 20)<<12 | bi.Ts(virtAddr, 0, shift))
}

// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#sv32algorithm
func (cpu *CPU) loadPTEsv32(virtAddr int, targetPTE, shift *int) {
	*targetPTE = 0
	var pte int

	pteAddr := cpu.Xint(bi.Ts(cpu.CSR.Satp, 0, 20)<<12 | bi.Ts(virtAddr, 22, 10)<<2)
	if !cpu.Bus.Read(pteAddr, &pte, 4) {
		trap.Enter(&cpu.State, ExceptionLoadAccessFault, virtAddr)
		return
	}

	isLeaf := bi.T(pte, PteR) == 1 || bi.T(pte, PteX) == 1

	if bi.T(pte, PteV) == 0 || // valid bit not set
		bi.T(pte, PteR) == 0 && bi.T(pte, PteW) == 1 || // reserved
		isLeaf && bi.Ts(pte, 10, 10) != 0 { // misaligned superpage
		return
	}

	*shift = 22

	if !isLeaf {
		pteAddr = cpu.Xint(bi.Ts(pte, 10, 20)<<12 | bi.Ts(virtAddr, 12, 10)<<2)
		if !cpu.Bus.Read(pteAddr, &pte, 4) {
			trap.Enter(&cpu.State, ExceptionLoadAccessFault, virtAddr)
			return
		}

		if bi.T(pte, PteV) == 0 || bi.T(pte, PteR) == 0 && !(bi.T(pte, PteW) == 0 && bi.T(pte, PteX) == 1) {
			return
		}

		*shift = 12
	}

	*targetPTE = pte
}
