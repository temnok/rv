package rv

const (
	// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#translation
	PteV = 0
	PteR = 1
	PteW = 2
	PteX = 3
	PteU = 4
	//PteG = 5
	PteA = 6
	PteD = 7
)

func (cpu *CPU) translateSv32(virtAddr int32, physAddr *int32, access int32) {
	// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_memory_privilege_in_mstatus_register
	epriv := cpu.priv
	if bit(cpu.csr.mstatus, mstatusMPRV) != 0 && access != AccessExecute {
		epriv = bits(cpu.csr.mstatus, mstatusMPP, 2)
	}

	// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#satp-mode
	if bit(cpu.csr.satp, satpMODE) == 0 || epriv == PrivM {
		*physAddr = int32(uint32(virtAddr) >> 2)
		return
	}

	pte, shift := cpu.tlb.lookup(virtAddr)
	if pte == 0 {
		if cpu.loadPTE(virtAddr, &pte, &shift); cpu.isTrapped {
			return
		}

		if pte != 0 {
			cpu.tlb.append(virtAddr, shift, pte)
		}
	}

	sum, mxr := bit(cpu.csr.mstatus, mstatusSUM), bit(cpu.csr.mstatus, mstatusMXR)

	if pte == 0 ||
		epriv == PrivU && bit(pte, PteU) == 0 ||
		epriv == PrivS && bit(pte, PteU) == 1 && !(sum == 1 && access != AccessExecute) ||
		access == AccessExecute && bit(pte, PteX) == 0 ||
		access == AccessRead && bit(pte, PteR) == 0 && !(mxr == 1 && bit(pte, PteX) == 1) ||
		access == AccessWrite && !(bit(pte, PteW) == 1 && bit(pte, PteD) == 1) ||
		bit(pte, PteA) == 0 {

		cpu.trapWithTval(ExceptionInstructionPageFault+access, virtAddr)
		return
	}

	*physAddr = pte&^0x3FF | bits(virtAddr, 2, shift-2)
}

// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#sv32algorithm
func (cpu *CPU) loadPTE(virtAddr int32, targetPTE, shift *int32) {
	*targetPTE = 0
	var pte int32

	pteAddr := bits(cpu.csr.satp, 0, 22)<<10 | bits(virtAddr, 22, 10)
	if !cpu.bus.read(pteAddr, &pte) {
		cpu.trapWithTval(ExceptionLoadAccessFault, virtAddr)
		return
	}

	isLeaf := bit(pte, PteR) != 0 || bit(pte, PteX) != 0

	if bit(pte, PteV) == 0 || // valid bit not set
		bit(pte, PteR) == 0 && bit(pte, PteW) == 1 || // reserved
		isLeaf && bits(pte, 10, 10) != 0 { // misaligned superpage
		return
	}

	*shift = 22

	if !isLeaf {
		pteAddr = bits(pte, 10, 22)<<10 | bits(virtAddr, 12, 10)
		if !cpu.bus.read(pteAddr, &pte) {
			cpu.trapWithTval(ExceptionLoadAccessFault, virtAddr)
			return
		}

		if bit(pte, PteV) == 0 || bit(pte, PteR) == 0 && !(bit(pte, PteW) == 0 && bit(pte, PteX) == 1) {
			return
		}

		*shift = 12
	}

	*targetPTE = pte
}
