package rv

import "fmt"

func Disasm(pc, opcode int32) string {
	buf := fmt.Sprintf("%08x\t", uint32(pc))

	if opcode&3 != 3 {
		opcode = int32(uint16(opcode))

		buf += fmt.Sprintf("%04x\t", opcode)

		decompressedOpcode := decompress(opcode)
		if decompressedOpcode == 0 {
			return buf + fmt.Sprintf(".half\t%x", int16(opcode))
		}

		opcode = decompressedOpcode
	} else {
		buf += fmt.Sprintf("%08x\t", uint32(opcode))
	}

	f7, rs2, rs1 := bits(opcode, 25, 7), bits(opcode, 20, 5), bits(opcode, 15, 5)
	f3, rd, op := bits(opcode, 12, 3), bits(opcode, 7, 5), bits(opcode, 2, 5)

	typeR := func(name string) string {
		return fmt.Sprintf("%v%v\t%v,%v,%v", buf, name, abiReg[rd], abiReg[rs1], abiReg[rs2])
	}

	if op == 0b_01100 {
		switch f7<<8 | f3<<5 {
		case 0b_0000000_000:
			return typeR("add")
		case 0b_0100000_000:
			return typeR("sub")
		case 0b_0000000_001:
			return typeR("sll")
		case 0b_0000000_010:
			return typeR("slt")
		case 0b_0000000_011:
			return typeR("sltu")
		case 0b_0000000_100:
			return typeR("xor")
		case 0b_0000000_101:
			return typeR("srl")
		case 0b_0100000_101:
			return typeR("sra")
		case 0b_0000000_110:
			return typeR("or")
		case 0b_0000000_111:
			return typeR("and")

		case 0b_0000001_000:
			return typeR("mul")
		case 0b_0000001_001:
			return typeR("mulh")
		case 0b_0000001_010:
			return typeR("mulhsu")
		case 0b_0000001_011:
			return typeR("mulhu")
		case 0b_0000001_100:
			return typeR("div")
		case 0b_0000001_101:
			return typeR("divu")
		case 0b_0000001_110:
			return typeR("rem")
		case 0b_0000001_111:
			return typeR("remu")
		}
	}

	return fmt.Sprintf("%v.word\t%x", buf, opcode)
}

var abiReg = []string{
	"zero", "ra", "sp", "gp", "tp", "t0", "t1", "t2",
	"s0", "s1", "a0", "a1", "a2", "a3", "a4", "a5",
	"a6", "a7", "s2", "s3", "s4", "s5", "s6", "s7",
	"s8", "s9", "s10", "s11", "t3", "t4", "t5", "t6",
}
