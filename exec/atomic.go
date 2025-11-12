package exec

import (
	"github.com/temnok/rv/instr"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

var atomics = []func(*state.CPU, instr.Op){
	0:  instr.Amoadd,
	1:  instr.Amoswap,
	2:  instr.Lr,
	3:  instr.Sc,
	4:  instr.Amoxor,
	8:  instr.Amoor,
	12: instr.Amoand,
	16: instr.Amomin,
	20: instr.Amomax,
	24: instr.Amominu,
	28: instr.Amomaxu,
}

func Atomic(cpu *state.CPU, op instr.Op) {
	i := op.F7() >> 2
	if i >= len(atomics) || atomics[i] == nil {
		trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
		return
	}

	atomics[i](cpu, op)
}
