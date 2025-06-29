package rv

// https://github.com/riscv/riscv-isa-manual/blob/main/src/priv-csrs.adoc#user-content-mcsrnames0

const (
	// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#misabase
	misaMXL = Xlen - 2

	// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_status_mstatus_and_mstatush_registers
	mstatusSXL = 34
	mstatusUXL = 32

	// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#satp
	satp     = 0x180
	satpMODE = Xlen - 1 - (Xlen/64)*3

	// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_status_mstatus_and_mstatush_registers
	mstatusSIE  = 1
	mstatusMIE  = 3
	mstatusSPIE = 5
	mstatusMPIE = 7
	mstatusSPP  = 8
	mstatusMPP  = 11
	mstatusMPRV = 17
	mstatusSUM  = 18
	mstatusMXR  = 19
	mstatusTVM  = 20
	mstatusTSR  = 22

	mcauseI = 31

	mipMSI = 3
	mipMTI = 7
	mipSEI = 9
)

type CSR struct {
	cycle, mtime, cycleh, mtimeh Xint

	sie, stvec, scounteren, sscratch, sepc, scause, stval, sip, satp Xint

	mstatus, misa, medeleg, mideleg, mie, mtvec, mcounteren, mstatush, medelegh Xint
	mscratch, mepc, mcause, mtval, mip                                          Xint
	mvendorid, marchid, mimpid, mhartid, mconfigptr                             Xint
}

func (cpu *CPU) csrAccess(i Xint) *Xint {
	if priv := bits(i, 8, 2); cpu.priv < priv {
		cpu.trap(ExceptionIllegalIstruction)
		return nil
	}

	if i == satp && cpu.priv == PrivS && bit(cpu.csr.mstatus, mstatusTVM) != 0 {
		cpu.trap(ExceptionIllegalIstruction)
		return nil
	}

	csr := &cpu.csr

	switch i {
	case 0x100:
		return &cpu.csr.mstatus

	case 0x104:
		return &cpu.csr.mie

	case 0x105:
		return &cpu.csr.stvec

	case 0x106:
		return &cpu.csr.scounteren

	case 0x140:
		return &cpu.csr.sscratch

	case 0x141:
		return &cpu.csr.sepc

	case 0x142:
		return &cpu.csr.scause

	case 0x143:
		return &cpu.csr.stval

	case 0x144:
		return &cpu.csr.sip

	case satp:
		return &csr.satp

	case 0x300:
		return &csr.mstatus

	case 0x301:
		return &csr.misa

	case 0x302:
		return &csr.medeleg

	case 0x303:
		return &csr.mideleg

	case 0x304:
		return &csr.mie

	case 0x305:
		return &csr.mtvec

	case 0x306:
		return &csr.mcounteren

	case 0x310:
		return &csr.mstatush

	case 0x312:
		if XlenIs32 {
			return &csr.medelegh
		}

	case 0x340:
		return &csr.mscratch

	case 0x341:
		return &csr.mepc

	case 0x342:
		return &csr.mcause

	case 0x343:
		return &csr.mtval

	case 0x344:
		return &csr.mip

	case 0xC00:
		return &cpu.csr.cycle

	case 0xC01:
		return &cpu.csr.mtime

	case 0xC02:
		return &cpu.csr.cycle

	case 0xC80:
		if XlenIs32 {
			return &cpu.csr.cycleh
		}

	case 0xC81:
		if XlenIs32 {
			return &cpu.csr.mtimeh
		}

	case 0xC82:
		if XlenIs32 {
			return &cpu.csr.cycleh
		}

	case 0xF11:
		return &csr.mvendorid

	case 0xF12:
		return &csr.marchid

	case 0xF13:
		return &csr.mimpid

	case 0xF14:
		return &csr.mhartid

	case 0xF15:
		return &csr.mconfigptr
	}

	cpu.trap(ExceptionIllegalIstruction)
	return nil
}

func (cpu *CPU) csrRead(i Xint, val *Xint) {
	ptr := cpu.csrAccess(i)
	if cpu.isTrapped {
		return
	}

	*val = *ptr
}

func (cpu *CPU) csrWrite(i, val Xint) {
	if readOnly := bits(i, 10, 2) == 0b_11; readOnly {
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	ptr := cpu.csrAccess(i)
	if cpu.isTrapped {
		return
	}

	csr := &cpu.csr

	switch ptr {
	case &csr.misa:
		return // ignore writes

	case &csr.mstatus:
		const mask = 0b_11<<mstatusSXL | 0b_11<<mstatusUXL
		val = Xint(int64(val)&^mask | int64(csr.mstatus)&mask) // copy SXL/UXL

	case &csr.satp:
		if XlenIs64 {
			const mask = 0b_0111 << satpMODE
			val = Xint(int64(val) &^ mask)
		}
	}

	*ptr = val
}
