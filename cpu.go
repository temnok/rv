package rv

import "time"

type CPU struct {
	x          [32]int32
	pc, nextPC int32
	csr        CSR

	trapped bool

	reserved        bool
	reservedAddress int32

	priv int32

	startNanoseconds int64

	bus Bus
}

// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#mcauses
const (
	//ExceptionInstructionAddressMisaligned = 0
	ExceptionInstructionAccessFault    = 1
	ExceptionIllegalIstruction         = 2
	ExceptionBreakpoint                = 3
	ExceptionLoadAddressMisaligned     = 4
	ExceptionLoadAccessFault           = 5
	ExceptionStoreAMOAddressMisaligned = 6
	ExceptionStoreAMOAccessFault       = 7
	ExceptionEnvironmentCallFromUMode  = 8
	ExceptionEnvironmentCallFromSMode  = 9
	ExceptionEnvironmentCallFromMMode  = 11
	ExceptionInstructionPageFault      = 12
	//ExceptionLoadPageFault                = 13
	//ExceptionStoreAMOPageFault            = 15

	PrivU = 0
	PrivS = 1
	PrivM = 3

	AccessExecute = 0
	AccessRead    = 1
	AccessWrite   = 3

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

func (cpu *CPU) init(startAddr int32, bus Bus) {
	const xlen32bit = 0b_01

	*cpu = CPU{
		pc:   startAddr,
		priv: PrivM,
		csr: CSR{
			misa: xlen32bit<<30 |
				1<<('i'-'a') | 1<<('m'-'a') | 1<<('a'-'a') | 1<<('c'-'a') |
				1<<('u'-'a') | 1<<('s'-'a'),
		},
		startNanoseconds: time.Now().UnixNano(),

		bus: bus,
	}
}

func (cpu *CPU) step() {
	cpu.trapped = false

	cpu.updateTimers()

	cpu.bus.notifyInterrupts()
	if cpu.trapOnPendingInterrupts() {
		cpu.pc = cpu.nextPC
		return
	}

	if opcode := int32(0); cpu.memFetch(cpu.pc, &opcode) {
		if compressed := opcode&0b11 != 0b11; compressed {
			opcode = decompress(opcode)
			cpu.nextPC = cpu.pc + 2
		} else {
			cpu.nextPC = cpu.pc + 4
		}

		cpu.exec(opcode)
	}

	cpu.pc = cpu.nextPC
}

func (cpu *CPU) updateTimers() {
	if cpu.csr.cycle++; cpu.csr.cycle == 0 {
		cpu.csr.cycleh++
	}

	ticks := (int64(time.Now().Nanosecond()) - cpu.startNanoseconds) / 10
	cpu.csr.mtime = int32(ticks)
	cpu.csr.mtimeh = int32(ticks >> 32)
}

// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#privstack
func (cpu *CPU) trapOnPendingInterrupts() bool {
	mi := cpu.csr.mip & cpu.csr.mie

	if mi == 0 {
		return false
	}

	for i := int32(12); i > 0; i-- {
		if bit(mi, i) == 0 {
			continue
		}

		priv := int32(PrivM)
		if bit(cpu.csr.mideleg, i) != 0 {
			priv = PrivS
		}

		if (priv == cpu.priv && bit(cpu.csr.mstatus, priv) != 0) || priv > cpu.priv {
			cpu.trap(-1<<mcauseI | i)
			return true
		}
	}

	return false
}

func (cpu *CPU) trap(cause int32) {
	if cpu.trapped {
		panic("double trap")
	}

	cpu.trapped = true

	isInterrupt := bit(cause, mcauseI) != 0
	causeID := bits(cause, 0, 5)

	deleg := cpu.csr.medeleg
	if isInterrupt {
		deleg = cpu.csr.mideleg
	}

	effectivePriv := int32(PrivM)
	if cpu.priv < PrivM && bit(deleg, causeID) != 0 {
		effectivePriv = PrivS
	}

	var tvec int32

	switch effectivePriv {
	case PrivM:
		mie := bit(cpu.csr.mstatus, mstatusMIE)
		cpu.csr.mstatus &^= 0b_11<<mstatusMPP | 1<<mstatusMPIE | 1<<mstatusMIE
		cpu.csr.mstatus |= cpu.priv<<mstatusMPP | mie<<mstatusMPIE

		cpu.csr.mepc = cpu.pc
		cpu.csr.mcause = cause
		cpu.csr.mtval = 0
		tvec = cpu.csr.mtvec

	case PrivS:
		sie := bit(cpu.csr.mstatus, mstatusSIE)
		cpu.csr.mstatus &^= 1<<mstatusSPP | 1<<mstatusSPIE | 1<<mstatusSIE
		cpu.csr.mstatus |= cpu.priv<<mstatusSPP | sie<<mstatusSPIE

		cpu.csr.sepc = cpu.pc
		cpu.csr.scause = cause
		cpu.csr.stval = 0
		tvec = cpu.csr.stvec
	}

	cpu.nextPC = tvec &^ 0b_11
	if bit(tvec, 0) != 0 && isInterrupt {
		cpu.nextPC += causeID * 4
	}

	cpu.priv = effectivePriv
}

// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#otherpriv
// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#privstack
func (cpu *CPU) ret(priv int32) {
	// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#virt-control
	trapSRET := cpu.priv == PrivS && bit(cpu.csr.mstatus, mstatusTSR) != 0

	if priv > cpu.priv || trapSRET {
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	switch priv {
	case PrivM:
		cpu.nextPC = cpu.csr.mepc
		cpu.priv = bits(cpu.csr.mstatus, mstatusMPP, 2)

		mie := bit(cpu.csr.mstatus, mstatusMPIE)
		cpu.csr.mstatus |= 1<<mstatusMPIE | mie<<mstatusMIE
		cpu.csr.mstatus &^= 0b_11 << mstatusMPP

	case PrivS:
		cpu.nextPC = cpu.csr.sepc
		cpu.priv = bits(cpu.csr.mstatus, mstatusSPP, 1)

		sie := bit(cpu.csr.mstatus, mstatusSPIE)
		cpu.csr.mstatus |= 1<<mstatusSPIE | sie<<mstatusSIE
		cpu.csr.mstatus &^= 1 << mstatusSPP
	}

	if cpu.priv != PrivM {
		cpu.csr.mstatus &^= 1 << mstatusMPRV
	}
}

func (cpu *CPU) translateSv32(virtAddr int32, physAddr *int32, access int32) bool {
	// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_memory_privilege_in_mstatus_register
	epriv := cpu.priv
	if bit(cpu.csr.mstatus, mstatusMPRV) != 0 && access != AccessExecute {
		epriv = bits(cpu.csr.mstatus, mstatusMPP, 2)
	}

	// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#satp-mode
	if bit(cpu.csr.satp, satpMODE) == 0 || epriv == PrivM {
		*physAddr = int32(uint32(virtAddr) >> 2)
		return true
	}

	// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#sv32algorithm
	var pte int32
	pteAddr := bits(cpu.csr.satp, 0, 22)<<10 | bits(virtAddr, 22, 10)
	if !cpu.bus.read(pteAddr, &pte) {
		cpu.trap(ExceptionLoadAccessFault)
		return false
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
			cpu.trap(ExceptionLoadAccessFault)
			return false
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

		cpu.trap(ExceptionInstructionPageFault + access)
		return false
	}

	*physAddr = pte&^0x3FF | lowAddrBits
	return true
}
