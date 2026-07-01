package ins

import (
	"fmt"
	"github.com/temnok/rv/isa"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

var sbiCallCount int

func ecall(cpu *state.CPU, op Op) {
	if cpu.Priv == 1 {
		eid, fid, a0 := cpu.X[isa.A7], cpu.X[isa.A6], cpu.X[isa.A0]

		if sbiCallCount++; sbiCallCount < 1000 {
			fmt.Printf("\r\n*** SBI call: EID=0x%x, FID=%v, a0=0x%x\r\n", eid, fid, a0)
		}

		err, res := 0, 0

		switch eid {
		case 0x10:
			switch fid {
			case 0:
				//res = 3<<24 | 0
				res = 3<<24 | 0
			case 3:
				res = 1
				//	switch a0 {
				//	case 0x4442434E, 0x54494d45:
				//		res = 1
				//	}
			}
		default:
			err = -2
		}

		cpu.X[isa.A0] = err
		if err == 0 {
			cpu.X[isa.A1] = res
		}
		return
	}

	trap.EnterWithoutTval(cpu, trap.EnvironmentCallFromUMode+cpu.Priv)
}
