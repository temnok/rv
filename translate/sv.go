package translate

import "github.com/temnok/rv/state"

const (
	// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#translation
	PteV = 0
	PteR = 1
	PteW = 2
	PteX = 3
	PteU = 4
	//PteG = 5
	PteA = 6
	PteD = 7
)

const (
	PrivU = 0
	PrivS = 1
	PrivM = 3

	AccessExecute = 0
	AccessRead    = 1
	AccessWrite   = 3
)

func Sv(cpu *state.CPU, virtAddr int, physAddr *int, access int) {
	sv39(cpu, virtAddr, physAddr, access)
}
