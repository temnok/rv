package csr

import (
	"github.com/temnok/rv/bi"
	"github.com/temnok/rv/state"
)

func FpDisabled(cpu *state.CPU) bool {
	return bi.Ts(cpu.CSR.Mstatus, MstatusFS, 2) == FSoff
}
