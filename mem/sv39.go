package mem

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/ram"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

const (
	accessFetch = 0
	accessLoad  = 1
	accessStore = 3

	// https://docs.riscv.org/reference/isa/v20260120/priv/supervisor.html#sv39pte
	pteV = 0
	pteR = 1
	pteW = 2
	pteX = 3
	pteU = 4
	pteG = 5
	pteA = 6
	pteD = 7

	leafMask = 1<<pteR | 1<<pteW | 1<<pteX
)

func translateSv39(cpu *state.CPU, virtAddr int, access int) int {
	// https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#3-1-1-6-4-memory-privilege-in-mstatus-register
	epriv := cpu.CSR.Priv
	if cpu.CSR.Mstatus>>csr.MstatusMPRV&1 == 1 && access != accessFetch {
		epriv = cpu.CSR.Mstatus >> csr.MstatusMPP & 3
	}

	// https://docs.riscv.org/reference/isa/v20260120/priv/supervisor.html#norm:satp_op_active
	// https://docs.riscv.org/reference/isa/v20260120/priv/supervisor.html#norm:satp-mode
	if epriv == csr.PrivM || cpu.CSR.Satp == 0 {
		return virtAddr
	}

	// https://docs.riscv.org/reference/isa/v20260120/priv/supervisor.html#addressing-and-memory-protection
	//if upper := virtAddr >> 38; upper != 0 && upper != -1 {
	//	trap.Enter(cpu, trap.PageFault+access, virtAddr)
	//	return 0
	//}

	pte, shift := loadPTEsv39(cpu, virtAddr, access)
	if trap.IsEntered(cpu) {
		return 0
	}

	sum, mxr := cpu.CSR.Mstatus>>csr.MstatusSUM&1, cpu.CSR.Mstatus>>csr.MstatusMXR&1

	if epriv == csr.PrivU && pte>>pteU&1 == 0 ||
		epriv == csr.PrivS && pte>>pteU&1 == 1 && !(sum == 1 && access != accessFetch) ||
		access == accessFetch && pte>>pteX&1 == 0 ||
		access == accessLoad && pte>>pteR&1 == 0 && !(mxr == 1 && pte>>pteX&1 == 1) ||
		access == accessStore && !(pte>>pteW&1 == 1 && pte>>pteD&1 == 1) ||
		pte>>pteA&1 == 0 {

		trap.Enter(cpu, trap.PageFault+access, virtAddr)
		return 0
	}

	// https://docs.riscv.org/reference/isa/v20260120/priv/supervisor.html#sv39pte
	return pte>>10&^(-1<<44)<<12 | virtAddr&^(-1<<shift)
}

// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#sv32algorithm
func loadPTEsv39(cpu *state.CPU, virtAddr, access int) (targetPTE, shift int) {
	pte := loadPTE(cpu, cpu.CSR.Satp, virtAddr, 30, access)
	if trap.IsEntered(cpu) || pte&leafMask != 0 {
		return pte, 30
	}

	pte = loadPTE(cpu, pte>>10, virtAddr, 21, access)
	if trap.IsEntered(cpu) || pte&leafMask != 0 {
		return pte, 21
	}

	return loadPTE(cpu, pte>>10, virtAddr, 12, access), 12
}

func loadPTE(cpu *state.CPU, ptNum, virtAddr, shift, access int) int {
	pteAddr := ptNum&^(-1<<44)<<12 | (virtAddr>>shift&511)<<3

	pte := ram.Read8(cpu, pteAddr)

	isLeaf := shift == 12 || pte&leafMask != 0

	if pte>>pteV&1 == 0 || // valid bit not set
		pte>>pteR&1 == 0 && pte>>pteW&1 == 1 || // reserved
		isLeaf && pte>>10&^(-1<<(shift-12)) != 0 { // misaligned page

		trap.Enter(cpu, trap.PageFault+access, virtAddr)
		return 0
	}

	return pte
}
