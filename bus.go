package rv

type Bus struct {
	ram   *RAM
	clint *CLINT
}

func (bus *Bus) read(addr int32, data *int32) bool {
	return bus.access(addr, data, false)
}

func (bus *Bus) write(addr int32, data int32) bool {
	return bus.access(addr, &data, true)
}

func (bus *Bus) access(addr int32, data *int32, write bool) bool {
	return bus.ram != nil && bus.ram.access(addr, data, write) ||
		bus.clint != nil && bus.clint.access(addr, data, write)
}
