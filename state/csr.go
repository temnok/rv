package state

type CSR struct {
	Fcsr, Cycle, Cycleh, Time, Timeh int

	Stvec, Scounteren, Sscratch, Sepc, Scause, Stval, Sip, Satp int

	Mstatus, Mstatush, Misa, Medeleg, Mideleg, Mie, Mtvec, Mcounteren int
	Mscratch, Mepc, Mcause, Mtval, Mip                                int
	Mvendorid, Marchid, Mimpid, Mhartid                               int
}
