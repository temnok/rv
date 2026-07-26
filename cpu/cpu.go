package cpu

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

func New(ramSize, startAddr int) *state.CPU {
	return &state.CPU{
		PC: startAddr,

		Priv: state.PrivM,

		CSR: state.CSR{
			Mstatus: 2<<csr.MstatusSXL | 2<<csr.MstatusUXL,
		},

		RAM: make([]int, ramSize/8),

		UARTInput: func() (byte, bool) {
			return 0, false
		},

		UARTOutput: func(byte) bool {
			return true
		},
	}
}
