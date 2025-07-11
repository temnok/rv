package rv

func immCB(instr int) int {
	a := bits(instr, 3, 2)
	b := bits(instr, 10, 2)
	c := bit(instr, 2)
	d := bits(instr, 5, 2)
	e := signBit(instr, 12)

	return e<<8 | d<<6 | c<<5 | b<<3 | a<<1
}

func immCI(instr int) int {
	a := bits(instr, 2, 5)
	b := signBit(instr, 12)

	return b<<5 | a
}

func immCI2(instr int) int {
	a := bits(instr, 4, 3)
	b := bit(instr, 12)
	c := bits(instr, 2, 2)

	return c<<6 | b<<5 | a<<2
}

func immCI4(instr int) int {
	a := bit(instr, 6)
	b := bit(instr, 2)
	c := bit(instr, 5)
	d := bits(instr, 3, 2)
	e := signBit(instr, 12)

	return e<<9 | d<<7 | c<<6 | b<<5 | a<<4
}

func immCIW(instr int) int {
	a := bit(instr, 6)
	b := bit(instr, 5)
	c := bits(instr, 11, 2)
	d := bits(instr, 7, 4)

	return d<<6 | c<<4 | b<<3 | a<<2
}

func immCJ(instr int) int {
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

func immCL(instr int) int {
	a := bit(instr, 6)
	b := bits(instr, 10, 3)
	c := bit(instr, 5)

	return c<<6 | b<<3 | a<<2
}

func immCSS(instr int) int {
	a := bits(instr, 9, 4)
	b := bits(instr, 7, 2)

	return b<<6 | a<<2
}
