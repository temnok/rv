package clint

import (
	"github.com/temnok/rv/bi"
	"github.com/temnok/rv/csr"
	re "github.com/temnok/rv/reg"
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
		baseAddr: baseAddr,
	}
}

func (clint *CLINT) Access(addr int, data *int, width int, write bool) bool {
	if addr -= clint.baseAddr; addr < 0 || addr >= 0x10000 {
		return false
	}

	switch reg, offset := addr&^7, addr&7; reg {
	case 0x0: // mswi
		re.Access(&clint.mswi, data, offset, width, write)
	case 0x4000 + 0x0000: // mtimecmp
		re.Access(&clint.mtimecmp, data, offset, width, write)

		if write {
			clint.clearTimerInterrupt()
		}
	case 0x4000 + 0x7FF8: // mtime
		re.Access(&clint.cpu.CSR.Time, data, offset, width, write)
	}

	return true
}

func (clint *CLINT) clearTimerInterrupt() {
	if uint(clint.cpu.CSR.Time) < uint(clint.mtimecmp) {
		clint.cpu.CSR.Mip &^= 1<<csr.MipMTI | 1<<csr.MipSTI
	}
}

func (clint *CLINT) NotifyInterrupts() {
	reg := &clint.cpu.CSR

	if bi.T(clint.mswi, 1) == 1 {
		reg.Mip |= 1 << csr.MipMSI
	}

	if uint(reg.Time) >= uint(clint.mtimecmp) {
		reg.Mip |= 1 << csr.MipMTI
	}
}
