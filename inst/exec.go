package inst

import (
	"github.com/temnok/rv/extc"
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

type context state.CPU

var (
	groups = []func(*state.CPU, Op){
		0:  execLoad,
		1:  execLoadFP,
		2:  illegal,
		3:  execFence,
		4:  execComputeI,
		6:  execComputeI32,
		7:  illegal,
		8:  execStore,
		9:  execStoreFP,
		10: illegal,
		11: execAtomic,
		12: execComputeR,
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
	opcodeSize := 4

	if isCompressed := opcode&3 != 3; isCompressed {
		compressedOpcode := int(uint16(opcode))

		if opcode = extc.Decompress(compressedOpcode); opcode == 0 {
			trap.Enter(cpu, trap.IllegalIstruction, compressedOpcode)
			return
		}

		opcodeSize = 2
	}

	cpu.Update.FollowPC = cpu.PC + opcodeSize

	ctx := (*context)(cpu)
	op := Op(opcode)
	f5, rd := op.f5(), op.rd()

	switch f5 {
	case 5:
		ctx.AUIPC(rd, imm.U(opcode))
	case 13:
		ctx.LUI(rd, imm.U(opcode))
	default:
		groups[f5](cpu, op)
	}

	if cpu.Update.Targets&state.UpdatePC == 0 {
		cpu.Update.Targets |= state.UpdatePC
		cpu.Update.PC = cpu.Update.FollowPC
	}
}
