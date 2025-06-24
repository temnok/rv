package rv

type PLIC struct {
	cpu      *CPU
	baseAddr int32

	priority  [32]int32
	pending   int32
	enable    int32
	threshold int32
	claim     int32
	claiming  int32

	AccessCount int
}

func (plic *PLIC) Init(cpu *CPU, baseAddr int32) {
	*plic = PLIC{
		cpu:      cpu,
		baseAddr: int32(uint32(baseAddr) >> 2),
	}
}

func (plic *PLIC) access(addr int32, data *int32, write bool) bool {
	if addr -= plic.baseAddr; addr < 0 || addr >= 0x400_0000/4 {
		return false
	}

	var reg *int32
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
		if addr > 0 && addr < int32(len(plic.priority)) {
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

	plic.AccessCount++
	return true
}

func (plic *PLIC) triggerInterrupt(source int32) {
	if source > 0 && source < 32 {
		plic.pending |= 1 << source
	}
}

func (plic *PLIC) notifyInterrupts() bool {
	maxPriority := int32(0)
	irq := int32(0)
	for i := int32(1); i < 32; i++ {
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
		return true
	}

	return false
}
