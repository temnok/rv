package isa

type MInstructions interface {
	DIV(rd, rs1, rs2 int)    // DIVide
	DIVU(rd, rs1, rs2 int)   // DIVide Unsigned
	DIVUW(rd, rs1, rs2 int)  // DIVide Unsigned Words
	DIVW(rd, rs1, rs2 int)   // DIVide Words
	MUL(rd, rs1, rs2 int)    // MULtiply
	MULH(rd, rs1, rs2 int)   // MULtiply High
	MULHSU(rd, rs1, rs2 int) // MULtiply High Signed and Unsigned
	MULHU(rd, rs1, rs2 int)  // MULtiply High Unsigned
	MULW(rd, rs1, rs2 int)   // MULtiply Words
	REM(rd, rs1, rs2 int)    // REMinder
	REMU(rd, rs1, rs2 int)   // REMinder of Unsigned
	REMUW(rd, rs1, rs2 int)  // REMinder of Unsigned Words
	REMW(rd, rs1, rs2 int)   // REMinder of Words
}
