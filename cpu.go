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

func (cpu *CPU) Init(xlen int, bus state.Bus, startAddr int) {
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

func (cpu *CPU) extD() bool {
	return cpu.CSR.Misa&1<<('d'-'a') != 0
}

func (cpu *CPU) Step() bool {
	//return cpu.debugStep()

	cpu.innerStep()
	return true
}

func (cpu *CPU) innerStep() int {
	cpu.updateState()

	if trap.OnPendingInterrupts(cpu.State); trap.IsEntered(cpu.State) {
		return 0
	}

	var opcode int
	if mem.Fetch(cpu.State, cpu.PC, &opcode); trap.IsEntered(cpu.State) {
		return 0
	}

	cpu.exec(opcode)

	return opcode
}
