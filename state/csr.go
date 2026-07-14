package state

type CSR struct {
	Fcsr int // https://docs.riscv.org/reference/isa/v20260120/unpriv/f-st-ext.html#fcsr

	Mcause     int
	Mcounteren int
	Mcycle     int
	Medeleg    int
	Menvcfg    int // https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#sec:menvcfg
	Mepc       int
	Mideleg    int
	Mie        int
	Mip        int
	Mscratch   int
	Mstatus    int
	Mtval      int
	Mtvec      int

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

func (csr *CSR) Mtime() uint {
	return uint(csr.Mcycle) / 20_000
}
