package dev

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

func ProcessIO(cpu *state.CPU) {
	if acceptsInput := cpu.Update.Targets&state.UpdateUart == 0 &&
		cpu.CSR.Uart>>csr.UartIP&2 == 0; acceptsInput {

		if b, hasInput := cpu.UARTInput(); hasInput {
			val := cpu.CSR.Uart&^(0xFF<<csr.UartRX) | int(b)<<csr.UartRX | 2<<csr.UartIP
			updateUart(cpu, val)
		}
	}

	if hasOutput := cpu.Update.Targets&state.UpdateUart == 0 &&
		cpu.CSR.Uart>>csr.UartIP&1 == 0; hasOutput {

		b := byte(cpu.CSR.Uart >> csr.UartTX)

		if outputAccepted := cpu.UARTOutput(b); outputAccepted {
			val := cpu.CSR.Uart&^(0xFF<<csr.UartTX) | 1<<csr.UartIP
			updateUart(cpu, val)
		}
	}
}
