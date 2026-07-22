package cpu

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

const (
	Misa = -1<<63 |
		1<<('i'-'a') | 1<<('m'-'a') | 1<<('a'-'a') | 1<<('c'-'a') |
		1<<('f'-'a') | ('d' - 'a') |
		1<<('u'-'a') | 1<<('s'-'a')
)

func New(startAddr int) *state.CPU {
	xl := 0b_10

	cpu := &state.CPU{
		StaticState: state.StaticState{
			PC: startAddr,

			Priv: state.PrivM,

			CSR: state.CSR{
				Mstatus: xl<<csr.MstatusSXL | xl<<csr.MstatusUXL,
			},
		},
	}

	return cpu
}
