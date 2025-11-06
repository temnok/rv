package instr

import (
	"github.com/temnok/rv/state"
)

// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#wfi
func Wfi(cpu *state.CPU, op Op) {
}
