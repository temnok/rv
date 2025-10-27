package rv

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

func (cpu *CPU) translateSv(virtAddr int, physAddr *int, access int) {
	if cpu.Xlen64() {
		cpu.translateSv39(virtAddr, physAddr, access)
	} else {
		cpu.translateSv32(virtAddr, physAddr, access)
	}
}
