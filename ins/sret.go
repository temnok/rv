package ins

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func sret(cpu *state.CPU, op Op) {
	trap.Exit(cpu, csr.PrivS)
}
