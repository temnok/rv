package instr

import "github.com/temnok/rv/state"

type instr = func(*state.CPU, Op)
