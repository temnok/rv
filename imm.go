package rv

func immI(opcode int32) int32 {
	return opcode >> 20
}

func immS(opcode int32) int32 {
	a := opcode >> 25
	b := bits(opcode, 7, 5)
	return a<<5 | b
}

func immB(opcode int32) int32 {
	a := -bit(opcode, 31)
	b := bits(opcode, 25, 6)
	c := bits(opcode, 8, 4)
	d := bit(opcode, 7)
	return a<<12 | d<<11 | b<<5 | c<<1
}

func immU(opcode int32) int32 {
	return bits(opcode, 12, 20) << 12
}

func immJ(opcode int32) int32 {
	a := -bit(opcode, 31)
	b := bits(opcode, 21, 10)
	c := bit(opcode, 20)
	d := bits(opcode, 12, 8)
	return a<<20 | d<<12 | c<<11 | b<<1
}

func bit(val, i int32) int32 {
	return (val >> i) & 1
}

func bits(val, i, n int32) int32 {
	return (val >> i) & (1<<n - 1)
}
