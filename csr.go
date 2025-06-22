package rv

// https://github.com/riscv/riscv-isa-manual/blob/main/src/priv-csrs.adoc#user-content-mcsrnames0

type CSR struct {
	cycle, mtime, cycleh, mtimeh int32

	mstatus, misa, medeleg, mideleg, mie, mtvec, mcounteren, mstatush, medelegh int32
	mscratch, mepc, mcause, mtval, mip                                          int32
	mvendorid, marchid, mimpid, mhartid, mconfigptr                             int32
}

func (cpu *CPU) csrAccess(i int32) *int32 {
	if priv := bits(i, 8, 2); cpu.priv < priv {
		cpu.trap(ExceptionIllegalIstruction)
		return nil
	}

	csr := &cpu.csr

	switch i {
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
		return &csr.medelegh

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
		return &cpu.csr.cycleh
	case 0xC81:
		return &cpu.csr.mtimeh
	case 0xC82:
		return &cpu.csr.cycleh

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

	default:
		cpu.trap(ExceptionIllegalIstruction)
		return nil
	}
}

func (cpu *CPU) csrRead(i int32, val *int32) bool {
	ptr := cpu.csrAccess(i)
	if ptr == nil {
		return false
	}

	*val = *ptr
	return true
}

func (cpu *CPU) csrWrite(i, val int32) bool {
	if readOnly := bits(i, 10, 2) == 0b_11; readOnly {
		cpu.trap(ExceptionIllegalIstruction)
		return false
	}

	ptr := cpu.csrAccess(i)
	if ptr == nil {
		return false
	}

	*ptr = val
	return true
}
