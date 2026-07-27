package dev

import (
	"github.com/temnok/rv/bit"
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

func processUartInput(cpu *state.CPU) {
	if acceptsInput := cpu.Update.Targets&state.UpdateUart == 0 &&
		bit.IsNotSet(cpu.CSR.Uart, csr.UartIPrx); acceptsInput {

		if b, hasInput := cpu.UARTInput(); hasInput {
			val := cpu.CSR.Uart&^(0xFF<<csr.UartRX) | int(b)<<csr.UartRX | 2<<csr.UartIP
			updateUart(cpu, val)
		}
	}
}
