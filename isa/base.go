package isa

type BaseInstructions interface { //   53 =
	baseComputationalInstructions   // 30 +
	baseControlTransferInstructions //  8 +
	baseLoadStoreInstructions       // 11 +
	baseMemoryOrderingInstructions  //  2 +
	baseSystemInstructions          //  2
}

type baseComputationalInstructions interface {
	ADD(rd, rs1, rs2 int)   // ADD
	ADDI(rd, rs1, imm int)  // ADD Immediate
	ADDIW(rd, rs1, imm int) // ADD Immediate, Word
	ADDW(rd, rs1, rs2 int)  // ADD, Word
	AND(rd, rs1, rs2 int)   // AND
	ANDI(rd, rs1, imm int)  // AND with Immediate
	AUIPC(rd, imm int)      // Add Upper Immediate to PC
	LUI(rd, imm int)        // Load Upper Immediate
	OR(rd, rs1, rs2 int)    // OR
	ORI(rd, rs1, imm int)   // OR with Immediate
	SLL(rd, rs1, rs2 int)   // Shift Left, Logical
	SLLI(rd, rs1, imm int)  // Shift Left, Logical by Immediate
	SLLIW(rd, rs1, imm int) // Shift Left, Logical by Immediate, Word
	SLLW(rd, rs1, rs2 int)  // Shift Left, Word
	SLT(rd, rs1, rs2 int)   // Set if Less Than
	SLTI(rd, rs1, imm int)  // Set if Less Than, with Immediate
	SLTIU(rd, rs1, imm int) // Set if Less Than, with Immediate, Unsigned
	SLTU(rd, rs1, rs2 int)  // Set if Less Than, Unsigned
	SRA(rd, rs1, rs2 int)   // Shift Right, Arithmetical
	SRAI(rd, rs1, imm int)  // Shift Right, Arithmetical by Immediate
	SRAIW(rd, rs1, imm int) // Shift Right, Arithmetical by Immediate, Word
	SRAW(rd, rs1, rs2 int)  // Shift Right, Arithmetical, Word
	SRL(rd, rs1, rs2 int)   // Shift Right, Logical
	SRLI(rd, rs1, imm int)  // Shift Right, Logical with Immediate
	SRLIW(rd, rs1, imm int) // Shift Right, Logical with Immediate, Word
	SRLW(rd, rs1, rs2 int)  // Shift Right, Logical, Word
	SUB(rd, rs1, rs2 int)   // SUBtract
	SUBW(rd, rs1, rs2 int)  // SUBtract, Word
	XOR(rd, rs1, rs2 int)   // eXclusive OR
	XORI(rd, rs1, imm int)  // eXclusive OR with Immediate
}

type baseControlTransferInstructions interface {
	BEQ(rs1, rs2, offset int)  // Branch if EQual
	BGE(rs1, rs2, offset int)  // Branch if Greater or Equal
	BGEU(rs1, rs2, offset int) // Branch if Greater or Equal, Unsigned operands
	BLT(rs1, rs2, offset int)  // Branch if Less Than
	BLTU(rs1, rs2, offset int) // Branch if Less Than, Unsigned operands
	BNE(rs1, rs2, offset int)  // Branch if Not Equal
	JAL(rd, offset int)        // Jump And Link by immediate offset
	JALR(rd, rs1, offset int)  // Jump And Link to address in register
}

type baseLoadStoreInstructions interface {
	LB(rd, rs1, offset int)  // Load Byte
	LBU(rd, rs1, offset int) // Load Byte, Unsigned
	LD(rd, rs1, offset int)  // Load Double word
	LH(rd, rs1, offset int)  // Load Half word
	LHU(rd, rs1, offset int) // Load Half word, Unsigned
	LW(rd, rs1, offset int)  // Load Word
	LWU(rd, rs1, offset int) // Load Word, Unsigned
	SB(rs2, rs1, offset int) // Store Byte
	SD(rs2, rs1, offset int) // Store Double word
	SH(rs2, rs1, offset int) // Store Half word
	SW(rs2, rs1, offset int) // Store Word
}

type baseMemoryOrderingInstructions interface {
	FENCE(pred, succ int) // FENCE
	FENCE_TSO()           // FENCE with Total Store Ordering
}

type baseSystemInstructions interface {
	EBREAK() // Environment BREAK
	ECALL()  // Environment CALL
}
