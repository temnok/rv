package imm

func CJ(instr int) int {
	a := bits(instr, 3, 3)
	b := bit(instr, 11)
	c := bit(instr, 2)
	d := bit(instr, 7)
	e := bit(instr, 6)
	f := bits(instr, 9, 2)
	g := bit(instr, 8)
	h := signBit(instr, 12)

	return h<<11 | g<<10 | f<<8 | e<<7 | d<<6 | c<<5 | b<<4 | a<<1
}
