package rv

type Bus []interface {
	access(addr Xint, data *Xint, width Xint, write bool) bool
}

func (bus Bus) read(addr Xint, data *Xint, width Xint) bool {
	return bus.access(addr, data, width, false)
}

func (bus Bus) write(addr Xint, data Xint, width Xint) bool {
	return bus.access(addr, &data, width, true)
}

func (bus Bus) access(addr Xint, data *Xint, width Xint, write bool) bool {
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
