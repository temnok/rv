package clint

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

type CLINT struct {
	cpu *state.CPU

	msip, mtimecmp int
}

func New(cpu *state.CPU) *CLINT {
	clint := &CLINT{
		cpu: cpu,
	}

	cpu.CSR.TimerCallbacks = append(cpu.CSR.TimerCallbacks, clint.Propagate)

	return clint
}

func (clint *CLINT) Access(addr int, width int, write bool, writeData int) int {
	addr &= 0xff_ffff

	var reg *int

	switch addr {
	case 0x0000: // msip
		reg = &clint.msip

	case 0x4000: // mtimecmp
		reg = &clint.mtimecmp

	case 0xBFF8: // mtime
		reg = &clint.cpu.CSR.Time

	default:
		return 0
	}

	val := *reg

	if write {
		*reg = writeData
		clint.Propagate()
	}

	return val
}

func (clint *CLINT) Propagate() {
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
