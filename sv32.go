package rv

func (cpu *CPU) translateSv32(virtAddr int, physAddr *int, access int) {
	// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_memory_privilege_in_mstatus_register
	epriv := cpu.Priv
	if bit(cpu.CSR.Mstatus, MstatusMPRV) == 1 && access != AccessExecute {
		epriv = bits(cpu.CSR.Mstatus, MstatusMPP, 2)
	}

	// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#satp-mode
	if bit(cpu.CSR.Satp, SatpMODE32) == 0 || epriv == PrivM {
		*physAddr = virtAddr
		return
	}

	pte, shift := cpu.TLB.lookup(virtAddr)
	if pte == 0 {
		if cpu.loadPTEsv32(virtAddr, &pte, &shift); cpu.isTrapped() {
			return
		}

		if pte != 0 {
			cpu.TLB.append(virtAddr, shift, pte)
		}
	}

	sum, mxr := bit(cpu.CSR.Mstatus, MstatusSUM), bit(cpu.CSR.Mstatus, MstatusMXR)

	if pte == 0 ||
		epriv == PrivU && bit(pte, PteU) == 0 ||
		epriv == PrivS && bit(pte, PteU) == 1 && !(sum == 1 && access != AccessExecute) ||
		access == AccessExecute && bit(pte, PteX) == 0 ||
		access == AccessRead && bit(pte, PteR) == 0 && !(mxr == 1 && bit(pte, PteX) == 1) ||
		access == AccessWrite && !(bit(pte, PteW) == 1 && bit(pte, PteD) == 1) ||
		bit(pte, PteA) == 0 {

		cpu.trapEnter(ExceptionPageFault+access, virtAddr)
		return
	}

	*physAddr = cpu.xint(bits(pte, 10, 20)<<12 | bits(virtAddr, 0, shift))
}

// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#sv32algorithm
func (cpu *CPU) loadPTEsv32(virtAddr int, targetPTE, shift *int) {
	*targetPTE = 0
	var pte int

	pteAddr := cpu.xint(bits(cpu.CSR.Satp, 0, 20)<<12 | bits(virtAddr, 22, 10)<<2)
	if !cpu.Bus.Read(pteAddr, &pte, 4) {
		cpu.trapEnter(ExceptionLoadAccessFault, virtAddr)
		return
	}

	isLeaf := bit(pte, PteR) == 1 || bit(pte, PteX) == 1

	if bit(pte, PteV) == 0 || // valid bit not set
		bit(pte, PteR) == 0 && bit(pte, PteW) == 1 || // reserved
		isLeaf && bits(pte, 10, 10) != 0 { // misaligned superpage
		return
	}

	*shift = 22

	if !isLeaf {
		pteAddr = cpu.xint(bits(pte, 10, 20)<<12 | bits(virtAddr, 12, 10)<<2)
		if !cpu.Bus.Read(pteAddr, &pte, 4) {
			cpu.trapEnter(ExceptionLoadAccessFault, virtAddr)
			return
		}

		if bit(pte, PteV) == 0 || bit(pte, PteR) == 0 && !(bit(pte, PteW) == 0 && bit(pte, PteX) == 1) {
			return
		}

		*shift = 12
	}

	*targetPTE = pte
}
