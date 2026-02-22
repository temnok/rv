package cpu

import (
	"github.com/temnok/rv/arch"
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

func New(startAddr int) *state.CPU {
	xl := 0b_10

	cpu := &state.CPU{
		StaticState: state.StaticState{
			Priv: state.PrivM,

			CSR: state.CSR{
				Misa: xl<<(arch.XLen-2) |
					1<<('i'-'a') | 1<<('m'-'a') | 1<<('a'-'a') | 1<<('c'-'a') |
					1<<('f'-'a') | ('d' - 'a') |
					1<<('u'-'a') | 1<<('s'-'a'),
			},
		},

		Update: state.UpdatedState{
			XReg: -1,
			CReg: -1,
		},
	}

	cpu.CSR.Mstatus = xl<<csr.MstatusSXL | xl<<csr.MstatusUXL
	cpu.Update.PC = startAddr

	return cpu
}
