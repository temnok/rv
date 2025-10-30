package rv

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/instr"
	"github.com/temnok/rv/trap"
)

func (cpu *CPU) execLoad(op instr.Op) {
	imm, rd, rs1 := imm.I(op.Code()), op.Rd(), op.Rs1()

	var val int

	switch op.F3() {
	case 0b_000: // lb
		cpu.memRead(cpu.X[rs1]+imm, &val, 1)
		cpu.Xset(rd, int(int8(val)))

	case 0b_001: // lh
		cpu.memRead(cpu.X[rs1]+imm, &val, 2)
		cpu.Xset(rd, int(int16(val)))

	case 0b_010: // lw
		cpu.memRead(cpu.X[rs1]+imm, &val, 4)
		cpu.Xset(rd, int(int32(val)))

	case 0b_011: // ld
		if !cpu.Xlen64() {
			trap.EnterWithoutTval(cpu.State, trap.IllegalIstruction)
			return
		}

		cpu.memRead(cpu.X[rs1]+imm, &val, 8)
		cpu.Xset(rd, val)

	case 0b_100: // lbu
		cpu.memRead(cpu.X[rs1]+imm, &val, 1)
		cpu.Xset(rd, int(uint8(val)))

	case 0b_101: // lhu
		cpu.memRead(cpu.X[rs1]+imm, &val, 2)
		cpu.Xset(rd, int(uint16(val)))

	case 0b_110: // lwu
		if !cpu.Xlen64() {
			trap.EnterWithoutTval(cpu.State, trap.IllegalIstruction)
			return
		}

		cpu.memRead(cpu.X[rs1]+imm, &val, 4)
		cpu.Xset(rd, int(uint32(val)))
	}

	if cpu.Update.XReg < 0 && !trap.IsEntered(cpu.State) {
		trap.EnterWithoutTval(cpu.State, trap.IllegalIstruction)
	}
}

func (cpu *CPU) execStore(op instr.Op) {
	imm, rs1, rs2 := imm.S(op.Code()), op.Rs1(), op.Rs2()

	switch op.F3() {
	case 0b_000: // sb
		cpu.memWrite(cpu.X[rs1]+imm, int(uint8(cpu.X[rs2])), 1)

	case 0b_001: // sh
		cpu.memWrite(cpu.X[rs1]+imm, int(uint16(cpu.X[rs2])), 2)

	case 0b_010: // sw
		cpu.memWrite(cpu.X[rs1]+imm, int(uint32(cpu.X[rs2])), 4)

	case 0b_011: // sd
		if !cpu.Xlen64() {
			trap.EnterWithoutTval(cpu.State, trap.IllegalIstruction)
			return
		}

		cpu.memWrite(cpu.X[rs1]+imm, cpu.X[rs2], 8)

	default:
		trap.EnterWithoutTval(cpu.State, trap.IllegalIstruction)
	}
}

func (cpu *CPU) execFence(op instr.Op) {
	imm, rd, rs1 := imm.I(op.Code()), op.Rd(), op.Rs1()

	switch op.F3() {
	case 0b_000: // fence
		if (imm&^0b_1111_1111) != 0 || rs1 != 0 || rd != 0 {
			trap.EnterWithoutTval(cpu.State, trap.IllegalIstruction)
		}

	case 0b_001: // fence.i
		if imm != 0 || rs1 != 0 || rd != 0 {
			trap.EnterWithoutTval(cpu.State, trap.IllegalIstruction)
		}

		cpu.Update.ICache.Clear()

	default:
		trap.EnterWithoutTval(cpu.State, trap.IllegalIstruction)
	}
}
