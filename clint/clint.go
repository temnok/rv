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

	cpu.CSR.TimerCallbacks = append(cpu.CSR.TimerCallbacks, clint.sync)

	return clint
}

func (clint *CLINT) Access(addr int, width int, write bool, writeData int) int {
	addr &= 0xff_ffff

	switch addr {
	case 0x0000: // msip
		return clint.accessReg(&clint.msip, write, writeData)

	case 0x4000: // mtimecmp
		return clint.accessReg(&clint.mtimecmp, write, writeData)

	case 0xBFF8: // mtime
		return clint.accessReg(&clint.cpu.CSR.Time, write, writeData)
	}

	return 0
}

func (clint *CLINT) accessReg(reg *int, write bool, writeData int) int {
	val := *reg

	if write {
		*reg = writeData
		clint.sync()
	}

	return val
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
