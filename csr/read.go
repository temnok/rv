package csr

import "github.com/temnok/rv/bit"

func Read(csr *Registers, reg int) (int, bool) {
	if csr.Priv < bit.GetN(reg, 8, 2) {
		return 0, false
	}

	var val int

	switch reg {
	case Fcsr:
		val = csr.Fcsr

	case Fflags:
		val = bit.GetN(csr.Fcsr, 0, 5)

	case Frm:
		val = bit.GetN(csr.Fcsr, 5, 3)

	case Cycle:
		val = csr.Mcycle

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
			1<<('f'-'a') | ('d' - 'a') |
			1<<('u'-'a') | 1<<('s'-'a')

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
		if csr.Priv == PrivS && bit.IsSet(csr.Mstatus, MstatusTVM) {
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
		const mask = 1<<MstatusSIE | 1<<MstatusSUM | 1<<MstatusMXR | 1<<MstatusSPP | 1<<MstatusSPIE | 3<<MstatusUXL
		val = csr.Mstatus & mask

	case Stimecmp:
		if csr.Priv == PrivS && (bit.IsNotSet(csr.Mcounteren, McounterenTM) || bit.IsNotSet(csr.Menvcfg, MenvcfgSTCE)) {
			return 0, false
		}
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
