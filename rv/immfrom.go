package rv

// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#_immediate_encoding_variants

// 31                    20                11                     0
// |a|                    b|                                       | I-code
// |                                       -a|                    b| imm
//
// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#_integer_register_immediate_instructions
func ImmFromI(code int32) int32 {
	return code >> 20
}

// 31          25                        12 11      7   5         0
// |a|          b|                        x|        c|            x| S-code
// |                                       -a|          b|        c| imm
//
// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#ldst
func ImmFromS(code int32) int32 {
	ab := code >> 25
	c := (code >> 7) & 0b_11111

	return ab<<5 | c
}

// 31          25                        12 11    8 7   5       1 0
// |a|          b|                        x|      c|d|            x| B-code
// |                                     -a|d|          b|      c|0| imm
//
// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#_conditional_branches
func ImmFromB(code int32) int32 {
	a := code >> 31
	b := (code >> 25) & 0b_111111
	c := (code >> 8) & 0b_1111
	d := (code >> 7) & 1

	return a<<12 | d<<11 | b<<5 | c<<1
}

// 31                                    12                       0
// |                                      a|                      x| U-code
// |                                      a|                      0| imm
//
// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#_integer_register_immediate_instructions
func ImmFromU(code int32) int32 {
	return code &^ 0b_111111111111
}

// 31                    20              12 11                  1 0
// |a|                  b|c|              d|                      x| J-code
// |                     -a|              d|c|                  b|0| imm
//
// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#_unconditional_jumps
func ImmFromJ(code int32) int32 {
	a := code >> 31
	b := (code >> 21) & 0b_1111111111
	c := (code >> 20) & 1
	d := (code >> 12) & 0b_11111111

	return a<<20 | d<<12 | c<<11 | b<<1
}
