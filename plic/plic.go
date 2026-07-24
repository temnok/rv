package plic

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
	"math/bits"
)

type PLIC struct {
	cpu *state.CPU

	pending int
	enable  int
}

func New(cpu *state.CPU) *PLIC {
	return &PLIC{
		cpu: cpu,
	}
}

func (plic *PLIC) Access(addr int, width int, write bool, writeData int) int {
	addr &= 0xff_ffff

	var val int

	switch addr {
	case 0x1000: // pending
		val = plic.pending

	case 0x2000: // enable
		val = plic.enable

		if write {
			plic.enable = writeData
			plic.sync()
		}

	case 0x20_0004: // claim
		if !write {
			val = bits.TrailingZeros(uint(plic.pending&plic.enable)) & 63
			plic.pending &^= 1 << val
		}

		plic.sync()
	}

	return val
}

func (plic *PLIC) PendInterrupt(interruptID int, enable bool) {
	if enable {
		plic.pending |= 1 << interruptID
	} else {
		plic.pending &^= 1 << interruptID
	}

	plic.sync()
}

func (plic *PLIC) sync() {
	if active := plic.pending & plic.enable; active != 0 {
		plic.cpu.CSR.Mip |= 1 << csr.MipSEIP
	} else {
		plic.cpu.CSR.Mip &^= 1 << csr.MipSEIP
	}
}
