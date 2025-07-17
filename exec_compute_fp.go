package rv

/*
#include <fenv.h>

float  add32(float a, float b)   { return a + b; }
double add64(double a, double b) { return a + b; }
float  sub32(float a, float b)   { return a - b; }
double sub64(double a, double b) { return a - b; }
float  mul32(float a, float b)   { return a * b; }
double mul64(double a, double b) { return a * b; }
float  div32(float a, float b)   { return a / b; }
double div64(double a, double b) { return a / b; }

*/
import "C"
import (
	"math"
)

const (
	f32boxingBits = -1 << 32

	// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#rm
	RmRNE = 0b_000
	RmRTZ = 0b_001
	RmRDN = 0b_010
	RmRUP = 0b_011
	RmRMM = 0b_100
	RmDYN = 0b_111
)

var rmToC = []C.int{
	RmRNE: C.FE_TONEAREST,
	RmRTZ: C.FE_TOWARDZERO,
	RmRDN: C.FE_DOWNWARD,
	RmRUP: C.FE_UPWARD,
	RmRMM: C.FE_TONEAREST,
	RmDYN: C.FE_TONEAREST,
}

func (cpu *CPU) execComputeFP(f7, rs2, rs1, f3, rd int) {
	if f7&1 == 1 && !cpu.xlen64() || bits(cpu.CSR.Mstatus, MstatusFS, 2) == FSoff {
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	switch f7 {
	case 0b_0000000: // fadd.s
		cpu.f32res(C.add32(cpu.f32arg(rs1, rs2, f3)))

	case 0b_0000001: // fadd.d
		cpu.f64res(C.add64(cpu.f64arg(rs1, rs2, f3)))

	case 0b_0000100: // fsub.s
		cpu.f32res(C.sub32(cpu.f32arg(rs1, rs2, f3)))

	case 0b_0000101: // fsub.d
		cpu.f64res(C.sub64(cpu.f64arg(rs1, rs2, f3)))

	case 0b_0001000: // fmul.s
		cpu.f32res(C.mul32(cpu.f32arg(rs1, rs2, f3)))

	case 0b_0001001: // fmul.d
		cpu.f64res(C.mul64(cpu.f64arg(rs1, rs2, f3)))

	case 0b_0001100: // fdiv.s
		cpu.f32res(C.div32(cpu.f32arg(rs1, rs2, f3)))

	case 0b_0001101: // fdiv.d
		cpu.f64res(C.div64(cpu.f64arg(rs1, rs2, f3)))

	case 0b_1110000: // fmv.x.w
		cpu.Updated.RegIndex = rd
		cpu.Updated.RegValue = int(int32(cpu.FReg[rs1]))
		return

	case 0b_1110001: // fmv.x.d
		cpu.Updated.RegIndex = rd
		cpu.Updated.RegValue = cpu.FReg[rs1]
		return

	case 0b_1111000: // fmv.w.x
		cpu.Updated.FRegValue = f32boxingBits | cpu.Reg[rs1]

	case 0b_1111001: // fmv.d.x
		cpu.Updated.FRegValue = cpu.Reg[rs1]

	default:
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	cpu.Updated.FRegUpdated = true
	cpu.Updated.FRegIndex = rd
	cpu.setUpdatedFflags()
}

func (cpu *CPU) f32arg(rs1, rs2, f3 int) (C.float, C.float) {
	cpu.prepareCfenv(f3)
	return C.float(cpu.f32(rs1)), C.float(cpu.f32(rs2))
}

func (cpu *CPU) f64arg(rs1, rs2, f3 int) (C.double, C.double) {
	cpu.prepareCfenv(f3)
	return C.double(cpu.f64(rs1)), C.double(cpu.f64(rs2))
}

func (cpu *CPU) f32res(res C.float) {
	cpu.Updated.FRegValue = f32bits(float32(res))
}

func (cpu *CPU) f64res(res C.double) {
	cpu.Updated.FRegValue = f64bits(float64(res))
}

func (cpu *CPU) f32(i int) float32 {
	return math.Float32frombits(uint32(cpu.FReg[i]))
}

func (cpu *CPU) f64(i int) float64 {
	return math.Float64frombits(uint64(cpu.FReg[i]))
}

func f32bits(val float32) int {
	return f32boxingBits | int(math.Float32bits(val))
}

func f64bits(val float64) int {
	return int(math.Float64bits(val))
}

func (cpu *CPU) prepareCfenv(rm int) {
	if rm == RmDYN {
		rm = bits(cpu.CSR.Fcsr, FcsrRM, 3)
	}

	C.fesetround(rmToC[rm])
	C.feclearexcept(C.FE_ALL_EXCEPT)
}

func (cpu *CPU) setUpdatedFflags() {
	ex := C.fetestexcept(C.FE_ALL_EXCEPT)

	if ex&C.FE_INEXACT != 0 {
		cpu.Updated.Fflags |= 1 << FcsrNX
	}

	if ex&C.FE_UNDERFLOW != 0 {
		cpu.Updated.Fflags |= 1 << FcsrUF
	}

	if ex&C.FE_OVERFLOW != 0 {
		cpu.Updated.Fflags |= 1 << FcsrOF
	}

	if ex&C.FE_DIVBYZERO != 0 {
		cpu.Updated.Fflags |= 1 << FcsrDZ
	}

	if ex&C.FE_INVALID != 0 {
		cpu.Updated.Fflags |= 1 << FcsrNV
	}
}
