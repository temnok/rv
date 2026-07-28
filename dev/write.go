package dev

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

func Write(cpu *state.CPU, addr int, val int) {
	switch addr {
	case UART_txdata:
		mask := 1<<csr.UartTP | 0xFF<<csr.UartTD
		val &= 0xFF
		updateUart(cpu, cpu.CSR.Uart&^mask|val<<csr.UartTD)

	case UART_ie:
		mask := 1<<csr.UartTE | 1<<csr.UartRE
		te := val & 1
		re := val >> 1 & 1
		updateUart(cpu, cpu.CSR.Uart&^mask|te<<csr.UartTE|re<<csr.UartRE)
	}
}
