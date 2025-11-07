package clint

import (
	"github.com/temnok/rv/bi"
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

type CLINT struct {
	cpu      *state.CPU
	baseAddr int

	mswi, mtimecmp, mtimecmph int
}

func New(cpu *state.CPU, baseAddr int) *CLINT {
	return &CLINT{
		cpu:      cpu,
		baseAddr: cpu.Xint(baseAddr),
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
	reg := &clint.cpu.CSR

	if bi.T(clint.mswi, 1) == 1 {
		reg.Mip |= 1 << csr.MipMSI
	}

	if uint(reg.Timeh) > uint(clint.mtimecmph) ||
		reg.Timeh == clint.mtimecmph && uint(reg.Time) >= uint(clint.mtimecmp) {

		reg.Mip |= 1 << csr.MipMTI
	}
}
