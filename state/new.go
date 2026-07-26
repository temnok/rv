package state

func New(ramSize int) *CPU {
	return &CPU{
		CSR: CSR{
			Mstatus: 2<<MstatusSXL | 2<<MstatusUXL, // 64-bit S- and U-modes,
		},

		Priv: PrivM,

		RAM: make([]int, ramSize/8),

		UARTInput: func() (byte, bool) {
			return 0, false // does not produce
		},

		UARTOutput: func(byte) bool {
			return true // silently consumes
		},
	}
}
