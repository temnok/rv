package rv

type CLINT struct {
	cpu                       *CPU
	mswi, mtimecmp, mtimecmph int32
}

func (clint *CLINT) init(cpu *CPU) {
	*clint = CLINT{
		cpu: cpu,
	}
}

func (clint *CLINT) access(addr int32, data *int32, write bool) bool {
	var reg *int32

	switch addr * 4 {
	case 0x0: // mswi
		reg = &clint.mswi
	case 0x4000 + 0x0000: // mtimecmp
		reg = &clint.mtimecmp
	case 0x4000 + 0x0000 + 4: // mtimecmph
		reg = &clint.mtimecmph
	case 0x4000 + 0x7FF8: // mtime
		reg = &clint.cpu.csr.mtime
	case 0x4000 + 0x7FF8 + 4: // mtimeh
		reg = &clint.cpu.csr.mtimeh
	default:
		return false
	}

	if write {
		*reg = *data
	} else {
		*data = *reg
	}

	return true
}

func (clint *CLINT) setPendingIterrupts() {
	csr := &clint.cpu.csr
	csr.mip |= bit(clint.mswi, 1) << mipMSI

	if csr.mtimeh > clint.mtimecmph ||
		csr.mtimeh == clint.mtimecmph && csr.mtime >= clint.mtimecmp {

		csr.mip |= 1 << mipMTI
	}
}
