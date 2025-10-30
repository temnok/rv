package rv

import (
	"github.com/temnok/rv/exec"
	"github.com/temnok/rv/instr"
	"github.com/temnok/rv/trap"
)

func (cpu *CPU) exec(opcode int) {
	opcodeSize := 4
	if isCompressed := opcode&3 != 3; isCompressed {
		opcodeSize = 2

		if cpu.decompress(&opcode); trap.IsEntered(cpu.State) {
			return
		}
	}
	cpu.Update.PC = cpu.Xint(cpu.PC + opcodeSize)

	switch op := instr.Op(opcode); op.F5() {
	case 0b_00000:
		cpu.execLoad(op)
	case 0b_00001:
		cpu.execLoadFP(op)
	case 0b_00011:
		cpu.execFence(op)
	case 0b_00100:
		exec.ComputeI(cpu.State, op)
	case 0b_00110:
		cpu.execComputeI64(op)
	case 0b_00101:
		instr.Auipc(cpu.State, op)
	case 0b_01000:
		cpu.execStore(op)
	case 0b_01001:
		cpu.execStoreFP(op)
	case 0b_01011:
		cpu.execAtomic(op)
	case 0b_01100:
		exec.ComputeR(cpu.State, op)
	case 0b_01110:
		cpu.execComputeR64(op)
	case 0b_01101:
		instr.Lui(cpu.State, op)
	case 0b_10000, 0b_10001, 0b_10010, 0b_10011, 0b_10100:
		cpu.execComputeFP(op)
	case 0b_11000:
		cpu.execBranch(op)
	case 0b_11001:
		instr.Jalr(cpu.State, op)
	case 0b_11011:
		instr.Jal(cpu.State, op)
	case 0b_11100:
		cpu.execSystem(op)
	default:
		trap.EnterWithoutTval(cpu.State, trap.IllegalIstruction)
	}
}
