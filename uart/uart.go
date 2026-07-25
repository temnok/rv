package uart

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

type UART struct {
	cpu *state.CPU
}

func New(cpu *state.CPU) *UART {
	return &UART{
		cpu: cpu,
	}
}

func (uart *UART) Access(addr int, width int, write bool, writeVal int) int {
	cpu := uart.cpu

	var val int

	switch addr {
	case 0x0C20_0004: // PLIC claim
		if (cpu.CSR.Uart>>csr.UartIP)&(cpu.CSR.Uart>>csr.UartIE)&3 != 0 {
			val = 1
		}

	case 0x1001_0000: // txdata
		val = (^cpu.CSR.Uart >> csr.UartIP) & 1 << 31

		if write {
			cpu.CSR.Uart = cpu.CSR.Uart&^(1<<csr.UartIP|0xFF<<csr.UartTX) | writeVal&0xFF<<csr.UartTX

			uart.sync()
		}

	case 0x1001_0004: // rxdata
		val = 1 << 31

		if !write && (cpu.CSR.Uart>>csr.UartIP)&2 != 0 {
			val = cpu.CSR.Uart >> csr.UartRX & 0xFF
			cpu.CSR.Uart &^= 2<<csr.UartIP | 0xFF<<csr.UartRX

			uart.sync()
		}

	case 0x1001_0010: // ie
		val = cpu.CSR.Uart >> csr.UartIE & 3

		if write {
			cpu.CSR.Uart = cpu.CSR.Uart&^(3<<csr.UartIE) | writeVal&3<<csr.UartIE

			uart.sync()
		}

	case 0x1001_0014: // ip
		val = cpu.CSR.Uart >> csr.UartIP & 3
	}

	return val
}

func (uart *UART) HasOutput() bool {
	return uart.cpu.CSR.Uart>>csr.UartIP&1 == 0
}

func (uart *UART) GetOutput() byte {
	b := byte(uart.cpu.CSR.Uart >> csr.UartTX)

	uart.cpu.CSR.Uart = uart.cpu.CSR.Uart&^(0xFF<<csr.UartTX) | 1<<csr.UartIP
	uart.sync()

	return b
}

func (uart *UART) AcceptsInput() bool {
	return uart.cpu.CSR.Uart>>csr.UartIP&2 == 0
}

func (uart *UART) SetInput(b byte) {
	uart.cpu.CSR.Uart = uart.cpu.CSR.Uart&^(0xFF<<csr.UartRX) | int(b)<<csr.UartRX | 2<<csr.UartIP

	uart.sync()
}

func (uart *UART) sync() {
	if (uart.cpu.CSR.Uart>>csr.UartIP)&(uart.cpu.CSR.Uart>>csr.UartIE)&3 != 0 {
		uart.cpu.CSR.Mip |= 1 << csr.MipSEIP
	} else {
		uart.cpu.CSR.Mip &^= 1 << csr.MipSEIP
	}
}
