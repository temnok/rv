package rv

type CLINT struct {
	cpu      *CPU
	baseAddr int

	mswi, mtimecmp, mtimecmph int
}

func (clint *CLINT) Init(cpu *CPU, baseAddr int) {
	*clint = CLINT{
		cpu:      cpu,
		baseAddr: cpu.xint(baseAddr),
	}
}

func (clint *CLINT) Access(addr int, data *int, width int, write bool) bool {
	if addr = (addr - clint.baseAddr) / 4; addr < 0 || addr >= 0x10000/4 || width < 4 {
		return false
	}

	var reg *int

	switch addr * 4 {
	case 0x0: // mswi
		reg = &clint.mswi
	case 0x4000 + 0x0000: // mtimecmp
		reg = &clint.mtimecmp
	case 0x4000 + 0x0000 + 4: // mtimecmph
		reg = &clint.mtimecmph
	case 0x4000 + 0x7FF8: // mtime
		reg = &clint.cpu.CSR.Time
	case 0x4000 + 0x7FF8 + 4: // mtimeh
		reg = &clint.cpu.CSR.Timeh
	}

	if write {
		if reg != nil {
			*reg = *data
		}
	} else {
		if reg != nil {
			*data = *reg
		} else {
			*data = 0
		}
	}

	return true
}

func (clint *CLINT) NotifyInterrupts() {
	csr := &clint.cpu.CSR

	if bit(clint.mswi, 1) == 1 {
		csr.Mip |= 1 << MipMSI
	}

	if uint(csr.Timeh) > uint(clint.mtimecmph) ||
		csr.Timeh == clint.mtimecmph && uint(csr.Time) >= uint(clint.mtimecmp) {

		csr.Mip |= 1 << MipMTI
	}
}
