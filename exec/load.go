package exec

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/instr"
	"github.com/temnok/rv/mem"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func Load(cpu *state.CPU, op instr.Op) {
	imm, rd, rs1 := imm.I(op.Code()), op.Rd(), op.Rs1()

	var val int

	switch op.F3() {
	case 0b_000: // lb
		mem.Read(cpu, cpu.X[rs1]+imm, &val, 1)
		cpu.Xset(rd, int(int8(val)))

	case 0b_001: // lh
		mem.Read(cpu, cpu.X[rs1]+imm, &val, 2)
		cpu.Xset(rd, int(int16(val)))

	case 0b_010: // lw
		mem.Read(cpu, cpu.X[rs1]+imm, &val, 4)
		cpu.Xset(rd, int(int32(val)))

	case 0b_011: // ld
		if !cpu.Xlen64() {
			trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
			return
		}

		mem.Read(cpu, cpu.X[rs1]+imm, &val, 8)
		cpu.Xset(rd, val)

	case 0b_100: // lbu
		mem.Read(cpu, cpu.X[rs1]+imm, &val, 1)
		cpu.Xset(rd, int(uint8(val)))

	case 0b_101: // lhu
		mem.Read(cpu, cpu.X[rs1]+imm, &val, 2)
		cpu.Xset(rd, int(uint16(val)))

	case 0b_110: // lwu
		if !cpu.Xlen64() {
			trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
			return
		}

		mem.Read(cpu, cpu.X[rs1]+imm, &val, 4)
		cpu.Xset(rd, int(uint32(val)))
	}

	if cpu.Update.XReg < 0 && !trap.IsEntered(cpu) {
		trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
	}
}
