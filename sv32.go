package rv

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

	// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#sv32algorithm
	var pte int32
	pteAddr := bits(cpu.csr.satp, 0, 22)<<10 | bits(virtAddr, 22, 10)
	if !cpu.bus.read(pteAddr, &pte) {
		cpu.trapWithTval(ExceptionLoadAccessFault, virtAddr)
		return
	}

	lowAddrBits := bits(virtAddr, 2, 20)
	pteR, pteW, pteX := bit(pte, PteR), bit(pte, PteW), bit(pte, PteX)
	isLeaf := pteR != 0 || pteX != 0

	pteInvalid := bit(pte, PteV) == 0 || // valid bit not set
		pteR == 0 && pteW == 1 || // reserved
		isLeaf && bits(pte, 10, 10) != 0 // misaligned superpage

	if !pteInvalid && !isLeaf {
		pteAddr = bits(pte, 10, 22)<<10 | bits(virtAddr, 12, 10)
		if !cpu.bus.read(pteAddr, &pte) {
			cpu.trapWithTval(ExceptionLoadAccessFault, virtAddr)
			return
		}

		lowAddrBits = bits(virtAddr, 2, 10)

		pteR, pteW, pteX = bit(pte, PteR), bit(pte, PteW), bit(pte, PteX)
		pteInvalid = bit(pte, PteV) == 0 ||
			pteR == 0 && !(pteW == 0 && pteX == 1)
	}

	pteU, pteD, pteA := bit(pte, PteU), bit(pte, PteD), bit(pte, PteA)
	sum, mxr := bit(cpu.csr.mstatus, mstatusSUM), bit(cpu.csr.mstatus, mstatusMXR)

	if pteInvalid ||
		epriv == PrivU && pteU == 0 ||
		epriv == PrivS && pteU == 1 && !(sum == 1 && access != AccessExecute) ||
		access == AccessExecute && pteX == 0 ||
		access == AccessRead && pteR == 0 && !(mxr == 1 && pteX == 1) ||
		access == AccessWrite && (pteW == 0 || pteD == 0) ||
		pteA == 0 {

		cpu.trapWithTval(ExceptionInstructionPageFault+access, virtAddr)
		return
	}

	*physAddr = pte&^0x3FF | lowAddrBits
}
