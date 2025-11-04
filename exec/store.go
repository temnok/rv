package exec

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/instr"
	"github.com/temnok/rv/mem"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func Store(cpu *state.State, op instr.Op) {
	imm, rs1, rs2 := imm.S(op.Code()), op.Rs1(), op.Rs2()

	switch op.F3() {
	case 0b_000: // sb
		mem.Write(cpu, cpu.X[rs1]+imm, int(uint8(cpu.X[rs2])), 1)

	case 0b_001: // sh
		mem.Write(cpu, cpu.X[rs1]+imm, int(uint16(cpu.X[rs2])), 2)

	case 0b_010: // sw
		mem.Write(cpu, cpu.X[rs1]+imm, int(uint32(cpu.X[rs2])), 4)

	case 0b_011: // sd
		if !cpu.Xlen64() {
			trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
			return
		}

		mem.Write(cpu, cpu.X[rs1]+imm, cpu.X[rs2], 8)

	default:
		trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
	}
}
