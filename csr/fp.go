package csr

import (
	"github.com/temnok/rv/state"
)

func FpDisabled(cpu *state.CPU) bool {
	return cpu.CSR.Mstatus>>MstatusFS&3 == FSoff
}
