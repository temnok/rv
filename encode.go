package rv

func encodeR(f7, rs2, rs1, f3, rd, op int) int {
	return f7<<25 | rs2<<20 | rs1<<15 | f3<<12 | rd<<7 | op<<2 | 3
}

func encodeI(imm, rs1, f3, rd, op int) int {
	return imm<<20 | rs1<<15 | f3<<12 | rd<<7 | op<<2 | 3
}

func encodeS(imm, rs2, rs1, f3, op int) int {
	a := bits(imm, 5, 7)
	b := bits(imm, 0, 5)
	return a<<25 | rs2<<20 | rs1<<15 | f3<<12 | b<<7 | op<<2 | 3
}

func encodeB(imm, rs2, rs1, f3, op int) int {
	a := bit(imm, 12)
	b := bit(imm, 11)
	c := bits(imm, 5, 6)
	d := bits(imm, 1, 4)
	return a<<31 | c<<25 | rs2<<20 | rs1<<15 | f3<<12 | d<<8 | b<<7 | op<<2 | 3
}

func encodeU(imm, rd, op int) int {
	return imm<<12 | rd<<7 | op<<2 | 3
}

func encodeJ(imm, rd, op int) int {
	a := bit(imm, 20)
	b := bits(imm, 12, 8)
	c := bit(imm, 11)
	d := bits(imm, 1, 10)
	return a<<31 | d<<21 | c<<20 | b<<12 | rd<<7 | op<<2 | 3
}
