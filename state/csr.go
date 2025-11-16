package state

type CSR struct {
	Cycle  int
	Cycleh int
	Fcsr   int
	Time   int
	Timeh  int

	Satp       int
	Scause     int
	Scounteren int
	Sepc       int
	Sip        int
	Sscratch   int
	Stval      int
	Stvec      int

	Marchid    int
	Mcause     int
	Mcounteren int
	Medeleg    int
	Mepc       int
	Mhartid    int
	Mideleg    int
	Mie        int
	Mimpid     int
	Mip        int
	Misa       int
	Mscratch   int
	Mstatus    int
	Mstatush   int
	Mtval      int
	Mtvec      int
	Mvendorid  int
}
