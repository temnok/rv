package rv

type Bus []interface {
	access(addr int, data *int, width int, write bool) bool
}

func (bus Bus) read(addr int, data *int, width int) bool {
	return bus.access(addr, data, width, false)
}

func (bus Bus) write(addr int, data int, width int) bool {
	return bus.access(addr, &data, width, true)
}

func (bus Bus) access(addr int, data *int, width int, write bool) bool {
	for _, device := range bus {
		if device.access(addr, data, width, write) {
			return true
		}
	}

	return false
}

func (bus Bus) notifyInterrupts() {
	for _, device := range bus {
		if target, ok := device.(interface{ notifyInterrupts() }); ok {
			target.notifyInterrupts()
		}
	}
}
