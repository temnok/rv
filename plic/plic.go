package plic

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
	"math/bits"
)

type PLIC struct {
	cpu      *state.CPU
	baseAddr int

	pending int
	enable  int
	claim   int
}

func New(cpu *state.CPU, baseAddr int) *PLIC {
	return &PLIC{
		cpu:      cpu,
		baseAddr: baseAddr,
	}
}

func (plic *PLIC) Access(addr int, data *int, width int, write bool) bool {
	if addr = addr - plic.baseAddr; addr < 0 || addr >= 0x200004+4 {
		return false
	}

	switch addr {
	case 0x1000: // pending
		if read := !write; read {
			*data = plic.pending
		}

	case 0x2000: // enable
		if write {
			plic.enable = *data

			plic.sync()
		} else {
			*data = plic.enable
		}

	case 0x200004: // claim
		if write {
			plic.claim = *data
		} else {
			*data = plic.claim
			plic.pending &^= 1 << plic.claim
		}

		plic.sync()
	}

	return true
}

func (plic *PLIC) TriggerInterrupt(source int) {
	plic.pending |= 1 << source
	plic.sync()
}

func (plic *PLIC) sync() {
	if active := plic.pending & plic.enable; active != 0 {
		plic.claim = bits.TrailingZeros(uint(active))
		plic.cpu.CSR.Mip |= 1 << csr.MipSEI
	} else {
		plic.claim = 0
		plic.cpu.CSR.Mip &^= 1 << csr.MipSEI
	}
}
