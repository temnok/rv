package ins

import (
	"github.com/temnok/rv/state"
)

func sb(cpu *state.CPU, op Op) {
	store(cpu, op, 1)
}
