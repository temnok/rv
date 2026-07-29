package device

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

func syncUartInput(cpu *state.CPU) {
	if acceptsInput := cpu.Update.Targets&state.UpdateUart == 0 &&
		cpu.CSR.Uart>>csr.UartRP&1 == 0; acceptsInput {

		if b, hasInput := cpu.UARTInput(); hasInput {
			val := cpu.CSR.Uart&^(0xFF<<csr.UartRD) | int(b)<<csr.UartRD | 1<<csr.UartRP
			updateUart(cpu, val)
		}
	}
}
