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
	case 0:
		cpu.execLoad(op)
	case 1:
		cpu.execLoadFP(op)
	case 3:
		exec.Fence(cpu.State, op)
	case 4:
		exec.ComputeI(cpu.State, op)
	case 5:
		instr.Auipc(cpu.State, op)
	case 6:
		exec.ComputeI64(cpu.State, op)
	case 8:
		cpu.execStore(op)
	case 9:
		cpu.execStoreFP(op)
	case 11:
		cpu.execAtomic(op)
	case 12:
		exec.ComputeR(cpu.State, op)
	case 13:
		instr.Lui(cpu.State, op)
	case 14:
		exec.ComputeR64(cpu.State, op)
	case 16, 17, 18, 19, 20:
		cpu.execComputeFP(op)
	case 24:
		exec.Branch(cpu.State, op)
	case 25:
		instr.Jalr(cpu.State, op)
	case 27:
		instr.Jal(cpu.State, op)
	case 28:
		exec.System(cpu.State, op)
	default:
		trap.EnterWithoutTval(cpu.State, trap.IllegalIstruction)
	}
}
