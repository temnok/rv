package ins

import (
	"github.com/temnok/rv/extc"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

var (
	groups = []func(*state.CPU, Op){
		0:  execLoad,
		1:  execLoadFP,
		2:  illegal,
		3:  execFence,
		4:  execComputeI,
		5:  auipc,
		6:  execComputeI32,
		7:  illegal,
		8:  execStore,
		9:  execStoreFP,
		10: illegal,
		11: execAtomic,
		12: execComputeR,
		13: lui,
		14: execComputeR32,
		15: illegal,
		16: execComputeFP,
		17: execComputeFP,
		18: execComputeFP,
		19: execComputeFP,
		20: execComputeFP,
		21: illegal,
		22: illegal,
		23: illegal,
		24: execBranch,
		25: jalr,
		26: illegal,
		27: jal,
		28: execSystem,
		29: illegal,
		30: illegal,
		31: illegal,
	}
)

func Exec(cpu *state.CPU, opcode int) {
	cpu.InstrCount++

	opcodeSize := 4

	if isCompressed := opcode&3 != 3; isCompressed {
		cpu.CInstrCount++

		compressedOpcode := int(uint16(opcode))

		if opcode = extc.Decompress(compressedOpcode); opcode == 0 {
			trap.Enter(cpu, trap.IllegalIstruction, compressedOpcode)
			return
		}

		opcodeSize = 2
	}

	cpu.Update.PC = cpu.PC + opcodeSize

	op := Op(opcode)
	f5 := op.f5()

	groups[f5](cpu, op)
}
