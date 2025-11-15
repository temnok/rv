package instr

import (
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

var atomics = []func(*state.CPU, Op){
	0:  amoadd,
	1:  amoswap,
	2:  lr,
	3:  sc,
	4:  amoxor,
	8:  amoor,
	12: amoand,
	16: amomin,
	20: amomax,
	24: amominu,
	28: amomaxu,
}

func Atomic(cpu *state.CPU, op Op) {
	i := op.F7() >> 2
	if i >= len(atomics) || atomics[i] == nil {
		trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
		return
	}

	atomics[i](cpu, op)
}
