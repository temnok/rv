package rv

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/mem"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func NewCPU(xlen int, startAddr int) *state.CPU {
	xl := xlen / 32

	cpu := &state.CPU{
		Fixed: state.Fixed{
			Xlen: xlen,
		},

		Static: state.Static{
			Priv: state.PrivM,

			CSR: state.CSR{
				Misa: xl<<(xlen-2) |
					1<<('i'-'a') | 1<<('m'-'a') | 1<<('a'-'a') | 1<<('c'-'a') |
					1<<('f'-'a') | ('d' - 'a') |
					1<<('u'-'a') | 1<<('s'-'a'),
			},
		},

		Update: state.Updated{
			XReg: -1,
			CReg: -1,
		},
	}

	cpu.CSR.Mstatus = cpu.Xint(xl<<csr.MstatusSXL | xl<<csr.MstatusUXL)
	cpu.Update.PC = cpu.Xint(startAddr)

	return cpu
}

func Step(cpu *state.CPU) bool {
	//return debugStep(cpu)

	innerStep(cpu)
	return true
}

func innerStep(cpu *state.CPU) int {
	updateState(cpu)

	if trap.OnPendingInterrupts(cpu); trap.IsEntered(cpu) {
		return 0
	}

	var opcode int
	if mem.Fetch(cpu, cpu.PC, &opcode); trap.IsEntered(cpu) {
		return 0
	}

	Exec(cpu, opcode)

	return opcode
}
