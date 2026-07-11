package translate

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func sv39(cpu *state.CPU, virtAddr int, physAddr *int, access int) {
	// https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#3-1-1-6-4-memory-privilege-in-mstatus-register
	epriv := cpu.Priv
	if cpu.CSR.Mstatus>>csr.MstatusMPRV&1 == 1 && access != state.AccessExecute {
		epriv = cpu.CSR.Mstatus >> csr.MstatusMPP & 3
	}

	// https://docs.riscv.org/reference/isa/v20260120/priv/supervisor.html#norm:satp_op_active
	// https://docs.riscv.org/reference/isa/v20260120/priv/supervisor.html#norm:satp-mode
	if epriv == state.PrivM || cpu.CSR.Satp == 0 {
		*physAddr = virtAddr
		return
	}

	// https://docs.riscv.org/reference/isa/v20260120/priv/supervisor.html#addressing-and-memory-protection
	if upper := virtAddr >> 38; upper != 0 && upper != -1 {
		trap.Enter(cpu, trap.PageFault+access, virtAddr)
		return
	}

	pte, shift := cpu.TLB.Lookup(virtAddr)
	if pte == 0 {
		if loadPTEsv39(cpu, virtAddr, &pte, &shift); trap.IsEntered(cpu) {
			return
		}

		if pte != 0 {
			cpu.TLB.Append(virtAddr, shift, pte)
		}
	}

	sum, mxr := cpu.CSR.Mstatus>>csr.MstatusSUM&1, cpu.CSR.Mstatus>>csr.MstatusMXR&1

	if pte == 0 ||
		epriv == state.PrivU && pte>>PteU&1 == 0 ||
		epriv == state.PrivS && pte>>PteU&1 == 1 && !(sum == 1 && access != state.AccessExecute) ||
		access == state.AccessExecute && pte>>PteX&1 == 0 ||
		access == state.AccessRead && pte>>PteR&1 == 0 && !(mxr == 1 && pte>>PteX&1 == 1) ||
		access == state.AccessWrite && !(pte>>PteW&1 == 1 && pte>>PteD&1 == 1) ||
		pte>>PteA&1 == 0 {

		trap.Enter(cpu, trap.PageFault+access, virtAddr)
		return
	}

	// https://docs.riscv.org/reference/isa/v20260120/priv/supervisor.html#sv39pte
	*physAddr = pte>>10&^(-1<<44)<<12 | virtAddr&^(-1<<shift)
}

// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#sv32algorithm
func loadPTEsv39(cpu *state.CPU, virtAddr int, targetPTE, shift *int) {
	*targetPTE = 0
	var pte int

	pteAddr := cpu.CSR.Satp&^(-1<<44)<<12 | (virtAddr>>30&511)<<3

	var ok bool
	if pte, ok = cpu.Bus.Read(pteAddr, 8); !ok {
		trap.Enter(cpu, trap.LoadAccessFault, virtAddr)
		return
	}

	isLeaf := pte>>PteR&1 == 1 || pte>>PteX&1 == 1

	if pte>>PteV&1 == 0 || // valid bit not set
		pte>>PteR&1 == 0 && pte>>PteW&1 == 1 || // reserved
		isLeaf && pte>>10&^(-1<<18) != 0 { // misaligned gigapage
		return
	}

	if isLeaf {
		*targetPTE = pte
		*shift = 30
		return
	}

	pteAddr = pte>>10&^(-1<<44)<<12 | (virtAddr>>21&511)<<3
	if pte, ok = cpu.Bus.Read(pteAddr, 8); !ok {
		trap.Enter(cpu, trap.LoadAccessFault, virtAddr)
		return
	}

	isLeaf = pte>>PteR&1 == 1 || pte>>PteX&1 == 1

	if pte>>PteV&1 == 0 || // valid bit not set
		pte>>PteR&1 == 0 && pte>>PteW&1 == 1 || // reserved
		isLeaf && pte>>10&^(-1<<9) != 0 { // misaligned megapage
		return
	}

	if isLeaf {
		*targetPTE = pte
		*shift = 21
		return
	}

	pteAddr = pte>>10&^(-1<<44)<<12 | (virtAddr>>12&511)<<3
	if pte, ok = cpu.Bus.Read(pteAddr, 8); !ok {
		trap.Enter(cpu, trap.LoadAccessFault, virtAddr)
		return
	}

	if pte>>PteV&1 == 0 ||
		pte>>PteR&1 == 0 && !(pte>>PteW&1 == 0 && pte>>PteX&1 == 1) {
		return
	}

	*targetPTE = pte
	*shift = 12
}
