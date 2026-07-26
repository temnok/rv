package dev

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

func Read(cpu *state.CPU, addr int) int {
	val := 0

	switch addr {
	case PLIC_claim:
		if (cpu.CSR.Uart>>csr.UartIP)&(cpu.CSR.Uart>>csr.UartIE)&3 != 0 {
			val = 1
		}

	case UART_txdata:
		val = (^cpu.CSR.Uart >> csr.UartIP) & 1 << 31

	case UART_rxdata:
		if (cpu.CSR.Uart>>csr.UartIP)&2 == 0 {
			val = 1 << 31
		} else {
			val = cpu.CSR.Uart >> csr.UartRX & 0xFF

			csr.UpdateUart(cpu, cpu.CSR.Uart&^2<<csr.UartIP|0xFF<<csr.UartRX)
		}

	case UART_ie:
		val = cpu.CSR.Uart >> csr.UartIE & 3

	case UART_ip:
		val = cpu.CSR.Uart >> csr.UartIP & 3
	}

	return val
}
