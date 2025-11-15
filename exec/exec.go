package exec

import (
	"github.com/temnok/rv/decompress"
	"github.com/temnok/rv/instr"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

var routes = []func(*state.CPU, instr.Op){
	0:  Load,
	1:  LoadFP,
	3:  Fence,
	4:  ComputeI,
	5:  instr.Auipc,
	6:  ComputeI64,
	8:  Store,
	9:  StoreFP,
	11: Atomic,
	12: ComputeR,
	13: instr.Lui,
	14: ComputeR64,
	16: ComputeFP,
	17: ComputeFP,
	18: ComputeFP,
	19: ComputeFP,
	20: ComputeFP,
	24: instr.Branch,
	25: instr.Jalr,
	27: instr.Jal,
	28: System,
}

func Exec(cpu *state.CPU, opcode int) {
	opcodeSize := 4
	if isCompressed := opcode&3 != 3; isCompressed {
		opcodeSize = 2

		if decompress.Decompress(cpu, &opcode); trap.IsEntered(cpu) {
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
