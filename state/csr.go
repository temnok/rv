package state

type CSR struct {
	Cycle int
	Fcsr  int
	Time  int

	Satp       int
	Scause     int
	Scounteren int
	Sepc       int
	Sip        int
	Sscratch   int
	Stimecmp   int
	Stval      int
	Stvec      int

	Marchid       int
	Mcause        int
	Mcounteren    int
	Mcountinhibit int
	Medeleg       int
	Menvcfg       int
	Mepc          int
	Mhartid       int
	Mideleg       int
	Mie           int
	Mimpid        int
	Mip           int
	Misa          int
	Mscratch      int
	Mstatus       int
	Mtval         int
	Mtvec         int
	Mvendorid     int

	TimerCallbacks []func()
}
