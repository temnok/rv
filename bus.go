package rv

type Bus struct {
	ram []int32
}

const (
	ramBaseAddr      = -1 << 31
	ramBaseAddrWords = 0x2000_0000
)

func (bus *Bus) init(ramSize int) {
	*bus = Bus{
		ram: make([]int32, ramSize/4),
	}
}

func (bus *Bus) read(addr int32, data *int32) bool {
	return bus._access(addr, data, false)
}

func (bus *Bus) write(addr int32, data int32) bool {
	return bus._access(addr, &data, true)
}

func (bus *Bus) _access(addr int32, data *int32, isWrite bool) bool {
	return bus._ramAccess(addr, data, isWrite)
}

func (bus *Bus) _ramAccess(addr int32, data *int32, isWrite bool) bool {
	if addr -= ramBaseAddrWords; addr < 0 || addr >= int32(len(bus.ram)) {
		return false
	}

	if isWrite {
		bus.ram[addr] = *data
	} else {
		*data = bus.ram[addr]
	}

	return true
}
