package device

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

func updateUart(cpu *state.CPU, val int) {
	cpu.Update.Targets |= state.UpdateUart
	cpu.Update.Uart = val

	pend := (val>>csr.UartTP)&(val>>csr.UartTE)&1 | (val>>csr.UartRP)&(val>>csr.UartRE)&1
	seip := cpu.CSR.Mip >> csr.MipSEIP & 1
	if pend != seip {
		if cpu.Update.Targets&state.UpdateMip == 0 {
			cpu.Update.Targets |= state.UpdateMip
			cpu.Update.Mip = cpu.CSR.Mip
		}

		cpu.Update.Mip &^= 1 << csr.MipSEIP
		cpu.Update.Mip |= (1 - seip) << csr.MipSEIP
	}
}
