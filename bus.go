package rv

type Bus []interface {
	access(addr int32, data *int32, write bool) bool
}

func (bus Bus) read(addr int32, data *int32) bool {
	return bus.access(addr, data, false)
}

func (bus Bus) write(addr int32, data int32) bool {
	return bus.access(addr, &data, true)
}

func (bus Bus) access(addr int32, data *int32, write bool) bool {
	for _, device := range bus {
		if device.access(addr, data, write) {
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
