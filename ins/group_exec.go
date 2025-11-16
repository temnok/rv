package ins

import (
	"github.com/temnok/rv/decompress"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

var insGroups = []func(*state.CPU, Op){
	0:  Load,
	1:  LoadFP,
	3:  Fence,
	4:  ComputeI,
	5:  Auipc,
	6:  ComputeI32,
	8:  Store,
	9:  StoreFP,
	11: Atomic,
	12: ComputeR,
	13: Lui,
	14: ComputeR32,
	16: ComputeFP,
	17: ComputeFP,
	18: ComputeFP,
	19: ComputeFP,
	20: ComputeFP,
	24: Branch,
	25: Jalr,
	27: Jal,
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
	cpu.Update.PC = cpu.Int(cpu.PC + opcodeSize)

	op := Op(opcode)
	f5 := op.F5()

	if f5 >= len(insGroups) || insGroups[f5] == nil {
		trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
		return
	}

	insGroups[f5](cpu, op)
}
