package rv

import "math"

const upper32ones = -1 << 32

func (cpu *CPU) fp32(i int) float32 {
	return math.Float32frombits(uint32(cpu.FReg[i]))
}

func (cpu *CPU) fp64(i int) float64 {
	return math.Float64frombits(uint64(cpu.FReg[i]))
}

func fp32bits(val float32) int {
	return upper32ones | int(math.Float32bits(val))
}

func fp64bits(val float64) int {
	return int(math.Float64bits(val))
}

func (cpu *CPU) execComputeFP(f7, rs2, rs1, f3, rd int) {
	if f7&1 == 1 && !cpu.xlen64() || bits(cpu.CSR.Mstatus, MstatusFS, 2) == FSoff {
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	switch f7 {
	case 0b_0000000: // fadd.s
		cpu.Updated.FRegValue = fp32bits(cpu.fp32(rs1) + cpu.fp32(rs2))

	case 0b_0000001: // fadd.d
		cpu.Updated.FRegValue = fp64bits(cpu.fp64(rs1) + cpu.fp64(rs2))

	case 0b_0000100: // fsub.s
		cpu.Updated.FRegValue = fp32bits(cpu.fp32(rs1) - cpu.fp32(rs2))

	case 0b_0000101: // fsub.d
		cpu.Updated.FRegValue = fp64bits(cpu.fp64(rs1) - cpu.fp64(rs2))

	case 0b_0001000: // fmul.s
		cpu.Updated.FRegValue = fp32bits(cpu.fp32(rs1) * cpu.fp32(rs2))

	case 0b_0001001: // fmul.d
		cpu.Updated.FRegValue = fp64bits(cpu.fp64(rs1) * cpu.fp64(rs2))

	case 0b_0001100: // fdiv.s
		cpu.Updated.FRegValue = fp32bits(cpu.fp32(rs1) / cpu.fp32(rs2))

	case 0b_0001101: // fdiv.d
		cpu.Updated.FRegValue = fp64bits(cpu.fp64(rs1) / cpu.fp64(rs2))

	case 0b_1110000: // fmv.x.w
		cpu.Updated.RegIndex = rd
		cpu.Updated.RegValue = int(int32(cpu.FReg[rs1]))
		return

	case 0b_1110001: // fmv.x.d
		cpu.Updated.RegIndex = rd
		cpu.Updated.RegValue = cpu.FReg[rs1]
		return

	case 0b_1111000: // fmv.w.x
		cpu.Updated.FRegValue = upper32ones | cpu.Reg[rs1]

	case 0b_1111001: // fmv.d.x
		cpu.Updated.FRegValue = cpu.Reg[rs1]

	default:
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	cpu.Updated.FRegUpdated = true
	cpu.Updated.FRegIndex = rd
}
