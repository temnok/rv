package isa

type BaseInstructions interface { //   53 =
	baseComputationalInstructions   // 30 +
	baseControlTransferInstructions //  8 +
	baseLoadStoreInstructions       // 11 +
	baseMemoryOrderingInstructions  //  2 +
	baseSystemInstructions          //  2
}

type baseComputationalInstructions interface {
	ADD(rd, rs1, rs2 int)
	ADDI(rd, rs1, imm int)
	ADDIW(rd, rs1, imm int)
	ADDW(rd, rs1, rs2 int)
	AND(rd, rs1, rs2 int)
	ANDI(rd, rs1, imm int)
	AUIPC(rd, imm int)
	LUI(rd, imm int)
	OR(rd, rs1, rs2 int)
	ORI(rd, rs1, imm int)
	SLL(rd, rs1, rs2 int)
	SLLI(rd, rs1, imm int)
	SLLIW(rd, rs1, imm int)
	SLLW(rd, rs1, rs2 int)
	SLT(rd, rs1, rs2 int)
	SLTI(rd, rs1, imm int)
	SLTIU(rd, rs1, imm int)
	SLTU(rd, rs1, rs2 int)
	SRA(rd, rs1, rs2 int)
	SRAI(rd, rs1, imm int)
	SRAIW(rd, rs1, imm int)
	SRAW(rd, rs1, rs2 int)
	SRL(rd, rs1, rs2 int)
	SRLI(rd, rs1, imm int)
	SRLIW(rd, rs1, imm int)
	SRLW(rd, rs1, rs2 int)
	SUB(rd, rs1, rs2 int)
	SUBW(rd, rs1, rs2 int)
	XOR(rd, rs1, rs2 int)
	XORI(rd, rs1, imm int)
}

type baseControlTransferInstructions interface {
	BEQ(rs1, rs2, offset int)
	BGE(rs1, rs2, offset int)
	BGEU(rs1, rs2, offset int)
	BLT(rs1, rs2, offset int)
	BLTU(rs1, rs2, offset int)
	BNE(rs1, rs2, offset int)
	JAL(rd, offset int)
	JALR(rd, rs1, offset int)
}

type baseLoadStoreInstructions interface {
	LB(rd, rs1, offset int)
	LBU(rd, rs1, offset int)
	LD(rd, rs1, offset int)
	LH(rd, rs1, offset int)
	LHU(rd, rs1, offset int)
	LW(rd, rs1, offset int)
	LWU(rd, rs1, offset int)
	SB(rs2, rs1, offset int)
	SD(rs2, rs1, offset int)
	SH(rs2, rs1, offset int)
	SW(rs2, rs1, offset int)
}

type baseMemoryOrderingInstructions interface {
	FENCE(pred, succ int)
	FENCE_TSO()
}

type baseSystemInstructions interface {
	EBREAK()
	ECALL()
}
