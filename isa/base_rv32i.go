package isa

type BaseRV32I interface { //          41 =
	baseComputationalInstructions   // 21 +
	baseControlTransferInstructions // 8 +
	baseLoadStoreInstructions       // 8 +
	baseMemoryOrderingInstructions  // 2 +
	baseSystemInstructions          // 2
}

type baseComputationalInstructions interface {
	ADD(rd, rs1, rs2 int)
	ADDI(rd, rs1, imm int)
	AND(rd, rs1, rs2 int)
	ANDI(rd, rs1, imm int)
	AUIPC(rd, imm int)
	LUI(rd, imm int)
	OR(rd, rs1, rs2 int)
	ORI(rd, rs1, imm int)
	SLL(rd, rs1, rs2 int)
	SLLI(rd, rs1, imm int)
	SLT(rd, rs1, rs2 int)
	SLTI(rd, rs1, imm int)
	SLTIU(rd, rs1, imm int)
	SLTU(rd, rs1, rs2 int)
	SRA(rd, rs1, rs2 int)
	SRAI(rd, rs1, imm int)
	SRL(rd, rs1, rs2 int)
	SRLI(rd, rs1, imm int)
	SUB(rd, rs1, rs2 int)
	XOR(rd, rs1, rs2 int)
	XORI(rd, rs1, imm int)
}

type baseControlTransferInstructions interface {
	BEQ(rs1, rs2, imm int)
	BGE(rs1, rs2, imm int)
	BGEU(rs1, rs2, imm int)
	BLT(rs1, rs2, imm int)
	BLTU(rs1, rs2, imm int)
	BNE(rs1, rs2, imm int)
	JAL(rd, imm int)
	JALR(rd, rs1, imm int)
}

type baseLoadStoreInstructions interface {
	LB(rd, rs1, imm int)
	LBU(rd, rs1, imm int)
	LH(rd, rs1, imm int)
	LHU(rd, rs1, imm int)
	LW(rd, rs1, imm int)
	SB(rs2, rs1, imm int)
	SH(rs2, rs1, imm int)
	SW(rs2, rs1, imm int)
}

type baseMemoryOrderingInstructions interface {
	FENCE(pred, succ int)
	FENCE_TSO()
}

type baseSystemInstructions interface {
	EBREAK()
	ECALL()
}
