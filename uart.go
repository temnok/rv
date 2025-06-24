package rv

type UART struct {
	baseAddr                         int32
	rx, tx                           UARTfifo
	txctrl, rxctrl, ip, ie, div, clk int32
	callback                         func(ch *byte, write bool) bool
}

type UARTfifo struct {
	buf  uint64
	size int32
}

func (uart *UART) init(callback func(ch *byte, write bool) bool) {
	*uart = UART{
		div:      3,
		callback: callback,
	}
}

func (uart *UART) access(addr int32, data *int32, write bool) bool {
	if addr -= uart.baseAddr; addr < 0 || addr > 8 {
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
			*data = ((8-uart.rx.size)/8)<<31 | uart.rx.get()
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

func (uart *UART) update() bool {
	ch := byte(uart.tx.buf & 0xFF)
	if uart.clk++; uart.clk >= uart.div {
		if bit(uart.txctrl, 1) == 1 && uart.tx.size > 0 && uart.callback(&ch, true) {
			uart.tx.get()
		}

		if bit(uart.rxctrl, 1) == 1 && uart.rx.size < 8 && uart.callback(&ch, false) {
			uart.rx.put(int32(ch))
		}

		uart.clk = 0
	}

	if (uart.txctrl>>16) > uart.tx.size && bit(uart.ie, 1) == 1 {
		uart.ip |= 1
	} else {
		uart.ip &^= 1
	}

	if (uart.rxctrl>>16) > uart.rx.size && bit(uart.ie, 2) == 1 {
		uart.ip |= 2
	} else {
		uart.ip &^= 2
	}

	return uart.ip != 0
}

func (fifo *UARTfifo) put(ch int32) {
	if fifo.size == 8 {
		return
	}

	fifo.buf |= uint64(byte(ch)) << (fifo.size * 8)
	fifo.size++
}

func (fifo *UARTfifo) get() int32 {
	if fifo.size == 0 {
		return 0
	}

	ch := int32(fifo.buf & 0xFF)
	fifo.buf >>= 8
	fifo.size--

	return ch
}
