package uart

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

type UART struct {
	cpu *state.CPU

	reg int
}

func New(cpu *state.CPU) *UART {
	return &UART{
		cpu: cpu,
	}
}

func (uart *UART) Access(addr int, width int, write bool, writeVal int) int {
	var val int

	switch addr {
	case 0x0200_0000: // txdata
		val = (^uart.reg >> csr.UartIP) & 1 << 31

		if write {
			uart.reg = uart.reg&^(1<<csr.UartIP|0xFF<<csr.UartTX) | writeVal&0xFF<<csr.UartTX

			uart.sync()
		}

	case 0x0200_0004: // rxdata
		val = 1 << 31

		if !write && (uart.reg>>csr.UartIP)&2 != 0 {
			val = uart.reg >> csr.UartRX & 0xFF
			uart.reg &^= 2<<csr.UartIP | 0xFF<<csr.UartRX

			uart.sync()
		}

	case 0x0200_0010: // ie
		val = uart.reg >> csr.UartIE & 3

		if write {
			uart.reg = uart.reg&^(3<<csr.UartIE) | writeVal&3<<csr.UartIE

			uart.sync()
		}

	case 0x0200_0014: // ip
		val = uart.reg >> csr.UartIP & 3

	case 0x0320_0004: // PLIC claim
		if (uart.reg>>csr.UartIP)&(uart.reg>>csr.UartIE)&3 != 0 {
			val = 1
		}
	}

	return val
}

func (uart *UART) HasOutput() bool {
	return uart.reg>>csr.UartIP&1 == 0
}

func (uart *UART) GetOutput() byte {
	b := byte(uart.reg >> csr.UartTX)

	uart.reg = uart.reg&^(0xFF<<csr.UartTX) | 1<<csr.UartIP
	uart.sync()

	return b
}

func (uart *UART) AcceptsInput() bool {
	return uart.reg>>csr.UartIP&2 == 0
}

func (uart *UART) SetInput(b byte) {
	uart.reg = uart.reg&^(0xFF<<csr.UartRX) | int(b)<<csr.UartRX | 2<<csr.UartIP

	uart.sync()
}

func (uart *UART) sync() {
	if (uart.reg>>csr.UartIP)&(uart.reg>>csr.UartIE)&3 != 0 {
		uart.cpu.CSR.Mip |= 1 << csr.MipSEIP
	} else {
		uart.cpu.CSR.Mip &^= 1 << csr.MipSEIP
	}
}
