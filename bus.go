package rv

type Bus struct {
	ram []byte
}

const (
	ramBaseAddr = -1 << 31
)

func (bus *Bus) init(ramSize int) {
	*bus = Bus{
		ram: make([]byte, ramSize),
	}
}

func (bus *Bus) read(addr, width int32, data *int32) bool {
	return bus._access(addr, width, data, false)
}

func (bus *Bus) write(addr, width int32, data int32) bool {
	return bus._access(addr, width, &data, true)
}

func (bus *Bus) _access(addr, width int32, data *int32, isWrite bool) bool {
	return bus._ramAccess(addr, width, data, isWrite)
}

func (bus *Bus) _ramAccess(addr, width int32, data *int32, isWrite bool) bool {
	if addr -= ramBaseAddr; addr < 0 || addr+width > int32(len(bus.ram)) {
		return false
	}

	if ram := bus.ram[addr : addr+width]; isWrite {
		val := *data
		for i := range ram {
			ram[i] = byte(val >> (i << 3))
		}
	} else {
		val := int32(0)
		for i, b := range ram {
			val |= int32(b) << (i << 3)
		}
		*data = val
	}

	return true
}
