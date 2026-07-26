package dev

import (
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

func HasOutput(cpu *state.CPU) bool {
	return cpu.CSR.Uart>>csr.UartIP&1 == 0
}

func GetOutput(cpu *state.CPU) byte {
	b := byte(cpu.CSR.Uart >> csr.UartTX)

	updateUart(cpu, cpu.CSR.Uart&^(0xFF<<csr.UartTX)|1<<csr.UartIP)

	return b
}

func AcceptsInput(cpu *state.CPU) bool {
	return cpu.CSR.Uart>>csr.UartIP&2 == 0
}

func SetInput(cpu *state.CPU, b byte) {
	updateUart(cpu, cpu.CSR.Uart&^(0xFF<<csr.UartRX)|int(b)<<csr.UartRX|2<<csr.UartIP)
}

func updateUart(cpu *state.CPU, uart int) {
	cpu.CSR.Uart = uart

	if (uart>>csr.UartIP)&(uart>>csr.UartIE)&3 != 0 {
		cpu.CSR.Mip |= 1 << csr.MipSEIP
	} else {
		cpu.CSR.Mip &^= 1 << csr.MipSEIP
	}
}
