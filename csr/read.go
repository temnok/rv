package csr

import (
	"github.com/temnok/rv/state"
)

func Read(cpu *state.CPU, reg int) (int, bool) {
	if cpu.Priv < reg>>8&3 {
		return 0, false
	}

	var val int

	switch csr := &cpu.CSR; reg {
	case Fcsr:
		val = csr.Fcsr
	case Fflags:
		val = csr.Fcsr & 0x1F
	case Frm:
		val = csr.Fcsr >> 5
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
		val = state.Misa
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
		if cpu.Priv == state.PrivS && csr.Mstatus>>MstatusTVM&1 == 1 {
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
		val = csr.Mie & (1<<MipSEI | 1<<MipSTI)
	case Sip:
		val = csr.Mip & (1 << MipSEI)
	case Sscratch:
		val = csr.Sscratch
	case Sstatus:
		val = csr.Mstatus & (1<<MstatusSIE | 1<<MstatusSUM | 1<<MstatusMXR | 1<<MstatusSPP | 1<<MstatusSPIE | 3<<MstatusUXL)
	case Stimecmp:
		if cpu.Priv == state.PrivS && (csr.Mcounteren>>McounterenTM&1 == 0 || csr.Menvcfg>>MenvcfgSTCE&1 == 0) {
			return 0, false
		}
		val = csr.Stimecmp
	case Stval:
		val = csr.Stval
	case Stvec:
		val = csr.Stvec
	case Time:
		val = int(csr.Mtime())
	default:
		return 0, false
	}

	return val, true
}
