package uart

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

type UART struct {
	cpu *state.CPU

	tx, rx, ie int
}

func New(cpu *state.CPU) *UART {
	return &UART{
		cpu: cpu,

		tx: -1,
		rx: -1,
	}
}

func (uart *UART) Access(addr int, width int, write bool, writeData int) int {
	addr &= 0x1ff_ffff

	var val int

	switch addr {
	case 0x00: // txdata
		if write {
			uart.tx = writeData & 0xFF

			uart.sync()
		} else {
			if uart.tx >= 0 {
				val = 1 << 31
			}
		}

	case 0x04: // rxdata
		if !write {
			if uart.rx < 0 {
				val = 1 << 31
			} else {
				val = uart.rx
				uart.rx = -1

				uart.sync()
			}
		}

	case 0x10: // ie
		val = uart.ie

		if write {
			uart.ie = writeData & 3

			uart.sync()
		}

	case 0x14: // ip
		val = uart.ip()

	case 0x120_0004: // PLIC claim
		if uart.ipAndIE() != 0 {
			val = 1
		}
	}

	return val
}

func (uart *UART) AcceptsInput() bool {
	return uart.rx < 0
}

func (uart *UART) SetInput(b byte) {
	uart.rx = int(b)
	uart.sync()
}

func (uart *UART) HasOutput() bool {
	return uart.tx >= 0
}

func (uart *UART) GetOutput() byte {
	b := byte(uart.tx)

	uart.tx = -1
	uart.sync()

	return b
}

func (uart *UART) ip() int {
	tp := -(uart.tx >> 63)
	rp := uart.rx>>63 + 1
	return rp<<1 | tp
}

func (uart *UART) ipAndIE() int {
	return uart.ip() & uart.ie
}

func (uart *UART) sync() {
	if uart.ipAndIE() != 0 {
		uart.cpu.CSR.Mip |= 1 << csr.MipSEIP
	} else {
		uart.cpu.CSR.Mip &^= 1 << csr.MipSEIP
	}
}
