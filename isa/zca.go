package isa

type ZcaInstructions interface { //                      34 =
	compressedComputationalInstructions   // 18 +
	compressedControlTransferInstructions //  6 +
	compressedLoadStoreInstructions       //  8 +
	compressedSpecialInstructions         //  2
}

type compressedComputationalInstructions interface {
	C_ADD(rd, rs2 int)      // ADD
	C_ADDI(rd, imm int)     // ADD Immediate
	C_ADDI16SP(rd, imm int) // ADD Immediate (multiples of 16) to Stack Pointer
	C_ADDI4SPN(rd, imm int) // ADD Immediate (multiples of 4) to Stack Pointer, Non-destructive
	C_ADDIW(rd, imm int)    // ADD Immediate, Word
	C_ADDW(rd, rs int)      // ADD, Word
	C_AND(rd, rs int)       // AND
	C_ANDI(rd, imm int)     // AND with Immediate
	C_LI(rd, imm int)       // Load Immediate
	C_LUI(rd, imm int)      // Load Upper Immediate
	C_MV(rd, rs int)        // MoVe
	C_OR(rd, rs int)        // OR
	C_SLLI(rd, shamt int)   // Shift Left, Logical by Immediate
	C_SRAI(rd, shamt int)   // Shift Right, Arithmetic by Immediate
	C_SRLI(rd, shamt int)   // Shift Right, Logical by Immediate
	C_SUB(rd, rs int)       // SUBtract
	C_SUBW(rd, rs int)      // SUBtract, Word
	C_XOR(rd, rs int)       // eXclusive OR
}

type compressedControlTransferInstructions interface {
	C_J(offset int)        // Jump by offset
	C_JAL(offset int)      // Jump And Link by offset
	C_JR(rs int)           // Jump by Register
	C_JALR(rs int)         // Jump And Link by Register
	C_BEQZ(rs, offset int) // Branch if EQual to Zero
	C_BNEZ(rs, offset int) // Branch if Not Equal to Zero
}

type compressedLoadStoreInstructions interface {
	C_LD(rd, rs, offset int)   // Load Double word
	C_LDSP(rd, offset int)     // Load Double word relative to Stack Pointer
	C_LW(rd, rs, offset int)   // Load Word
	C_LWSP(rd, offset int)     // Load Word relative to Stack Pointer
	C_SD(rs2, rs1, offset int) // Store Double word
	C_SDSP(rs, offset int)     // Store Double word relative to Stack Pointer
	C_SW(rs2, rs1, offset int) // Store Word
	C_SWSP(rs, offset int)     // Store Word relative to Stack Pointer
}

type compressedSpecialInstructions interface {
	C_EBREAK() // Environment Break
	C_NOP()    // No-OP
}
