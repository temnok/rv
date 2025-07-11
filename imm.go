package rv

// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#immtypes

func immB(opcode int) int {
	a := bits(opcode, 8, 4)
	b := bits(opcode, 25, 6)
	c := bit(opcode, 7)
	d := signBit(opcode, 31)

	return d<<12 | c<<11 | b<<5 | a<<1
}

func immI(opcode int) int {
	return int(int32(opcode) >> 20)
}

func immJ(opcode int) int {
	a := bits(opcode, 21, 10)
	b := bit(opcode, 20)
	c := bits(opcode, 12, 8)
	d := signBit(opcode, 31)

	return d<<20 | c<<12 | b<<11 | a<<1
}

func immS(opcode int) int {
	a := bits(opcode, 7, 5)
	b := int(int32(opcode)) >> 25

	return b<<5 | a
}

func immU(opcode int) int {
	return int(int32(bits(opcode, 12, 20) << 12))
}
