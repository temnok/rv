package rv

type PLIC struct {
	cpu      *CPU
	baseAddr int

	priority  [32]int
	pending   int
	enable    int
	threshold int
	claim     int
	claiming  int
}

func (plic *PLIC) Init(cpu *CPU, baseAddr int) {
	*plic = PLIC{
		cpu:      cpu,
		baseAddr: Xint(baseAddr),
	}
}

func (plic *PLIC) access(addr int, data *int, width int, write bool) bool {
	if addr = (addr - plic.baseAddr) / 4; addr < 0 || addr >= 0x400_0000/4 || width < 4 {
		return false
	}

	var reg *int
	val := *data

	switch addr {
	case 0x1000 / 4:
		reg = &plic.pending
		val &^= 1

	case 0x2000 / 4:
		reg = &plic.enable
		val &^= 1

	case 0x200000 / 4:
		reg = &plic.threshold

	case 0x200004 / 4:
		reg = &plic.claim

		if write {
			if val > 0 && val < 32 {
				plic.claiming &^= 1 << val
			}
		} else {
			if plic.claim > 0 && plic.claim < 32 && bit(plic.pending, plic.claim) == 1 {
				plic.claiming |= 1 << plic.claim
			}
		}
	default:
		if addr > 0 && addr < len(plic.priority) {
			reg = &plic.priority[addr]
		}
	}

	if write {
		if reg != nil {
			*reg = val
		}
	} else {
		if reg != nil {
			*data = *reg
		} else {
			*data = 0
		}
	}

	return true
}

func (plic *PLIC) triggerInterrupt(source int) {
	if source > 0 && source < 32 && bit(plic.claiming, source) == 0 && bit(plic.pending, source) == 0 {
		plic.pending |= 1 << source
	}
}

func (plic *PLIC) notifyInterrupts() {
	maxPriority := 0
	irq := 0
	for i := 1; i < 32; i++ {
		if bit(plic.claiming, i) == 1 {
			plic.pending &^= 1 << i
		} else if bit(plic.enable, i) == 1 &&
			bit(plic.pending, i) == 1 &&
			plic.priority[i] >= maxPriority &&
			plic.priority[i] >= plic.threshold {

			irq = i
			maxPriority = plic.priority[i]
		}
	}

	plic.claim = irq

	if irq > 0 {
		plic.cpu.csr.mip |= 1 << mipSEI
	}
}
