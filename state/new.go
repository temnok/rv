package state

import "github.com/temnok/rv/csr"

func New(ramSize int) *CPU {
	return &CPU{
		CSR: csr.Registers{
			Priv:    csr.PrivM,
			Mstatus: 2<<csr.MstatusSXL | 2<<csr.MstatusUXL, // 64-bit S- and U-modes,
		},

		RAM: make([]int, ramSize/8),

		UARTInput: func() (byte, bool) {
			return 0, false // does not produce
		},

		UARTOutput: func(byte) bool {
			return true // silently consumes
		},
	}
}
