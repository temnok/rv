package trap

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

func IsEntered(cpu *state.CPU) bool {
	return cpu.Update.Targets&state.UpdateCause != 0
}

func EnterWithoutTval(cpu *state.CPU, cause int) {
	Enter(cpu, cause, 0)
}

func Enter(cpu *state.CPU, cause, tval int) {
	cpu.TrapCount++

	if IsEntered(cpu) {
		panic("double Trap")
	}

	isInterrupt := cause>>csr.McauseI&1 == 1
	causeID := cause & 31

	deleg := cpu.CSR.Medeleg
	if isInterrupt {
		deleg = cpu.CSR.Mideleg
	}

	priv := state.PrivM
	clearBits := 3<<csr.MstatusMPP | 1<<csr.MstatusMPIE | 1<<csr.MstatusMIE
	setBits := cpu.Priv<<csr.MstatusMPP | (cpu.CSR.Mstatus>>csr.MstatusMIE&1)<<csr.MstatusMPIE
	tvec := cpu.CSR.Mtvec

	if privS := cpu.Priv <= state.PrivS && deleg>>causeID&1 == 1; privS {
		priv = state.PrivS
		clearBits = 1<<csr.MstatusSPP | 1<<csr.MstatusSPIE | 1<<csr.MstatusSIE
		setBits = cpu.Priv<<csr.MstatusSPP | (cpu.CSR.Mstatus>>csr.MstatusSIE&1)<<csr.MstatusSPIE
		tvec = cpu.CSR.Stvec
	}

	cpu.Update.Targets = state.UpdatePriv | state.UpdateMstatus | state.UpdateEpc | state.UpdateCause | state.UpdateTval

	cpu.Update.Priv = priv
	cpu.Update.Mstatus = cpu.CSR.Mstatus&^clearBits | setBits
	cpu.Update.Cause = cause
	cpu.Update.Tval = tval

	cpu.Update.PC = tvec &^ 3
	if tvec&1 == 1 && isInterrupt {
		cpu.Update.PC += causeID * 4
	}
}
