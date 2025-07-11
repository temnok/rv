package rv

func immCI3(instr int) int {
	a := bits(instr, 5, 2)
	b := bit(instr, 12)
	c := bits(instr, 2, 3)

	return c<<6 | b<<5 | a<<3
}

func immCL3(instr int) int {
	a := bits(instr, 10, 3)
	b := bits(instr, 5, 2)

	return b<<6 | a<<3
}

func immCSS3(instr int) int {
	a := bits(instr, 10, 3)
	b := bits(instr, 7, 3)

	return b<<6 | a<<3
}
