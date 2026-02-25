package state

type TLB struct {
	entries                []TLBEntry
	LookupCount, MissCount int
}

type TLBEntry struct {
	virtAddr, pte int
}

const tlbSize = 16

func (tlb *TLB) Flush() {
	tlb.entries = tlb.entries[:0]
}

func (tlb *TLB) Lookup(virtAddr int) (int, int) {
	tlb.LookupCount++

	for i, e := range tlb.entries {
		shift := e.virtAddr & 0xFFF
		if virtAddr>>shift == e.virtAddr>>shift {
			copy(tlb.entries[1:i+1], tlb.entries[:i])
			tlb.entries[0] = e
			return e.pte, shift
		}
	}

	tlb.MissCount++
	return 0, 0
}

func (tlb *TLB) Append(virtAddr, shift, pte int) {
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
