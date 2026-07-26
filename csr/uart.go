package csr

import "github.com/temnok/rv/state"

func UpdateUart(cpu *state.CPU, val int) {
	cpu.Update.Targets |= state.UpdateUart
	cpu.Update.Uart = val

	seip := cpu.CSR.Mip >> MipSEIP & 1
	if pend := ((val>>UartIP)&(val>>UartIE)&3 != 0); pend != (seip == 1) {
		if cpu.Update.Targets&state.UpdateMip == 0 {
			cpu.Update.Targets |= state.UpdateMip
			cpu.Update.Mip = cpu.CSR.Mip
		}

		cpu.Update.Mip &^= 1 << MipSEIP
		cpu.Update.Mip |= (1 - seip) << MipSEIP
	}
}
