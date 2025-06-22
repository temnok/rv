package rv

type CSR struct {
	satp                                                             int32
	mstatus, misa, medeleg, mideleg, mie, mtvec, mepc, mcause, mtval int32
	mhartid                                                          int32
}

func (cpu *CPU) csrAccess(i int32) *int32 {
	if priv := bits(i, 8, 2); cpu.priv < priv {
		cpu.trap(ExceptionIllegalIstruction)
		return nil
	}

	csr := &cpu.csr

	switch i {
	case 0x180:
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
	case 0x341:
		return &csr.mepc
	case 0x342:
		return &csr.mcause
	case 0x343:
		return &csr.mtval
	case 0xf14:
		return &csr.mhartid

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
