package rv

type TLB struct {
	entries []TLBEntry
}

type TLBEntry struct {
	virtAddr, pte Xint
}

const tlbSize = 8

func (tlb *TLB) flush() {
	tlb.entries = tlb.entries[:0]
}

func (tlb *TLB) lookup(virtAddr Xint) (Xint, Xint) {
	for i, e := range tlb.entries {
		shift := e.virtAddr & 0xFFF
		if virtAddr>>shift == e.virtAddr>>shift {
			copy(tlb.entries[1:i+1], tlb.entries[:i])
			tlb.entries[0] = e
			return e.pte, shift
		}
	}

	return 0, 0
}

func (tlb *TLB) append(virtAddr, shift, pte Xint) {
	if tlbSize == 0 {
		return
	}

	tlb.entries = append(tlb.entries, TLBEntry{})
	n := len(tlb.entries)
	copy(tlb.entries[1:n], tlb.entries[:n-1])
	tlb.entries[0] = TLBEntry{
		virtAddr: virtAddr>>shift<<shift | shift,
		pte:      pte,
	}

	if n > tlbSize {
		tlb.entries = tlb.entries[:tlbSize]
	}
}
