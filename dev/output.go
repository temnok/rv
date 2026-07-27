package dev

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

func processUartOutput(cpu *state.CPU) {
	if hasOutput := cpu.Update.Targets&state.UpdateUart == 0 &&
		cpu.CSR.Uart>>csr.UartIPtx&1 == 0; hasOutput {

		b := byte(cpu.CSR.Uart >> csr.UartTX)

		if outputAccepted := cpu.UARTOutput(b); outputAccepted {
			val := cpu.CSR.Uart&^(0xFF<<csr.UartTX) | 1<<csr.UartIP
			updateUart(cpu, val)
		}
	}
}
