package cpu

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

func New(xlen int, startAddr int) *state.CPU {
	xl := xlen / 32

	cpu := &state.CPU{
		FixedState: state.FixedState{
			Xlen: xlen,
		},

		StaticState: state.StaticState{
			Priv: state.PrivM,

			CSR: state.CSR{
				Misa: xl<<(xlen-2) |
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

	cpu.CSR.Mstatus = cpu.Xint(xl<<csr.MstatusSXL | xl<<csr.MstatusUXL)
	cpu.Update.PC = cpu.Xint(startAddr)

	return cpu
}
