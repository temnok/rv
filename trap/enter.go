package trap

import (
	"github.com/temnok/rv/arch"
	"github.com/temnok/rv/bi"
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

func IsEntered(cpu *state.CPU) bool {
	return cpu.Update.TrapEnter
}

func EnterWithoutTval(cpu *state.CPU, cause int) {
	Enter(cpu, cause, 0)
}

func Enter(cpu *state.CPU, cause, tval int) {
	cpu.TrapCount++

	if IsEntered(cpu) {
		panic("double Trap")
	}

	mcauseI := arch.XMask
	isInterrupt := bi.T(cause, mcauseI) == 1
	causeID := bi.Ts(cause, 0, 5)

	deleg := cpu.CSR.Medeleg
	if isInterrupt {
		deleg = cpu.CSR.Mideleg
	}

	effectivePriv := state.PrivM
	if cpu.Priv <= state.PrivS && bi.T(deleg, causeID) == 1 {
		effectivePriv = state.PrivS
	}

	cpu.Update.TrapEnter = true
	cpu.Update.TrapPriv = effectivePriv
	cpu.Update.TrapXepc = cpu.PC
	cpu.Update.TrapXcause = cause
	cpu.Update.TrapXtval = tval

	var tvec int

	switch effectivePriv {
	case state.PrivM:
		mie := bi.T(cpu.CSR.Mstatus, csr.MstatusMIE)
		cpu.Update.TrapMstatus = cpu.CSR.Mstatus&^(3<<csr.MstatusMPP|1<<csr.MstatusMPIE|1<<csr.MstatusMIE) |
			cpu.Priv<<csr.MstatusMPP | mie<<csr.MstatusMPIE

		tvec = cpu.CSR.Mtvec

	case state.PrivS:
		sie := bi.T(cpu.CSR.Mstatus, csr.MstatusSIE)
		cpu.Update.TrapMstatus = cpu.CSR.Mstatus&^(1<<csr.MstatusSPP|1<<csr.MstatusSPIE|1<<csr.MstatusSIE) |
			cpu.Priv<<csr.MstatusSPP | sie<<csr.MstatusSPIE

		tvec = cpu.CSR.Stvec
	}

	cpu.Update.TrapPC = tvec &^ 3
	if bi.T(tvec, 0) == 1 && isInterrupt {
		cpu.Update.TrapPC += causeID * 4
	}
}
