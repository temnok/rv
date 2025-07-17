package rv

import (
	"math"
	"math/big"
)

const upper32ones = -1 << 32

// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#_rounding_modes
var roundingModes = []big.RoundingMode{
	RmRNE: big.ToNearestEven,
	RmRTZ: big.ToZero,
	RmRDN: big.ToNegativeInf,
	RmRUP: big.ToPositiveInf,
	RmRMM: big.ToNearestAway,
}

func (cpu *CPU) matchRoundMode(rm int) big.RoundingMode {
	if rm > RmRMM {
		rm = bits(cpu.CSR.Fcsr, FcsrRM, 3)
	}

	return roundingModes[rm]
}

func (cpu *CPU) fp32(f *big.Float, i, rm int) {
	f.SetPrec(23)
	f.SetMode(cpu.matchRoundMode(rm))
	f.SetFloat64(float64(math.Float32frombits(uint32(cpu.FReg[i]))))
}

func (cpu *CPU) fp64(f *big.Float, i, rm int) {
	f.SetPrec(52)
	f.SetMode(cpu.matchRoundMode(rm))
	f.SetFloat64(math.Float64frombits(uint64(cpu.FReg[i])))
}

func (cpu *CPU) execComputeFP(f7, rs2, rs1, f3, rd int) {
	if f7&1 == 1 && !cpu.xlen64() || bits(cpu.CSR.Mstatus, MstatusFS, 2) == FSoff {
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	//switch f7 {
	//case 0b_0000000: // fadd.s
	//	a, b := cpu.fp32(rs1), cpu.fp32(rs2)
	//	c := a + b
	//	cpu.Updated.FRegValue = fp32bits(c)
	//
	//case 0b_0000001: // fadd.d
	//	a, b := cpu.fp64(rs1), cpu.fp64(rs2)
	//	c := a + b
	//	cpu.Updated.FRegValue = fp64bits(c)
	//
	//case 0b_0000100: // fsub.s
	//	a, b := cpu.fp32(rs1), cpu.fp32(rs2)
	//	c := a - b
	//	cpu.Updated.FRegValue = fp32bits(c)
	//
	//case 0b_0000101: // fsub.d
	//	a, b := cpu.fp64(rs1), cpu.fp64(rs2)
	//	c := a - b
	//	cpu.Updated.FRegValue = fp64bits(c)
	//
	//case 0b_0001000: // fmul.s
	//	a, b := cpu.fp32(rs1), cpu.fp32(rs2)
	//	c := a * b
	//	cpu.Updated.FRegValue = fp32bits(c)
	//
	//case 0b_0001001: // fmul.d
	//	a, b := cpu.fp64(rs1), cpu.fp64(rs2)
	//	c := a * b
	//	cpu.Updated.FRegValue = fp64bits(c)
	//
	//case 0b_0001100: // fdiv.s
	//	a, b := cpu.fp32(rs1), cpu.fp32(rs2)
	//	c := a / b
	//	cpu.Updated.FRegValue = fp32bits(c)
	//
	//case 0b_0001101: // fdiv.d
	//	a, b := cpu.fp64(rs1), cpu.fp64(rs2)
	//	c := a / b
	//	cpu.Updated.FRegValue = fp64bits(c)
	//
	//case 0b_1110000: // fmv.x.w
	//	cpu.Updated.RegIndex = rd
	//	cpu.Updated.RegValue = int(int32(cpu.FReg[rs1]))
	//	return
	//
	//case 0b_1110001: // fmv.x.d
	//	cpu.Updated.RegIndex = rd
	//	cpu.Updated.RegValue = cpu.FReg[rs1]
	//	return
	//
	//case 0b_1111000: // fmv.w.x
	//	cpu.Updated.FRegValue = upper32ones | cpu.Reg[rs1]
	//
	//case 0b_1111001: // fmv.d.x
	//	cpu.Updated.FRegValue = cpu.Reg[rs1]
	//
	//default:
	//	cpu.trap(ExceptionIllegalIstruction)
	//	return
	//}
	//
	//cpu.Updated.FRegUpdated = true
	//cpu.Updated.FRegIndex = rd
}
