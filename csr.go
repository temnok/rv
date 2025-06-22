package rv

type CSR struct {
	satp                                                       int32
	mstatus, misa, medeleg, mideleg, mie, mtvec, mepc, mhartid int32

	// TODO: remove
	pmpcfg0, pmpaddr0, mnstatus int32
}

func (cpu *CPU) csrAccess(i int32) *int32 {
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
	case 0xf14:
		return &csr.mhartid

	// TODO: remove
	case 0x3A0:
		return &csr.pmpcfg0
	case 0x3B0:
		return &csr.pmpaddr0
	case 0x744:
		return &csr.mnstatus

	default:
		cpu.instrIllegal = true
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
		cpu.instrIllegal = true
		return false
	}

	ptr := cpu.csrAccess(i)
	if ptr == nil {
		return false
	}

	*ptr = val
	return true
}
