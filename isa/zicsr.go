package isa

type Zicsr interface {
	CSRRW(rd, csr, rs1 int)
	CSRRC(rd, csr, rs1 int)
	CSRRS(rd, csr, rs1 int)
	CSRRWI(rd, csr, imm int)
	CSRRCI(rd, csr, imm int)
	CSRRSI(rd, csr, imm int)
}
