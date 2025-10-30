package rv

import (
	"github.com/temnok/rv/instr"
	"github.com/temnok/rv/trap"
)

func (cpu *CPU) execBranch(op instr.Op) {
	switch op.F3() {
	case 0b_000:
		instr.Beq(cpu.State, op)
	case 0b_001:
		instr.Bne(cpu.State, op)
	case 0b_100:
		instr.Blt(cpu.State, op)
	case 0b_101:
		instr.Bge(cpu.State, op)
	case 0b_110:
		instr.Bltu(cpu.State, op)
	case 0b_111:
		instr.Bgeu(cpu.State, op)
	default:
		trap.EnterWithoutTval(cpu.State, ExceptionIllegalIstruction)
	}
}
