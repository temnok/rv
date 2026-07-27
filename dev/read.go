package dev

import (
	"github.com/temnok/rv/bit"
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

const (
	PLIC_claim  = 0x0C20_0004
	UART_txdata = 0x1001_0000
	UART_rxdata = 0x1001_0004
	UART_ie     = 0x1001_0010
	UART_ip     = 0x1001_0014
)

func Read(cpu *state.CPU, addr int) int {
	val := 0

	switch addr {
	case PLIC_claim:
		if bit.GetN(cpu.CSR.Uart, csr.UartIP, 2)&bit.GetN(cpu.CSR.Uart, csr.UartIE, 2) != 0 {
			val = 1
		}

	case UART_txdata:
		if bit.IsNotSet(cpu.CSR.Uart, csr.UartIPtx) {
			val = 1 << 31
		}

	case UART_rxdata:
		if bit.IsNotSet(cpu.CSR.Uart, csr.UartIPrx) {
			val = 1 << 31
		} else {
			val = bit.GetN(cpu.CSR.Uart, csr.UartRX, 8)

			updateUart(cpu, cpu.CSR.Uart&^2<<csr.UartIP|0xFF<<csr.UartRX)
		}

	case UART_ie:
		val = bit.GetN(cpu.CSR.Uart, csr.UartIE, 2)

	case UART_ip:
		val = bit.GetN(cpu.CSR.Uart, csr.UartIP, 2)
	}

	return val
}
