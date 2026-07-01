package clint

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/reg"
	"github.com/temnok/rv/state"
)

type CLINT struct {
	cpu      *state.CPU
	baseAddr int

	msip, mtimecmp int
}

func New(cpu *state.CPU, baseAddr int) *CLINT {
	clint := &CLINT{
		cpu:      cpu,
		baseAddr: baseAddr,
	}

	cpu.CSR.TimerCallbacks = append(cpu.CSR.TimerCallbacks, clint.sync)

	return clint
}

func (clint *CLINT) Access(addr int, data *int, width int, write bool) bool {
	if addr -= clint.baseAddr; addr < 0 || addr >= 0x10000 || addr&7 != 0 {
		return false
	}

	switch addr {
	case 0x0000: // msip
		reg.Access(&clint.msip, data, 0, width, write)

		if write {
			clint.sync()
		}

	case 0x4000: // mtimecmp
		reg.Access(&clint.mtimecmp, data, 0, width, write)

		if write {
			clint.sync()
		}

	case 0xBFF8: // mtime
		reg.Access(&clint.cpu.CSR.Time, data, 0, width, write)
	}

	return true
}

func (clint *CLINT) sync() {
	mip := clint.cpu.CSR.Mip

	if clint.msip&1 == 0 {
		mip &^= 1 << csr.MipMSI
	} else {
		mip |= 1 << csr.MipMSI
	}

	if uint(clint.cpu.CSR.Time) < uint(clint.mtimecmp) {
		mip &^= 1<<csr.MipMTI | 1<<csr.MipSTI
	} else {
		mip |= 1<<csr.MipMTI | 1<<csr.MipSTI
	}

	clint.cpu.CSR.Mip = mip
}
