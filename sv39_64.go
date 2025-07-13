package rv

func (cpu *CPU) translateSv39(virtAddr int, physAddr *int, access int) {
	if upper := virtAddr >> 38; upper != 0 && upper != -1 {
		cpu.trapWithTval(ExceptionPageFault+access, virtAddr)
		return
	}

	// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_memory_privilege_in_mstatus_register
	epriv := cpu.priv
	if bit(cpu.csr.mstatus, mstatusMPRV) == 1 && access != AccessExecute {
		epriv = bits(cpu.csr.mstatus, mstatusMPP, 2)
	}

	// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#satp-mode
	if bits(cpu.csr.satp, satpMODE64, 4) == 0 || epriv == PrivM {
		*physAddr = virtAddr
		return
	}

	pte, shift := cpu.tlb.lookup(virtAddr)
	if pte == 0 {
		if cpu.loadPTEsv39(virtAddr, &pte, &shift); cpu.isTrapped() {
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

		cpu.trapWithTval(ExceptionPageFault+access, virtAddr)
		return
	}

	*physAddr = bits(pte, 10, 44)<<12 | bits(virtAddr, 0, shift)
}

// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#sv32algorithm
func (cpu *CPU) loadPTEsv39(virtAddr int, targetPTE, shift *int) {
	*targetPTE = 0
	var pte int

	pteAddr := bits(cpu.csr.satp, 0, 44)<<12 | bits(virtAddr, 30, 9)<<3

	//panic(fmt.Sprintf("*** oops: virtAddr:%x, pteAddr:%x, pte:%x", uint(virtAddr), uint(pteAddr), uint(pte)))

	if !cpu.bus.read(pteAddr, &pte, 8) {
		cpu.trapWithTval(ExceptionLoadAccessFault, virtAddr)
		return
	}

	isLeaf := bit(pte, PteR) == 1 || bit(pte, PteX) == 1

	if bit(pte, PteV) == 0 || // valid bit not set
		bit(pte, PteR) == 0 && bit(pte, PteW) == 1 || // reserved
		isLeaf && bits(pte, 10, 18) != 0 { // misaligned gigapage
		return
	}

	if isLeaf {
		*targetPTE = pte
		*shift = 30
		return
	}

	pteAddr = bits(pte, 10, 44)<<12 | bits(virtAddr, 21, 9)<<3
	if !cpu.bus.read(pteAddr, &pte, 8) {
		cpu.trapWithTval(ExceptionLoadAccessFault, virtAddr)
		return
	}

	isLeaf = bit(pte, PteR) == 1 || bit(pte, PteX) == 1

	if bit(pte, PteV) == 0 || // valid bit not set
		bit(pte, PteR) == 0 && bit(pte, PteW) == 1 || // reserved
		isLeaf && bits(pte, 10, 9) != 0 { // misaligned megapage
		return
	}

	if isLeaf {
		*targetPTE = pte
		*shift = 21
		return
	}

	pteAddr = bits(pte, 10, 44)<<12 | bits(virtAddr, 12, 9)<<3
	if !cpu.bus.read(pteAddr, &pte, 8) {
		cpu.trapWithTval(ExceptionLoadAccessFault, virtAddr)
		return
	}

	if bit(pte, PteV) == 0 ||
		bit(pte, PteR) == 0 && !(bit(pte, PteW) == 0 && bit(pte, PteX) == 1) {
		return
	}

	*targetPTE = pte
	*shift = 12
}
