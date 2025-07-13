package rv

type Bus []interface {
	Access(addr int, data *int, width int, write bool) bool
}

func (bus Bus) Read(addr int, data *int, width int) bool {
	return bus.Access(addr, data, width, false)
}

func (bus Bus) Write(addr int, data int, width int) bool {
	return bus.Access(addr, &data, width, true)
}

func (bus Bus) Access(addr int, data *int, width int, write bool) bool {
	for _, device := range bus {
		if device.Access(addr, data, width, write) {
			return true
		}
	}

	return false
}

func (bus Bus) NotifyInterrupts() {
	for _, device := range bus {
		if target, ok := device.(interface{ NotifyInterrupts() }); ok {
			target.NotifyInterrupts()
		}
	}
}
