package dev

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

func Write(cpu *state.CPU, addr int, val int) {
	switch addr {
	case UART_txdata:
		csr.UpdateUart(cpu, cpu.CSR.Uart&^(1<<csr.UartIP|0xFF<<csr.UartTX)|val&0xFF<<csr.UartTX)

	case UART_ie:
		csr.UpdateUart(cpu, cpu.CSR.Uart&^(3<<csr.UartIE)|val&3<<csr.UartIE)
	}
}
