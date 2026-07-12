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
	claim   int
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

	case 0x200004: // claim
		val = plic.claim

		if write {
			plic.claim = writeData
		} else {
			plic.pending &^= 1 << plic.claim
		}

		plic.sync()
	}

	return val
}

func (plic *PLIC) PendInterrupt(source int, enable bool) {
	if mask := 1 << source; enable {
		plic.pending |= mask
	} else {
		plic.pending &^= mask
	}

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
