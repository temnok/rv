package rv

import (
	"github.com/temnok/rv/exec"
	"github.com/temnok/rv/instr"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

var routes = []func(*state.CPU, instr.Op){
	0:  exec.Load,
	1:  exec.LoadFP,
	3:  exec.Fence,
	4:  exec.ComputeI,
	5:  instr.Auipc,
	6:  exec.ComputeI64,
	8:  exec.Store,
	9:  exec.StoreFP,
	11: exec.Atomic,
	12: exec.ComputeR,
	13: instr.Lui,
	14: exec.ComputeR64,
	16: exec.ComputeFP,
	17: exec.ComputeFP,
	18: exec.ComputeFP,
	19: exec.ComputeFP,
	20: exec.ComputeFP,
	24: exec.Branch,
	25: instr.Jalr,
	27: instr.Jal,
	28: exec.System,
}

func Exec(cpu *state.CPU, opcode int) {
	opcodeSize := 4
	if isCompressed := opcode&3 != 3; isCompressed {
		opcodeSize = 2

		if decompress(cpu, &opcode); trap.IsEntered(cpu) {
			return
		}
	}
	cpu.Update.PC = cpu.Xint(cpu.PC + opcodeSize)

	op := instr.Op(opcode)
	f5 := op.F5()

	if f5 >= len(routes) || routes[f5] == nil {
		trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
		return
	}

	routes[f5](cpu, op)
}
