package rv

type CLINT struct {
	cpu      *CPU
	baseAddr Xint

	mswi, mtimecmp, mtimecmph Xint
}

func (clint *CLINT) Init(cpu *CPU, baseAddr Xint) {
	*clint = CLINT{
		cpu:      cpu,
		baseAddr: baseAddr,
	}
}

func (clint *CLINT) access(addr Xint, data *Xint, width Xint, write bool) bool {
	if addr = (addr - clint.baseAddr) / 4; addr < 0 || addr >= 0x10000/4 || width < 4 {
		return false
	}

	var reg *Xint

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

func (clint *CLINT) notifyInterrupts() {
	csr := &clint.cpu.csr

	if bit(clint.mswi, 1) == 1 {
		csr.mip |= 1 << mipMSI
	}

	if Xuint(csr.mtimeh) > Xuint(clint.mtimecmph) ||
		csr.mtimeh == clint.mtimecmph && Xuint(csr.mtime) >= Xuint(clint.mtimecmp) {

		csr.mip |= 1 << mipMTI
	}
}
