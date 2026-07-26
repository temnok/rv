package state

import "github.com/temnok/rv/csr"

type CPU struct {
	PC          int
	X, F        [32]int
	CSR         csr.Registers
	Reservation int

	TLB TLB
	RAM []int

	UARTInput  func() (byte, bool)
	UARTOutput func(byte) bool

	Update Updates
}

func (cpu *CPU) Xset(reg, val int) {
	cpu.Update.Targets |= UpdateXreg
	cpu.Update.Xreg = reg
	cpu.Update.Xval = val
}

func (cpu *CPU) Fset(reg, val int) {
	cpu.Update.Targets |= UpdateFreg | UpdateMstatus
	cpu.Update.Freg = reg
	cpu.Update.Fval = val
	cpu.Update.Mstatus = cpu.CSR.Mstatus | -1<<csr.MstatusSD | 3<<csr.MstatusFS
}

func (cpu *CPU) Cset(reg, val int) bool {
	if reg>>10&3 == 3 {
		return false
	}

	if _, ok := csr.Read(&cpu.CSR, reg); ok {
		cpu.Update.Targets |= UpdateCreg
		cpu.Update.Creg = reg
		cpu.Update.Cval = val

		if reg == csr.Fflags || reg == csr.Frm || reg == csr.Fcsr {
			cpu.Update.Targets |= UpdateMstatus
			cpu.Update.Mstatus = cpu.CSR.Mstatus | -1<<csr.MstatusSD | 3<<csr.MstatusFS
		}

		return true
	}

	return false
}
