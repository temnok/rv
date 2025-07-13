package rv

type UART struct {
	plic        *PLIC
	baseAddr    int
	interruptID int
	callback    func(ch *byte, write bool) bool

	rx, tx                           UARTfifo
	txctrl, rxctrl, ip, ie, div, clk int
}

type UARTfifo struct {
	buf  uint64
	size int
}

func (uart *UART) Init(plic *PLIC, baseAddr int, interuptID int, callback func(ch *byte, write bool) bool) {
	*uart = UART{
		plic:        plic,
		baseAddr:    baseAddr,
		interruptID: interuptID,
		callback:    callback,

		div: 3,
	}
}

func (uart *UART) Access(addr int, data *int, width int, write bool) bool {
	if addr = (addr - uart.baseAddr) / 4; addr < 0 || addr >= 8 || width != 4 {
		return false
	}

	switch addr {
	case 0: // txdata
		if write {
			uart.tx.put(*data)
		} else {
			*data = (uart.tx.size / 8) << 31
		}
	case 1: // rxdata
		if !write {
			size := uart.rx.size
			*data = ((8-size)/8)<<31 | uart.rx.get()
		}
	case 2: // txctrl
		if write {
			uart.txctrl = *data & 0x0007_0003
		} else {
			*data = uart.txctrl
		}
	case 3: // rxctrl
		if write {
			uart.rxctrl = *data & 0x0007_0001
		} else {
			*data = uart.rxctrl
		}
	case 4: // ie
		if write {
			uart.ie = *data & 3
		} else {
			*data = uart.ie
		}
	case 5: // ip
		if !write {
			*data = uart.ip
		}
	case 6: // div
		if write {
			uart.div = *data & 0xFFFF
		} else {
			*data = uart.div
		}
	case 7: // unused
		if !write {
			*data = 0
		}
	}

	return true
}

func (uart *UART) NotifyInterrupts() {
	ch := byte(uart.tx.buf & 0xFF)
	if uart.clk++; uart.clk >= uart.div {
		if bit(uart.txctrl, 0) == 1 && uart.tx.size > 0 {
			if uart.callback != nil && uart.callback(&ch, true) {
				uart.tx.get()
			}
		}

		if bit(uart.rxctrl, 0) == 1 && uart.rx.size < 8 {
			if uart.callback != nil && uart.callback(&ch, false) {
				uart.rx.put(int(ch))
			}
		}

		uart.clk = 0
	}

	if (uart.txctrl>>16) > uart.tx.size && bit(uart.ie, 0) == 1 {
		uart.ip |= 1
	} else {
		uart.ip &^= 1
	}

	if (uart.rxctrl>>16) < uart.rx.size && bit(uart.ie, 1) == 1 {
		uart.ip |= 2
	} else {
		uart.ip &^= 2
	}

	if uart.ip != 0 && uart.plic != nil {
		uart.plic.TriggerInterrupt(uart.interruptID)
	}
}

func (fifo *UARTfifo) put(ch int) {
	if fifo.size == 8 {
		return
	}

	fifo.buf |= uint64(byte(ch)) << (fifo.size * 8)
	fifo.size++
}

func (fifo *UARTfifo) get() int {
	if fifo.size == 0 {
		return 0
	}

	ch := int(fifo.buf & 0xFF)
	fifo.buf >>= 8
	fifo.size--

	return ch
}
