package uart

import (
	"github.com/temnok/rv/plic"
	"io"
)

type UART struct {
	plic        *plic.PLIC
	interruptID int

	input  chan byte
	output io.Writer

	rx, ip, ie int
}

func New(plic *plic.PLIC, interuptID int, r io.Reader, w io.Writer) *UART {
	input := make(chan byte)

	go func() {
		buf := []byte{0}

		for {
			if n, _ := r.Read(buf); n > 0 {
				input <- buf[0]
			}
		}
	}()

	return &UART{
		plic:        plic,
		interruptID: interuptID,
		input:       input,
		output:      w,

		rx: -1,
	}
}

func (uart *UART) Access(addr int, width int, write bool, writeData int) int {
	addr &= 0xff_ffff

	var val int

	switch addr {
	case 0x00: // txdata
		if write {
			uart.output.Write([]byte{byte(writeData)})
			uart.sync()
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
			uart.ie = writeData

			uart.sync()
		}

	case 0x14: // ip
		val = uart.ip
	}

	return val
}

func (uart *UART) Input() int {
	select {
	case char := <-uart.input:
		uart.rx = int(char)
		uart.sync()

		return int(char)

	default:
		return -1
	}
}

func (uart *UART) sync() {
	if uart.ie&1 == 1 {
		uart.ip |= 1
	} else {
		uart.ip &^= 1
	}

	if uart.rx >= 0 && uart.ie&2 == 2 {
		uart.ip |= 2
	} else {
		uart.ip &^= 2
	}

	uart.plic.PendInterrupt(uart.interruptID, uart.ip != 0)
}
