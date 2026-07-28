package csr

func Read(csr *Registers, reg int) (int, bool) {
	if csr.Priv < reg>>8&3 {
		return 0, false
	}

	var val int

	switch reg {
	case Cycle:
		val = csr.Mcycle

	case Fcsr:
		val = csr.Fcsr

	case Fflags:
		val = csr.Fcsr & 0x1F

	case Frm:
		val = csr.Fcsr >> 5 & 7

	case Instret:
		val = csr.Mcycle

	case Marchid:

	case Mcause:
		val = csr.Mcause

	case Mcounteren:
		val = csr.Mcounteren

	case Mcountinhibit:

	case Mcycle:
		val = csr.Mcycle

	case Medeleg:
		val = csr.Medeleg

	case Menvcfg:
		val = csr.Menvcfg

	case Mepc:
		val = csr.Mepc

	case Mhartid:

	case Mideleg:
		val = csr.Mideleg

	case Mie:
		val = csr.Mie

	case Mimpid:

	case Minstret:
		val = csr.Mcycle

	case Mip:
		val = csr.Mip

	case Misa:
		val = -1<<63 |
			1<<('i'-'a') | 1<<('m'-'a') | 1<<('a'-'a') | 1<<('c'-'a') |
			1<<('f'-'a') | 1<<('d'-'a') | 1<<('u'-'a') | 1<<('s'-'a')

	case Mscratch:
		val = csr.Mscratch

	case Mstatus:
		val = csr.Mstatus

	case Mtval:
		val = csr.Mtval

	case Mtvec:
		val = csr.Mtvec

	case Mvendorid:

	case Satp:
		if csr.Priv == PrivS && csr.Mstatus>>MstatusTVM&1 == 1 {
			return 0, false
		}
		val = csr.Satp

	case Scause:
		val = csr.Scause

	case Scounteren:
		val = csr.Scounteren

	case Sepc:
		val = csr.Sepc

	case Sie:
		val = csr.Mie & csr.Mideleg

	case Sip:
		val = csr.Mip & csr.Mideleg

	case Sscratch:
		val = csr.Sscratch

	case Sstatus:
		mask := 1<<MstatusSIE | 1<<MstatusSPIE | 1<<MstatusSPP | 1<<MstatusSUM | 1<<MstatusMXR | 3<<MstatusUXL
		val = csr.Mstatus & mask

	case Stimecmp:
		val = csr.Stimecmp

	case Stval:
		val = csr.Stval

	case Stvec:
		val = csr.Stvec

	case Time:
		val = McycleToMtime(csr.Mcycle)

	default:
		return 0, false
	}

	return val, true
}

func McycleToMtime(mcycle int) int {
	return mcycle / 20_000
}
