package uart

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

type UART struct {
	cpu *state.CPU

	reg int
}

const (
	UartTX = 0
	UartRX = 8
	UartIE = 16
	UartIP = 24
)

func New(cpu *state.CPU) *UART {
	return &UART{
		cpu: cpu,

		reg: 1 << UartIP,
	}
}

func (uart *UART) Access(addr int, width int, write bool, writeVal int) int {
	addr &= 0x1ff_ffff

	var val int

	switch addr {
	case 0x00: // txdata
		val = (^uart.reg >> UartIP) & 1 << 31

		if write {
			uart.reg = uart.reg&^(1<<UartIP|0xFF<<UartTX) | writeVal&0xFF<<UartTX

			uart.sync()
		}

	case 0x04: // rxdata
		val = 1 << 31

		if !write && (uart.reg>>UartIP)&2 != 0 {
			val = uart.reg >> UartRX & 0xFF
			uart.reg &^= 2<<UartIP | 0xFF<<UartRX

			uart.sync()
		}

	case 0x10: // ie
		val = uart.reg >> UartIE & 3

		if write {
			uart.reg = uart.reg&^(3<<UartIE) | writeVal&3<<UartIE

			uart.sync()
		}

	case 0x14: // ip
		val = uart.reg >> UartIP & 3

	case 0x120_0004: // PLIC claim
		if (uart.reg>>UartIP)&(uart.reg>>UartIE)&3 != 0 {
			val = 1
		}
	}

	return val
}

func (uart *UART) HasOutput() bool {
	return uart.reg>>UartIP&1 == 0
}

func (uart *UART) GetOutput() byte {
	b := byte(uart.reg >> UartTX)

	uart.reg = uart.reg&^(0xFF<<UartTX) | 1<<UartIP
	uart.sync()

	return b
}

func (uart *UART) AcceptsInput() bool {
	return uart.reg>>UartIP&2 == 0
}

func (uart *UART) SetInput(b byte) {
	uart.reg = uart.reg&^(0xFF<<UartRX) | int(b)<<UartRX | 2<<UartIP

	uart.sync()
}

func (uart *UART) sync() {
	if (uart.reg>>UartIP)&(uart.reg>>UartIE)&3 != 0 {
		uart.cpu.CSR.Mip |= 1 << csr.MipSEIP
	} else {
		uart.cpu.CSR.Mip &^= 1 << csr.MipSEIP
	}
}
