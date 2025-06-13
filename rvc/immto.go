package rvc

// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#_immediate_encoding_variants

// 31                    20                                       0
// |a a a a a a a a a a a a a a a a a a a a a b b b b b b b b b b b| imm
// |a b b b b b b b b b b b|                                       | I-code
//
// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#_integer_register_immediate_instructions
func immToI(imm int32) int32 {
	return imm << 20
}

// 31          25                                   7   5         0
// |a a a a a a a a a a a a a a a a a a a a a b b b b b b|c c c c c| imm
// |a b b b b b b|                         |c c c c c|             | S-code
//
// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#ldst
func immToS(imm int32) int32 {
	ab := imm >> 5
	c := imm & 0b_11111

	return ab<<25 | c<<7
}

// 31          25                        12 11    8 7   5       1 0
// |a a a a a a a a a a a a a a a a a a a a|d|b b b b b b|c c c c|0| imm
// |a|b b b b b b|                         |c c c c|d|             | B-code
//
// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#_conditional_branches
func immToB(imm int32) int32 {
	a := imm >> 12
	b := (imm >> 5) & 0b_111111
	c := (imm >> 1) & 0b_1111
	d := (imm >> 11) & 1

	return a<<31 | b<<25 | c<<8 | d<<7
}

// 31                                    12                       0
// |a b b b b b b b b b b b b b b b b b b b|0 0 0 0 0 0 0 0 0 0 0 0| imm
// |a b b b b b b b b b b b b b b b b b b b|                       | U-code
//
// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#_integer_register_immediate_instructions
func immToU(imm int32) int32 {
	return imm &^ 0b_111111111111
}

// 31                    20              12 11                  1 0
// |a a a a a a a a a a a a|d d d d d d d d|c|b b b b b b b b b b|0| imm
// |a|b b b b b b b b b b|c|d d d d d d d d|                       | J-code
//
// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#_unconditional_jumps
func immToJ(imm int32) int32 {
	a := imm >> 20
	b := (imm >> 1) & 0b_1111111111
	c := (imm >> 11) & 1
	d := (imm >> 12) & 0b_11111111

	return a<<31 | b<<21 | c<<20 | d<<12
}
