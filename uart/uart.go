package uart

import (
	"github.com/temnok/rv/plic"
)

type UART struct {
	plic        *plic.PLIC
	interruptID int

	tx, rx, ie int
}

func New(plic *plic.PLIC, interuptID int) *UART {
	return &UART{
		plic:        plic,
		interruptID: interuptID,

		tx: -1,
		rx: -1,
	}
}

func (uart *UART) Access(addr int, width int, write bool, writeData int) int {
	addr &= 0xff_ffff

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
	return (rp<<1 | tp) & uart.ie
}

func (uart *UART) sync() {
	uart.plic.PendInterrupt(uart.interruptID, uart.ip() != 0)
}
