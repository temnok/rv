package trap

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

func IsEntered(cpu *state.CPU) bool {
	return cpu.Update.Targets&(state.UpdateMepc|state.UpdateSepc) != 0
}

func Enter(cpu *state.CPU, cause, tval int) {
	isInterrupt := cause>>csr.McauseI&1 == 1
	causeID := cause & 0x3F

	var deleg int
	if isInterrupt {
		deleg = cpu.CSR.Mideleg
	} else {
		deleg = cpu.CSR.Medeleg
	}

	targets := state.UpdatePriv | state.UpdateMstatus | state.UpdatePC

	var tvec int
	if cpu.CSR.Priv == csr.PrivM || deleg>>causeID&1 == 0 {
		cpu.Update.Priv = csr.PrivM
		cpu.Update.Targets = targets | state.UpdateMepc | state.UpdateMcause | state.UpdateMtval
		cpu.Update.Mstatus = cpu.CSR.Mstatus&^(3<<csr.MstatusMPP|1<<csr.MstatusMPIE|1<<csr.MstatusMIE) |
			cpu.CSR.Priv<<csr.MstatusMPP | (cpu.CSR.Mstatus>>csr.MstatusMIE&1)<<csr.MstatusMPIE
		tvec = cpu.CSR.Mtvec
	} else {
		cpu.Update.Priv = csr.PrivS
		cpu.Update.Targets = targets | state.UpdateSepc | state.UpdateScause | state.UpdateStval
		cpu.Update.Mstatus = cpu.CSR.Mstatus&^(1<<csr.MstatusSPP|1<<csr.MstatusSPIE|1<<csr.MstatusSIE) |
			cpu.CSR.Priv<<csr.MstatusSPP | (cpu.CSR.Mstatus>>csr.MstatusSIE&1)<<csr.MstatusSPIE
		tvec = cpu.CSR.Stvec
	}

	base := tvec &^ 3
	if tvec&1 == 1 && isInterrupt {
		cpu.Update.PC = base + causeID*4
	} else {
		cpu.Update.PC = base
	}

	cpu.Update.Xepc = cpu.PC
	cpu.Update.Xcause = cause
	cpu.Update.Xtval = tval
}
