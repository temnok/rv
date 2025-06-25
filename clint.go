package rv

type CLINT struct {
	cpu      *CPU
	baseAddr int32

	mswi, mtimecmp, mtimecmph int32

	AccessCount, InterruptCount int
}

func (clint *CLINT) Init(cpu *CPU, baseAddr int32) {
	*clint = CLINT{
		cpu:      cpu,
		baseAddr: int32(uint32(baseAddr) >> 2),
	}
}

func (clint *CLINT) access(addr int32, data *int32, write bool) bool {
	if addr -= clint.baseAddr; addr < 0 || addr >= 0x10000/4 {
		return false
	}

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

	clint.AccessCount++

	return true
}

func (clint *CLINT) notifyInterrupts() {
	csr := &clint.cpu.csr

	if bit(clint.mswi, 1) == 1 {
		csr.mip |= 1 << mipMSI
	}

	if uint32(csr.mtimeh) > uint32(clint.mtimecmph) ||
		csr.mtimeh == clint.mtimecmph && uint32(csr.mtime) >= uint32(clint.mtimecmp) {

		if csr.mip&(1<<mipMTI) == 0 {
			clint.InterruptCount++
		}

		csr.mip |= 1 << mipMTI
	}
}
