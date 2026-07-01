package uart

import (
	"github.com/temnok/rv/plic"
	"github.com/temnok/rv/terminal"
)

type UART struct {
	plic        *plic.PLIC
	baseAddr    int
	interruptID int

	tx, rx, ip, ie int
}

func New(plic *plic.PLIC, baseAddr int, interuptID int) *UART {
	return &UART{
		plic:        plic,
		baseAddr:    baseAddr,
		interruptID: interuptID,

		tx: -1,
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
			uart.tx = *data & 0xFF

			uart.sync()
		} else {
			if uart.tx >= 0 {
				*data = 1 << 31
			} else {
				*data = 0
			}
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

func (uart *UART) Sync(term *terminal.Terminal) {
	sync := false

	if char, ok := term.GetChar(); ok {
		uart.rx = int(char)
		sync = true
	}

	if uart.tx >= 0 {
		term.PutChar(byte(uart.tx))
		uart.tx = -1
		sync = true
	}

	if sync {
		uart.sync()
	}
}

func (uart *UART) sync() {
	if uart.tx < 0 && uart.ie&1 == 1 {
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
