package trap

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

func IsEntered(cpu *state.CPU) bool {
	return cpu.Update.Targets&(state.UpdateMcause|state.UpdateScause) != 0
}

func Enter(cpu *state.CPU, cause, tval int) {
	isInterrupt := cause>>csr.McauseI&1 == 1
	causeID := cause & 0x3F

	deleg := cpu.CSR.Medeleg
	if isInterrupt {
		deleg = cpu.CSR.Mideleg
	}

	priv := csr.PrivM
	clr := 3<<csr.MstatusMPP | 1<<csr.MstatusMPIE | 1<<csr.MstatusMIE
	set := cpu.CSR.Priv<<csr.MstatusMPP | (cpu.CSR.Mstatus>>csr.MstatusMIE&1)<<csr.MstatusMPIE
	tvec := cpu.CSR.Mtvec
	cpu.Update.Targets = state.UpdateMepc | state.UpdateMcause | state.UpdateMtval

	if privS := cpu.CSR.Priv <= csr.PrivS && deleg>>causeID&1 == 1; privS {
		priv = csr.PrivS
		clr = 1<<csr.MstatusSPP | 1<<csr.MstatusSPIE | 1<<csr.MstatusSIE
		set = cpu.CSR.Priv<<csr.MstatusSPP | (cpu.CSR.Mstatus>>csr.MstatusSIE&1)<<csr.MstatusSPIE
		tvec = cpu.CSR.Stvec
		cpu.Update.Targets = state.UpdateSepc | state.UpdateScause | state.UpdateStval
	}

	cpu.Update.Targets |= state.UpdatePC | state.UpdatePriv | state.UpdateMstatus
	cpu.Update.Priv = priv
	cpu.Update.Mstatus = cpu.CSR.Mstatus&^clr | set
	cpu.Update.Xepc = cpu.PC
	cpu.Update.Xcause = cause
	cpu.Update.Xtval = tval

	cpu.Update.PC = tvec &^ 3
	if tvec&1 == 1 && isInterrupt {
		cpu.Update.PC += causeID * 4
	}
}
