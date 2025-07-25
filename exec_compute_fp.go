package rv

/*
#cgo LDFLAGS: -lm

#include <fenv.h>
#include <math.h>
#include <stdint.h>

double nan32() { int32_t nan32bits = 0x7fc00000; return *(float*)&nan32bits; }
double nan64() { int64_t nan64bits = 0x7ff8000000000000; return *(double*)&nan64bits; }
float  fixNan32(float r, float a, float b) { return !isnan(r)? r : !isnan(a)? a : !isnan(b)? b : nan32(); }
double fixNan64(double r, double a, double b) { return !isnan(r)? r : !isnan(a)? a : !isnan(b)? b : nan64(); }

float  fadd_s(float a, float b)    { return a + b; }
double fadd_d(double a, double b)  { return a + b; }
float  fsub_s(float a, float b)    { return a - b; }
double fsub_d(double a, double b)  { return a - b; }
float  fmul_s(float a, float b)    { return a * b; }
double fmul_d(double a, double b)  { return a * b; }
float  fdiv_s(float a, float b)    { return a / b; }
double fdiv_d(double a, double b)  { return a / b; }

float  fsqrt_s(float a)            { return sqrtf(a); }
double fsqrt_d(double a)           { return sqrt(a); }

float  fmin_s(float a, float b)    { return a==0 && b==0? (copysign(1, a)<0? a : b) : fixNan32(fminf(a, b), a, b); }
float  fmax_s(float a, float b)    { return a==0 && b==0? (copysign(1, a)>0? a : b) : fixNan32(fmaxf(a, b), a, b); }
double fmin_d(double a, double b)  { return a==0 && b==0? (copysign(1, a)<0? a : b) : fixNan64(fmin(a, b), a, b); }
double fmax_d(double a, double b)  { return a==0 && b==0? (copysign(1, a)>0? a : b) : fixNan64(fmax(a, b), a, b); }

float  fmadd_s(float a, float b, float c)     { return a*b + c; }
float  fmsub_s(float a, float b, float c)     { return a*b - c; }
float  fnmadd_s(float a, float b, float c)    { return -a*b - c; }
float  fnmsub_s(float a, float b, float c)    { return -a*b + c; }
double fmadd_d(double a, double b, double c)  { return a*b + c; }
double fmsub_d(double a, double b, double c)  { return a*b - c; }
double fnmadd_d(double a, double b, double c) { return -a*b - c; }
double fnmsub_d(double a, double b, double c) { return -a*b + c; }

float    fcvt_s_d(double a)    { return isnan(a)? nan32() : (float)a; }
float    fcvt_s_l(int64_t a)   { return (float)a; }
float    fcvt_s_lu(uint64_t a) { return (float)a; }
float    fcvt_s_w(int32_t a)   { return (float)a; }
float    fcvt_s_wu(uint32_t a) { return (float)a; }
double   fcvt_d_s(float a)     { return isnan(a)? nan64() : (double)a; }
double   fcvt_d_l(int64_t a)   { return (double)a; }
double   fcvt_d_lu(uint64_t a) { return (double)a; }
double   fcvt_d_w(int32_t a)   { return (double)a; }
double   fcvt_d_wu(uint32_t a) { return (double)a; }

void ce() { feclearexcept(FE_INEXACT); feraiseexcept(FE_INVALID); }

int32_t  fcvt_w_s(float a)     { return isnan(a)? INT32_MAX :
	(a = rintf(a)) <= INT32_MIN? ce(), INT32_MIN : a >= INT32_MAX? ce(), INT32_MAX : (int32_t)a; }
uint32_t fcvt_wu_s(float a)    { return isnan(a)? UINT32_MAX :
	(a = rintf(a)) < 0? ce(), 0 : a >= UINT32_MAX? ce(), UINT32_MAX : (uint32_t)a; }
int64_t  fcvt_l_s(float a)     { return isnan(a)? INT64_MAX :
	(a = rint(a)) <= INT64_MIN? ce(), INT64_MIN : a >= INT64_MAX? ce(), INT64_MAX : (int64_t)a; }
uint64_t fcvt_lu_s(float a)    { return isnan(a)? UINT64_MAX :
	(a = rint(a)) < 0? ce(), 0 : a >= UINT64_MAX? ce(), UINT64_MAX : (uint64_t)a; }
int32_t  fcvt_w_d(double a)    { return isnan(a)? INT32_MAX :
	(a = rintf(a)) <= INT32_MIN? ce(), INT32_MIN : a >= INT32_MAX? ce(), INT32_MAX : (int32_t)a; }
uint32_t fcvt_wu_d(double a)   { return isnan(a)? UINT32_MAX :
	(a = rintf(a)) < 0? ce(), 0 : a >= UINT32_MAX? ce(), UINT32_MAX : (uint32_t)a; }
int64_t  fcvt_l_d(double a)    { return isnan(a)? INT64_MAX :
	(a = rint(a)) <= INT64_MIN? ce(), INT64_MIN : a >= INT64_MAX? ce(), INT64_MAX : (int64_t)a; }
uint64_t fcvt_lu_d(double a)   { return isnan(a)? UINT64_MAX :
	(a = rint(a)) < 0? ce(), 0 : a >= UINT64_MAX? ce(), UINT64_MAX : (uint64_t)a; }

*/
import "C"
import (
	"math"
)

const (
	f32boxingBits = -1 << 32
	f32signMask   = 1 << 31
	f64signMask   = -1 << 63

	nan32bits = -1<<32 | 0x7fc00000

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

func (cpu *CPU) execComputeFP(f7, rs2, rs1, f3, rd, op int) {
	if f7&1 == 1 && !cpu.extD() || cpu.fpDisabled() {
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	rs3 := 0
	if op != 0 {
		rs3 = f7 >> 2
		f7 = 0b_10000000 | f7&3
	}

	switch f7 {
	case 0b_0000000: // fadd.s
		cpu.f32set(rd, C.fadd_s(cpu.f32arg2(rs1, rs2, f3)))

	case 0b_0000001: // fadd.d
		cpu.f64set(rd, C.fadd_d(cpu.f64arg2(rs1, rs2, f3)))

	case 0b_0000100: // fsub.s
		cpu.f32set(rd, C.fsub_s(cpu.f32arg2(rs1, rs2, f3)))

	case 0b_0000101: // fsub.d
		cpu.f64set(rd, C.fsub_d(cpu.f64arg2(rs1, rs2, f3)))

	case 0b_0001000: // fmul.s
		cpu.f32set(rd, C.fmul_s(cpu.f32arg2(rs1, rs2, f3)))

	case 0b_0001001: // fmul.d
		cpu.f64set(rd, C.fmul_d(cpu.f64arg2(rs1, rs2, f3)))

	case 0b_0001100: // fdiv.s
		cpu.f32set(rd, C.fdiv_s(cpu.f32arg2(rs1, rs2, f3)))

	case 0b_0001101: // fdiv.d
		cpu.f64set(rd, C.fdiv_d(cpu.f64arg2(rs1, rs2, f3)))

	case 0b_0010000:
		switch f3 {
		case 0b_000: // fsgnj.s
			cpu.f32setBits(rd, cpu.f32bits(rs1)&^f32signMask|cpu.f32bits(rs2)&f32signMask)

		case 0b_001: // fsgnjn.s
			cpu.f32setBits(rd, cpu.f32bits(rs1)&^f32signMask|(^cpu.f32bits(rs2))&f32signMask)

		case 0b_010: // fsgnjx.s
			cpu.f32setBits(rd, cpu.f32bits(rs1)^cpu.f32bits(rs2)&f32signMask)
		}

	case 0b_0010001:
		switch f3 {
		case 0b_000: // fsgnj.d
			cpu.f64setBits(rd, cpu.F[rs1]&^f64signMask|cpu.F[rs2]&f64signMask)

		case 0b_001: // fsgnjn.d
			cpu.f64setBits(rd, cpu.F[rs1]&^f64signMask|(^cpu.F[rs2])&f64signMask)

		case 0b_010: // fsgnjx.d
			cpu.f64setBits(rd, cpu.F[rs1]^cpu.F[rs2]&f64signMask)
		}

	case 0b_0010100:
		switch f3 {
		case 0b_000: // fmin.s
			cpu.f32set(rd, C.fmin_s(cpu.f32arg2(rs1, rs2, -1)))

		case 0b_001: // fmax.s
			cpu.f32set(rd, C.fmax_s(cpu.f32arg2(rs1, rs2, -1)))
		}

	case 0b_0010101:
		switch f3 {
		case 0b_000: // fmin.d
			cpu.f64set(rd, C.fmin_d(cpu.f64arg2(rs1, rs2, -1)))

		case 0b_001: // fmax.d
			cpu.f64set(rd, C.fmax_d(cpu.f64arg2(rs1, rs2, -1)))
		}

	case 0b_0100000:
		switch rs2 {
		case 0b_00001: // fcvt.s.d
			cpu.f32set(rd, C.fcvt_s_d(cpu.f64arg(rs1, f3)))
		}

	case 0b_0100001:
		switch rs2 {
		case 0b_00000: // fcvt.d.s
			cpu.f64set(rd, C.fcvt_d_s(cpu.f32arg(rs1, f3)))
		}

	case 0b_0101100: // fsqrt.s
		if rs2 == 0 {
			cpu.f32set(rd, C.fsqrt_s(cpu.f32arg(rs1, f3)))
		}

	case 0b_0101101: // fsqrt.d
		if rs2 == 0 {
			cpu.f64set(rd, C.fsqrt_d(cpu.f64arg(rs1, f3)))
		}

	case 0b_1010000: // https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#_single_precision_floating_point_compare_instructions
		switch f3 {
		case 0b_000: // fle.s
			a, b := cpu.f32(rs1), cpu.f32(rs2)
			cpu.fsetCmp(rd, a <= b, isNaN32(a) || isNaN32(b))

		case 0b_001: // flt.s
			a, b := cpu.f32(rs1), cpu.f32(rs2)
			cpu.fsetCmp(rd, a < b, isNaN32(a) || isNaN32(b))

		case 0b_010: // feq.s
			a, b := cpu.f32(rs1), cpu.f32(rs2)
			cpu.fsetCmp(rd, a == b, isSNaN32(a) || isSNaN32(b))
		}

	case 0b_1010001:
		switch f3 {
		case 0b_000: // fle.d
			a, b := cpu.f64(rs1), cpu.f64(rs2)
			cpu.fsetCmp(rd, a <= b, isNaN64(a) || isNaN64(b))

		case 0b_001: // flt.d
			a, b := cpu.f64(rs1), cpu.f64(rs2)
			cpu.fsetCmp(rd, a < b, isNaN64(a) || isNaN64(b))

		case 0b_010: // feq.d
			a, b := cpu.f64(rs1), cpu.f64(rs2)
			cpu.fsetCmp(rd, a == b, isSNaN64(a) || isSNaN64(b))
		}

	// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#_single_precision_floating_point_conversion_and_move_instructions
	case 0b_1100000:
		switch rs2 {
		case 0b_00000: // fcvt.w.s
			cpu.fxset(rd, int(C.fcvt_w_s(cpu.f32arg(rs1, f3))))

		case 0b_00001: // fcvt.wu.s
			cpu.fxset(rd, int(int32(C.fcvt_wu_s(cpu.f32arg(rs1, f3)))))

		case 0b_00010: // fcvt.l.s
			cpu.fxset(rd, int(C.fcvt_l_s(cpu.f32arg(rs1, f3))))

		case 0b_00011: // fcvt.lu.s
			cpu.fxset(rd, int(C.fcvt_lu_s(cpu.f32arg(rs1, f3))))
		}

	case 0b_1100001:
		switch rs2 {
		case 0b_00000: // fcvt.w.d
			cpu.fxset(rd, int(C.fcvt_w_d(cpu.f64arg(rs1, f3))))

		case 0b_00001: // fcvt.wu.d
			cpu.fxset(rd, int(int32(C.fcvt_wu_d(cpu.f64arg(rs1, f3)))))

		case 0b_00010: // fcvt.l.d
			cpu.fxset(rd, int(C.fcvt_l_d(cpu.f64arg(rs1, f3))))

		case 0b_00011: // fcvt.lu.d
			cpu.fxset(rd, int(C.fcvt_lu_d(cpu.f64arg(rs1, f3))))
		}

	case 0b_1101000:
		switch rs2 {
		case 0b_00000: // fcvt.s.w
			cpu.f32set(rd, C.fcvt_s_w(C.int32_t(cpu.fxget(rs1, f3))))

		case 0b_00001: // fcvt.s.wu
			cpu.f32set(rd, C.fcvt_s_wu(C.uint32_t(cpu.fxget(rs1, f3))))

		case 0b_00010: // fcvt.s.l
			cpu.f32set(rd, C.fcvt_s_l(C.int64_t(cpu.fxget(rs1, f3))))

		case 0b_00011: // fcvt.s.lu
			cpu.f32set(rd, C.fcvt_s_lu(C.uint64_t(cpu.fxget(rs1, f3))))
		}

	case 0b_1101001:
		switch rs2 {
		case 0b_00000: // fcvt.d.w
			cpu.f64set(rd, C.fcvt_d_w(C.int32_t(cpu.fxget(rs1, f3))))

		case 0b_00001: // fcvt.d.wu
			cpu.f64set(rd, C.fcvt_d_wu(C.uint32_t(cpu.fxget(rs1, f3))))

		case 0b_00010: // fcvt.d.l
			cpu.f64set(rd, C.fcvt_d_l(C.int64_t(cpu.fxget(rs1, f3))))

		case 0b_00011: // fcvt.d.lu
			cpu.f64set(rd, C.fcvt_d_lu(C.uint64_t(cpu.fxget(rs1, f3))))
		}

	case 0b_1110000:
		switch rs2<<3 | f3 {
		case 0b_00000_000: // fmv.x.w
			cpu.xset(rd, int(int32(cpu.F[rs1])))

		case 0b_00000_001: // fclass.s
			cpu.xset(rd, fclass_s(cpu.f32(rs1)))
		}

	case 0b_1110001:
		switch rs2<<3 | f3 {
		case 0b_00000_000: // fmv.x.d
			cpu.xset(rd, cpu.F[rs1])

		case 0b_00000_001: // fclass.d
			cpu.xset(rd, fclass_d(cpu.f64(rs1)))
		}

	case 0b_1111000: // fmv.w.x
		cpu.f32setBits(rd, cpu.X[rs1])

	case 0b_1111001: // fmv.d.x
		cpu.f64setBits(rd, cpu.X[rs1])

	case 0b_10000000:
		switch op {
		case 0b_10000: // fmadd.s
			cpu.f32set(rd, C.fmadd_s(cpu.f32arg3(rs1, rs2, rs3, f3)))

		case 0b_10001: // fmsub.s
			cpu.f32set(rd, C.fmsub_s(cpu.f32arg3(rs1, rs2, rs3, f3)))

		case 0b_10010: // fnmsub.s
			cpu.f32set(rd, C.fnmsub_s(cpu.f32arg3(rs1, rs2, rs3, f3)))

		case 0b_10011: // fnmadd.s
			cpu.f32set(rd, C.fnmadd_s(cpu.f32arg3(rs1, rs2, rs3, f3)))
		}

	case 0b_10000001:
		switch op {
		case 0b_10000: // fmadd.d
			cpu.f64set(rd, C.fmadd_d(cpu.f64arg3(rs1, rs2, rs3, f3)))

		case 0b_10001: // fmsub.d
			cpu.f64set(rd, C.fmsub_d(cpu.f64arg3(rs1, rs2, rs3, f3)))

		case 0b_10010: // fnmsub.d
			cpu.f64set(rd, C.fnmsub_d(cpu.f64arg3(rs1, rs2, rs3, f3)))

		case 0b_10011: // fnmadd.d
			cpu.f64set(rd, C.fnmadd_d(cpu.f64arg3(rs1, rs2, rs3, f3)))
		}
	}

	if cpu.Updated.XReg < 0 && cpu.Updated.FReg < 0 {
		cpu.trap(ExceptionIllegalIstruction)
	}
}

func (cpu *CPU) f32arg(rs1, f3 int) C.float {
	cpu.prepareCfenv(f3)
	return C.float(cpu.f32(rs1))
}

func (cpu *CPU) f32arg2(rs1, rs2, f3 int) (C.float, C.float) {
	cpu.prepareCfenv(f3)
	return C.float(cpu.f32(rs1)), C.float(cpu.f32(rs2))
}

func (cpu *CPU) f32arg3(rs1, rs2, rs3, f3 int) (C.float, C.float, C.float) {
	cpu.prepareCfenv(f3)
	return C.float(cpu.f32(rs1)), C.float(cpu.f32(rs2)), C.float(cpu.f32(rs3))
}

func (cpu *CPU) f64arg(rs1, f3 int) C.double {
	cpu.prepareCfenv(f3)
	return C.double(cpu.f64(rs1))
}

func (cpu *CPU) f64arg2(rs1, rs2, f3 int) (C.double, C.double) {
	cpu.prepareCfenv(f3)
	return C.double(cpu.f64(rs1)), C.double(cpu.f64(rs2))
}

func (cpu *CPU) f64arg3(rs1, rs2, rs3, f3 int) (C.double, C.double, C.double) {
	cpu.prepareCfenv(f3)
	return C.double(cpu.f64(rs1)), C.double(cpu.f64(rs2)), C.double(cpu.f64(rs3))
}

func (cpu *CPU) f32set(rd int, res C.float) {
	cpu.Updated.FReg = rd
	cpu.Updated.FVal = f32boxingBits | int(math.Float32bits(float32(res)))
	if uint(cpu.Updated.FVal) == 0xffffffffffc00000 {
		cpu.Updated.FVal &^= f32signMask
	}

	cpu.setUpdatedFflags()
}

func (cpu *CPU) f64set(rd int, res C.double) {
	cpu.Updated.FReg = rd
	cpu.Updated.FVal = int(math.Float64bits(float64(res)))
	if uint(cpu.Updated.FVal) == 0xfff8000000000000 {
		cpu.Updated.FVal &^= f64signMask
	}

	cpu.setUpdatedFflags()
}

func (cpu *CPU) fxget(rs1, f3 int) int {
	cpu.prepareCfenv(f3)
	return cpu.X[rs1]
}

func (cpu *CPU) fxset(rd, res int) {
	cpu.xset(rd, res)
	cpu.setUpdatedFflags()
}

func (cpu *CPU) fsetCmp(rd int, res, nv bool) {
	if res {
		cpu.xset(rd, 1)
	} else {
		cpu.xset(rd, 0)
	}
	if nv {
		cpu.Updated.Fflags |= 1 << FflagsNV
	}
}

func (cpu *CPU) f32setBits(rd, bits int) {
	cpu.Updated.FReg = rd
	cpu.Updated.FVal = f32boxingBits | bits
}

func (cpu *CPU) f64setBits(rd, bits int) {
	cpu.Updated.FReg = rd
	cpu.Updated.FVal = bits
}

func (cpu *CPU) f32(i int) float32 {
	return math.Float32frombits(uint32(cpu.f32bits(i)))
}

// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#nanboxing
func (cpu *CPU) f32bits(i int) int {
	val := cpu.F[i]
	if val>>32 != -1 {
		return nan32bits
	}
	return val
}

func (cpu *CPU) f64(i int) float64 {
	return math.Float64frombits(uint64(cpu.F[i]))
}

func (cpu *CPU) prepareCfenv(rm int) {
	if rm >= 0 {
		if rm == RmDYN {
			rm = bits(cpu.CSR.Fcsr, FcsrRM, 3)
		}

		C.fesetround(rmToC[rm])
	}

	C.feclearexcept(C.FE_ALL_EXCEPT)
}

func (cpu *CPU) setUpdatedFflags() {
	ex := C.fetestexcept(C.FE_ALL_EXCEPT)

	if ex&C.FE_INEXACT != 0 {
		cpu.Updated.Fflags |= 1 << FflagsNX
	}

	if ex&C.FE_UNDERFLOW != 0 {
		cpu.Updated.Fflags |= 1 << FflagsUF
	}

	if ex&C.FE_OVERFLOW != 0 {
		cpu.Updated.Fflags |= 1 << FflagsOF
	}

	if ex&C.FE_DIVBYZERO != 0 {
		cpu.Updated.Fflags |= 1 << FflagsDZ
	}

	if ex&C.FE_INVALID != 0 {
		cpu.Updated.Fflags |= 1 << FflagsNV
	}
}

func isNaN32(a float32) bool {
	return a != a
}

func isNaN64(a float64) bool {
	return a != a
}

func isSNaN32(a float32) bool {
	return a != a && math.Float32bits(a)&(1<<22) == 0
}

func isSNaN64(a float64) bool {
	return a != a && math.Float64bits(a)&(1<<51) == 0
}

func fclass_s(a float32) int {
	const smallestNormal = 1.1754943508222875e-38 // 2**-126

	i := 0

	switch {
	case math.IsInf(float64(a), -1):
		i = 0

	case a <= -smallestNormal:
		i = 1

	case a < 0:
		i = 2

	case math.Signbit(float64(a)):
		i = 3

	case a == 0:
		i = 4

	case a < smallestNormal:
		i = 5

	case float64(a) < math.Inf(1):
		i = 6

	case a > 0:
		i = 7

	case isSNaN32(a):
		i = 8

	default:
		i = 9
	}

	return 1 << i
}

func fclass_d(a float64) int {
	const smallestNormal = 2.2250738585072014e-308 // 2**-1022

	i := 0

	switch {
	case math.IsInf(a, -1):
		i = 0

	case a <= -smallestNormal:
		i = 1

	case a < 0:
		i = 2

	case math.Signbit(a):
		i = 3

	case a == 0:
		i = 4

	case a < smallestNormal:
		i = 5

	case a < math.Inf(1):
		i = 6

	case a > 0:
		i = 7

	case isSNaN64(a):
		i = 8

	default:
		i = 9
	}

	return 1 << i
}

func fcvt_l_d(a float64) int {
	return int(C.fcvt_l_d(C.double(a)))
}

func fcvt_l_s(a float32) int {
	return int(C.fcvt_l_s(C.float(a)))
}

func fcvt_lu_d(a float64) int {
	return int(C.fcvt_lu_d(C.double(a)))
}

func fcvt_lu_s(a float32) int {
	return int(C.fcvt_lu_s(C.float(a)))
}

func fcvt_w_d(a float64) int {
	return int(C.fcvt_w_d(C.double(a)))
}

func fcvt_w_s(a float32) int {
	return int(C.fcvt_w_s(C.float(a)))
}

func fcvt_wu_d(a float64) int {
	return int(C.fcvt_wu_d(C.double(a)))
}

func fcvt_wu_s(a float32) int {
	return int(C.fcvt_wu_s(C.float(a)))
}
