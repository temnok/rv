package rv

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/mem"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

type CPU struct {
	*state.State
}

const (
	PrivM = 3
)

func Init(cpu *CPU, xlen int, bus state.Bus, startAddr int) {
	xl := xlen / 32

	*cpu = CPU{
		State: &state.State{
			Bus: bus,

			Fixed: state.Fixed{
				Xlen: xlen,
			},

			Static: state.Static{
				Priv: PrivM,

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
		},
	}

	cpu.CSR.Mstatus = cpu.Xint(xl<<csr.MstatusSXL | xl<<csr.MstatusUXL)
	cpu.Update.PC = cpu.Xint(startAddr)
}

func Step(cpu *state.State) bool {
	//return debugStep(cpu)

	innerStep(cpu)
	return true
}

func innerStep(cpu *state.State) int {
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
