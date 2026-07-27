package dev

import (
	"github.com/temnok/rv/bit"
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

func processUartOutput(cpu *state.CPU) {
	if hasOutput := cpu.Update.Targets&state.UpdateUart == 0 &&
		bit.IsNotSet(cpu.CSR.Uart, csr.UartIPtx); hasOutput {

		b := byte(bit.GetN(cpu.CSR.Uart, csr.UartTX, 8))

		if outputAccepted := cpu.UARTOutput(b); outputAccepted {
			val := cpu.CSR.Uart&^(0xFF<<csr.UartTX) | 1<<csr.UartIP
			updateUart(cpu, val)
		}
	}
}
