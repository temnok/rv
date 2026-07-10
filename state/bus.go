package state

type Bus []func(addr int, data *int, width int, write bool) bool

func (bus Bus) Read(addr int, width int) (int, bool) {
	var data int
	ok := bus.Access(addr, &data, width, false)
	return data, ok
}

func (bus Bus) Write(addr int, data int, width int) bool {
	return bus.Access(addr, &data, width, true)
}

func (bus Bus) Access(addr int, data *int, width int, write bool) bool {
	for _, device := range bus {
		if device(addr, data, width, write) {
			return true
		}
	}

	return false
}
