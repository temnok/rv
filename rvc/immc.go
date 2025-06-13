package rvc

// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#_compressed_instruction_formats

// 31	                                 12           6 5     2   0
// |  	                                 |a|         |b b b b b|   | CI-code
// |a a a a a a a a a a a a a a a a a a a a a a a a a a a|b b b b b| imm
//
// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#_integer_constant_generation_instructions
func immFromCI(code int32) int32 {
	a := (code >> 12) & 1
	b := (code >> 2) & 0b_11111

	return -a<<5 | b
}

// 31                                    12           6 5 4   2   0
// |  	                                 |a|         |b b b|c c|   | CI-code
// |0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0|c c|a|b b b|0 0| imm
//
// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#_stack_pointer_based_loads_and_stores
func immFromCIx4(code int32) int32 {
	a := (code >> 12) & 1
	b := (code >> 4) & 0b_111
	c := (code >> 2) & 0b_11

	return c<<6 | a<<5 | b<<2
}

// 31                                    12     9   7 6 5 4 3 2   0
// |  	                                 |a|         |b|c|d d|e|   | CI-code
// |a a a a a a a a a a a a a a a a a a a a a a a|d d|c|e|b|0 0 0 0| imm
//
// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#_integer_register_immediate_operations
func immFromCIx16(code int32) int32 {
	a := (code >> 12) & 1
	b := (code >> 6) & 1
	c := (code >> 5) & 1
	d := (code >> 3) & 0b_11
	e := (code >> 2) & 1

	return -a<<9 | d<<7 | c<<6 | e<<5 | b<<4
}

// 31                                    12     9   7 6       2   0
// |  	                                 |a a a a|b b|             | CI-code
// |0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0|b b|a a a a|0 0| imm
//
// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#_stack_pointer_based_loads_and_stores
func immFromCSS(code int32) int32 {
	a := (code >> 9) & 0b_1111
	b := (code >> 7) & 0b_11

	return b<<6 | a<<2
}

// 31                                      11       7 6 5 4 3 2   0
// |  	                                 |a a|b b b b|c|d|         | CIW-code
// |0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0|b b b b|a a|d|c|0 0| imm
//
// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#_integer_register_immediate_operations
func immFromCIW(code int32) int32 {
	a := (code >> 11) & 0b_11
	b := (code >> 7) & 0b_1111
	c := (code >> 6) & 1
	d := (code >> 5) & 1

	return b<<6 | a<<4 | d<<3 | c<<2
}

// 31                                        10       6 5   3 2   0
// |  	                                 |a a a|     |b|c|         | CL-code
// |0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0|c|a a a|b|0 0| imm
//
// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#_register_based_loads_and_stores
func immFromCL(code int32) int32 {
	a := (code >> 10) & 0b_111
	b := (code >> 6) & 1
	c := (code >> 5) & 1

	return c<<6 | a<<3 | b<<2
}

// 31                                    12  10   8   6 5   3 2 1 0
// |  	                                 |a|b b|     |c c|d d|e|   | CB-code
// |a a a a a a a a a a a a a a a a a a a a a a a a|c c|e|b b|d d|0| imm
//
// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#_control_transfer_instructions_2
func immFromCB(code int32) int32 {
	a := (code >> 12) & 1
	b := (code >> 10) & 0b_11
	c := (code >> 5) & 0b_11
	d := (code >> 3) & 0b_11
	e := (code >> 2) & 1

	return -a<<8 | c<<6 | e<<5 | b<<3 | d<<1
}

// 31                                    12 11  9 8 7 6 5 4 3 2 1 0
// |  	                                 |a|b|c c|d|e|f|g g g|h|   | CJ-code
// |a a a a a a a a a a a a a a a a a a a a a|d|c c|f|e|h|b|g g g|0| imm
//
// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#_control_transfer_instructions_2
func immFromCJ(code int32) int32 {
	a := (code >> 12) & 1
	b := (code >> 11) & 1
	c := (code >> 9) & 0b_11
	d := (code >> 8) & 1
	e := (code >> 7) & 1
	f := (code >> 6) & 1
	g := (code >> 3) & 0b_111
	h := (code >> 2) & 1

	return -a<<11 | d<<10 | c<<8 | f<<7 | e<<6 | h<<5 | b<<4 | g<<1
}
