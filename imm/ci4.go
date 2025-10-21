package imm

func CI4(instr int) int {
	a := bit(instr, 6)
	b := bit(instr, 2)
	c := bit(instr, 5)
	d := bits(instr, 3, 2)
	e := signBit(instr, 12)

	return e<<9 | d<<7 | c<<6 | b<<5 | a<<4
}
