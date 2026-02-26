package clint

import (
	"github.com/temnok/rv/csr"
	re "github.com/temnok/rv/reg"
	"github.com/temnok/rv/state"
)

type CLINT struct {
	cpu      *state.CPU
	baseAddr int

	msip, mtimecmp, mtimecmph int
}

func New(cpu *state.CPU, baseAddr int) *CLINT {
	clint := &CLINT{
		cpu:      cpu,
		baseAddr: baseAddr,
	}

	cpu.CSR.TimerCallbacks = append(cpu.CSR.TimerCallbacks, clint.notifyInterrupts)

	return clint
}

func (clint *CLINT) Access(addr int, data *int, width int, write bool) bool {
	if addr -= clint.baseAddr; addr < 0 || addr >= 0x10000 {
		return false
	}

	switch reg, offset := addr&^7, addr&7; reg {
	case 0x0: // msip
		re.Access(&clint.msip, data, offset, width, write)
		if write {
			clint.initSoftwareInterrupt()
		}

	case 0x4000: // mtimecmp
		re.Access(&clint.mtimecmp, data, offset, width, write)

		if write {
			clint.clearTimerInterrupt()
		}
	case 0xBFF8: // mtime
		re.Access(&clint.cpu.CSR.Time, data, offset, width, write)
	}

	return true
}

func (clint *CLINT) clearTimerInterrupt() {
	if uint(clint.cpu.CSR.Time) < uint(clint.mtimecmp) {
		clint.cpu.CSR.Mip &^= 1<<csr.MipMTI | 1<<csr.MipSTI
	}
}

func (clint *CLINT) initSoftwareInterrupt() {
	if clint.msip&1 == 0 {
		clint.cpu.CSR.Mip &^= 1 << csr.MipMSI
	} else {
		clint.cpu.CSR.Mip |= 1 << csr.MipMSI
	}
}

func (clint *CLINT) notifyInterrupts() {
	reg := &clint.cpu.CSR

	if uint(reg.Time) >= uint(clint.mtimecmp) {
		reg.Mip |= 1 << csr.MipMTI
	}
}
