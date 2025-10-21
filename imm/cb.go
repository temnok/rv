package imm

func CB(instr int) int {
	a := bits(instr, 3, 2)
	b := bits(instr, 10, 2)
	c := bit(instr, 2)
	d := bits(instr, 5, 2)
	e := signBit(instr, 12)

	return e<<8 | d<<6 | c<<5 | b<<3 | a<<1
}
