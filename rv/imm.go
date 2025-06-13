package rv

// 31                    20                                       0
// |a b b b b b b b b b b b|                                       | I-code
// |a a a a a a a a a a a a a a a a a a a a a b b b b b b b b b b b| imm
//
// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#_integer_register_immediate_instructions
func immFromI(code int32) int32 {
	return code >> 20
}

// 31          25                                   7   5         0
// |a b b b b b b|                         |c c c c c|             | S-code
// |a a a a a a a a a a a a a a a a a a a a a b b b b b b|c c c c c| imm
//
// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#ldst
func immFromS(code int32) int32 {
	ab := code >> 25
	c := (code >> 7) & 0b_11111

	return ab<<5 | c
}

// 31          25                        12 11    8 7   5       1 0
// |a|b b b b b b|                         |c c c c|d|             | B-code
// |a a a a a a a a a a a a a a a a a a a a|d|b b b b b b|c c c c|0| imm
//
// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#_conditional_branches
func immFromB(code int32) int32 {
	a := code >> 31
	b := (code >> 25) & 0b_111111
	c := (code >> 8) & 0b_1111
	d := (code >> 7) & 1

	return a<<12 | d<<11 | b<<5 | c<<1
}

// 31                                    12                       0
// |a b b b b b b b b b b b b b b b b b b b|                       | U-code
// |a b b b b b b b b b b b b b b b b b b b|0 0 0 0 0 0 0 0 0 0 0 0| imm
//
// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#_integer_register_immediate_instructions
func immFromU(code int32) int32 {
	return code &^ 0b_111111111111
}

// 31                    20              12 11                  1 0
// |a|b b b b b b b b b b|c|d d d d d d d d|                       | J-code
// |a a a a a a a a a a a a|d d d d d d d d|c|b b b b b b b b b b|0| imm
//
// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#_unconditional_jumps
func immFromJ(code int32) int32 {
	a := code >> 31
	b := (code >> 21) & 0b_1111111111
	c := (code >> 20) & 1
	d := (code >> 12) & 0b_11111111

	return a<<20 | d<<12 | c<<11 | b<<1
}
