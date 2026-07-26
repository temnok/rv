package state

import (
	"github.com/temnok/rv/csr"
)

func UpdateCSR(cpu *CPU, reg, val int) bool {
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
