package dev

import "github.com/temnok/rv/state"

func Process(cpu *state.CPU) {
	processTimer(cpu)
	processUartInput(cpu)
	processUartOutput(cpu)
}
