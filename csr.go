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

	mcauseI = Xlen - 1

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

func (cpu *CPU) csrAccess(i Xint, val *Xint, write bool) {
	if write && bits(i, 10, 2) == 0b_11 {
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	if priv := bits(i, 8, 2); cpu.priv < priv {
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	if i == satp && cpu.priv == PrivS && bit(cpu.csr.mstatus, mstatusTVM) != 0 {
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	csr := &cpu.csr
	writeVal := *val
	var reg *Xint

	switch i {
	case 0x100:
		reg = &csr.mstatus

	case 0x104:
		reg = &csr.mie

	case 0x105:
		reg = &csr.stvec

	case 0x106:
		reg = &csr.scounteren

	case 0x140:
		reg = &csr.sscratch

	case 0x141:
		reg = &csr.sepc

	case 0x142:
		reg = &csr.scause

	case 0x143:
		reg = &csr.stval

	case 0x144:
		reg = &csr.sip

	case satp:
		reg = &csr.satp

	case 0x300:
		reg = &csr.mstatus

	case 0x301:
		ignoreWrites := csr.misa
		reg = &ignoreWrites

	case 0x302:
		reg = &csr.medeleg

	case 0x303:
		reg = &csr.mideleg

	case 0x304:
		reg = &csr.mie

	case 0x305:
		reg = &csr.mtvec

	case 0x306:
		reg = &csr.mcounteren

	case 0x310:
		if Xlen32 {
			reg = &csr.mstatush
		}

	case 0x312:
		if Xlen32 {
			reg = &csr.medelegh
		}

	case 0x340:
		reg = &csr.mscratch

	case 0x341:
		reg = &csr.mepc

	case 0x342:
		reg = &csr.mcause

	case 0x343:
		reg = &csr.mtval

	case 0x344:
		reg = &csr.mip

	case 0xC00:
		reg = &csr.cycle

	case 0xC01:
		reg = &csr.mtime

	case 0xC02:
		reg = &csr.cycle

	case 0xC80:
		if Xlen32 {
			reg = &csr.cycleh
		}

	case 0xC81:
		if Xlen32 {
			reg = &csr.mtimeh
		}

	case 0xC82:
		if Xlen32 {
			reg = &csr.cycleh
		}

	case 0xF11:
		reg = &csr.mvendorid

	case 0xF12:
		reg = &csr.marchid

	case 0xF13:
		reg = &csr.mimpid

	case 0xF14:
		reg = &csr.mhartid

	case 0xF15:
		reg = &csr.mconfigptr
	}

	if reg == nil {
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	if write {
		*reg = writeVal
	} else {
		*val = *reg
	}
}

func (cpu *CPU) csrRead(i Xint, val *Xint) {
	cpu.csrAccess(i, val, false)
}

func (cpu *CPU) csrWrite(i, val Xint) {
	cpu.csrAccess(i, &val, true)
}
