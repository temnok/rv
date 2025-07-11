package rv

func immCIW(instr int) int {
	a := bits(instr, 11, 2)
	b := bits(instr, 7, 4)
	c := bit(instr, 6)
	d := bit(instr, 5)
	return b<<6 | a<<4 | d<<3 | c<<2
}

func immCL(instr int) int {
	a := bits(instr, 10, 3)
	b := bit(instr, 6)
	c := bit(instr, 5)
	return c<<6 | a<<3 | b<<2
}

func immCLx8(instr int) int {
	a := bits(instr, 10, 3)
	b := bits(instr, 5, 2)
	return b<<6 | a<<3
}

func immCI(instr int) int {
	a := signedBit(instr, 12)
	b := bits(instr, 2, 5)
	return a<<5 | b
}

func immCI2(instr int) int {
	a := bit(instr, 12)
	b := bits(instr, 4, 3)
	c := bits(instr, 2, 2)
	return c<<6 | a<<5 | b<<2
}

func immCIx8(instr int) int {
	a := bit(instr, 12)
	b := bits(instr, 5, 2)
	c := bits(instr, 2, 3)
	return c<<6 | a<<5 | b<<3
}

func immCI4(instr int) int {
	a := signedBit(instr, 12)
	b := bit(instr, 6)
	c := bit(instr, 5)
	d := bits(instr, 3, 2)
	e := bit(instr, 2)
	return a<<9 | d<<7 | c<<6 | e<<5 | b<<4
}

func immCJ(instr int) int {
	a := signedBit(instr, 12)
	b := bit(instr, 11)
	c := bits(instr, 9, 2)
	d := bit(instr, 8)
	e := bit(instr, 7)
	f := bit(instr, 6)
	g := bits(instr, 3, 3)
	h := bit(instr, 2)
	return a<<11 | d<<10 | c<<8 | f<<7 | e<<6 | h<<5 | b<<4 | g<<1
}

func immCB(instr int) int {
	a := signedBit(instr, 12)
	b := bits(instr, 10, 2)
	c := bits(instr, 5, 2)
	d := bits(instr, 3, 2)
	e := bit(instr, 2)
	return a<<8 | c<<6 | e<<5 | b<<3 | d<<1
}

func immCSS(instr int) int {
	a := bits(instr, 9, 4)
	b := bits(instr, 7, 2)
	return b<<6 | a<<2
}

func immCSSx8(instr int) int {
	a := bits(instr, 10, 3)
	b := bits(instr, 7, 3)
	return b<<6 | a<<3
}
