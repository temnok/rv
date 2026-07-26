package csr

func Write(csr *Registers, reg, val int) {
	switch reg {
	case Fcsr:
		csr.Fcsr = val & 0xFF

	case Fflags:
		mask := 0x1F
		csr.Fcsr = csr.Fcsr&^mask | val&mask

	case Frm:
		mask := 0xE0
		csr.Fcsr = csr.Fcsr&^mask | (val<<5)&mask

	case Mcause:
		csr.Mcause = val

	case Mcounteren:
		csr.Mcounteren = val

	case Medeleg:
		csr.Medeleg = val

	case Menvcfg:
		csr.Menvcfg = val

	case Mepc:
		csr.Mepc = val

	case Mideleg:
		csr.Mideleg = val

	case Mie:
		csr.Mie = val

	case Mip:
		csr.Mip = val

	case Mscratch:
		csr.Mscratch = val

	case Mstatus:
		mask := ^(3<<MstatusSXL | 3<<MstatusUXL)
		csr.Mstatus = csr.Mstatus&^mask | val&mask

	case Mtval:
		csr.Mtval = val

	case Mtvec:
		csr.Mtvec = val

	case Satp:
		csr.Satp = val

	case Scause:
		csr.Scause = val

	case Scounteren:
		csr.Scounteren = val

	case Sepc:
		csr.Sepc = val

	case Sie:
		csr.Mie = csr.Mie&^csr.Mideleg | val&csr.Mideleg

	case Sscratch:
		csr.Sscratch = val

	case Sstatus:
		mask := 1<<MstatusSIE | 1<<MstatusSUM | 1<<MstatusMXR | 1<<MstatusSPP
		csr.Mstatus = csr.Mstatus&^mask | val&mask

	case Stimecmp:
		csr.Stimecmp = val

	case Stval:
		csr.Stval = val

	case Stvec:
		csr.Stvec = val
	}
}
