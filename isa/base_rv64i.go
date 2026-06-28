package isa

type BaseRV64I interface { //             53 =
	BaseRV32I                          // 41 +
	baseRV64IComputationalInstructions // 9 +
	baseRV64ILoadStoreInstructions     // 3
}

type baseRV64IComputationalInstructions interface {
	ADDIW(rd, rs1, imm int)
	ADDW(rd, rs1, rs2 int)
	SLLIW(rd, rs1, imm int)
	SLLW(rd, rs1, rs2 int)
	SRAIW(rd, rs1, imm int)
	SRAW(rd, rs1, rs2 int)
	SRLIW(rd, rs1, imm int)
	SRLW(rd, rs1, rs2 int)
	SUBW(rd, rs1, rs2 int)
}

type baseRV64ILoadStoreInstructions interface {
	LD(rd, rs1, imm int)
	LWU(rd, rs1, imm int)
	SD(rs2, rs1, imm int)
}
