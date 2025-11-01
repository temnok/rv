package rv

import (
	"github.com/temnok/rv/bi"
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

type CPU struct {
	*state.State
}

type ICache struct {
	VirtAddr, PhysAddr, Value int
}

// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#mcauses
const (
	PageSize = 1 << 12

	PrivU = 0
	PrivS = 1
	PrivM = 3

	AccessExecute = 0
	AccessRead    = 1
	AccessWrite   = 3
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

	if cpu.trapOnPendingInterrupts(); trap.IsEntered(cpu.State) {
		return 0
	}

	var opcode int
	if cpu.memFetch(cpu.PC, &opcode); trap.IsEntered(cpu.State) {
		return 0
	}

	cpu.exec(opcode)

	return opcode
}

// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#privstack
func (cpu *CPU) trapOnPendingInterrupts() {
	cpu.Bus.NotifyInterrupts()

	mi := cpu.CSR.Mip & cpu.CSR.Mie

	if mi == 0 {
		return
	}

	for i := 12; i > 0; i-- {
		if bi.T(mi, i) == 0 {
			continue
		}

		priv := PrivM
		if bi.T(cpu.CSR.Mideleg, i) == 1 {
			priv = PrivS
		}

		mcauseI := cpu.Xlen - 1
		if (priv == cpu.Priv && bi.T(cpu.CSR.Mstatus, priv) == 1) || priv > cpu.Priv {
			trap.EnterWithoutTval(cpu.State, -1<<mcauseI|i)

			return
		}
	}
}
