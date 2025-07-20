package rv

import (
	"math"
	"testing"
)

// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#int_conv

func Test_fcvt_w(t *testing.T) {
	testF32toInt(t, "fcvt_w", fcvt_w_s, fcvt_w_d, map[float32]int{
		float32(math.MinInt32):     math.MinInt32,
		float32(math.MaxInt32):     math.MaxInt32,
		float32(math.MinInt32) * 2: math.MinInt32,
		float32(math.Inf(-1)):      math.MinInt32,
		float32(math.MaxInt32) * 2: math.MaxInt32,
		float32(math.Inf(1)):       math.MaxInt32,
		float32(math.NaN()):        math.MaxInt32,
	})
}

func Test_fcvt_wu(t *testing.T) {
	testF32toInt(t, "fcvt_wu", fcvt_wu_s, fcvt_wu_d, map[float32]int{
		float32(0):                  0,
		float32(math.MaxUint32):     math.MaxUint32,
		float32(-1):                 0,
		float32(math.Inf(-1)):       0,
		float32(math.MaxUint32) * 2: math.MaxUint32,
		float32(math.Inf(1)):        math.MaxUint32,
		float32(math.NaN()):         math.MaxUint32,
	})
}

func Test_fcvt_l(t *testing.T) {
	testF32toInt(t, "fcvt_l", fcvt_l_s, fcvt_l_d, map[float32]int{
		float32(math.MinInt64):     math.MinInt64,
		float32(math.MaxInt64):     math.MaxInt64,
		float32(math.MinInt64) * 2: math.MinInt64,
		float32(math.Inf(-1)):      math.MinInt64,
		float32(math.MaxInt64) * 2: math.MaxInt64,
		float32(math.Inf(1)):       math.MaxInt64,
		float32(math.NaN()):        math.MaxInt64,
	})
}

func Test_fcvt_lu(t *testing.T) {
	testF32toInt(t, "fcvt_lu", fcvt_lu_s, fcvt_lu_d, map[float32]int{
		float32(0):                    0,
		float32(math.MaxUint64):       -1,
		float32(math.MaxUint64+1) / 2: math.MinInt64,
		float32(-1):                   0,
		float32(math.Inf(-1)):         0,
		float32(math.MaxUint64) * 2:   -1,
		float32(math.Inf(1)):          -1,
		float32(math.NaN()):           -1,
	})
}

func testF32toInt(t *testing.T, fName string, f32 func(float32) int, f64 func(float64) int, tests map[float32]int) {
	for arg, want := range tests {
		if got := f32(arg); got != want {
			t.Errorf("%v_s(%v):\nwant %v\n got %v", fName, arg, want, got)
		}

		if got := f64(float64(arg)); got != want {
			t.Errorf("%v_d(%v):\nwant %v\n got %v", fName, arg, want, got)
		}
	}
}
