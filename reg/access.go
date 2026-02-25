package reg

func Access(reg, data *int, offset, width int, isWrite bool) {
	if isWrite {
		write(reg, *data, offset, width)
	} else {
		read(*reg, data, offset, width)
	}
}

func read(reg int, data *int, offset, width int) {
	if width == 8 {
		*data = reg
		return
	}

	*data = (reg >> (offset * 8)) &^ (-1 << (width * 8))
}

func write(reg *int, data int, offset, width int) {
	if width == 8 {
		*reg = data
		return
	}

	mask := (1<<(width*8) - 1) << (offset * 8)
	*reg = *reg&^mask | (data<<(offset*8))&mask
}
