package uart

import (
	"github.com/temnok/rv/plic"
	"io"
)

type UART struct {
	plic        *plic.PLIC
	baseAddr    int
	interruptID int

	input  chan byte
	output io.Writer

	rx, ip, ie int
}

func New(plic *plic.PLIC, baseAddr int, interuptID int, r io.Reader, w io.Writer) *UART {
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
		baseAddr:    baseAddr,
		interruptID: interuptID,
		input:       input,
		output:      w,

		rx: -1,
	}
}

func (uart *UART) Access(addr int, data *int, width int, write bool) bool {
	if addr = addr - uart.baseAddr; addr < 0 || addr >= 32 || width != 4 {
		return false
	}

	switch addr {
	case 0x00: // txdata
		if write {
			uart.output.Write([]byte{byte(*data)})
			uart.sync()
		}

	case 0x04: // rxdata
		if !write {
			if uart.rx < 0 {
				*data = 1 << 31
			} else {
				*data = uart.rx
				uart.rx = -1

				uart.sync()
			}
		}

	case 0x10: // ie
		if write {
			uart.ie = *data

			uart.sync()
		} else {
			*data = uart.ie
		}

	case 0x14: // ip
		if !write {
			*data = uart.ip
		}
	}

	return true
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
