package state

type CSR struct {
	Fcsr int

	Marchid       int
	Mcause        int
	Mcounteren    int
	Mcountinhibit int
	Mcycle        int
	Medeleg       int
	Menvcfg       int
	Mepc          int
	Mhartid       int
	Mideleg       int
	Mie           int
	Mimpid        int
	Minstret      int
	Mip           int
	Misa          int
	Mscratch      int
	Mstatus       int
	Mtime         int
	Mtval         int
	Mtvec         int
	Mvendorid     int

	Satp       int
	Scause     int
	Scounteren int
	Sepc       int
	Sip        int
	Sscratch   int
	Stimecmp   int
	Stval      int
	Stvec      int
}
