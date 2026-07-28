package dev

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

func processUartOutput(cpu *state.CPU) {
	if hasOutput := cpu.Update.Targets&state.UpdateUart == 0 &&
		cpu.CSR.Uart>>csr.UartTP&1 == 0; hasOutput {

		b := byte(cpu.CSR.Uart >> csr.UartTD)

		if outputAccepted := cpu.UARTOutput(b); outputAccepted {
			val := cpu.CSR.Uart&^(0xFF<<csr.UartTD) | 1<<csr.UartTP
			updateUart(cpu, val)
		}
	}
}
