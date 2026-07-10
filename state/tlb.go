package state

type TLB struct {
	entries                [8]tlbEntry
	LookupCount, MissCount int
}

type tlbEntry struct {
	vas, pte int
}

func (tlb *TLB) Flush() {
	tlb.entries[0].vas = 0
}

func (tlb *TLB) Lookup(virtAddr int) (int, int) {
	tlb.LookupCount++

	for i, node := range tlb.entries {
		if node.vas == 0 {
			break
		}

		shift := node.vas & 63
		if node.vas>>shift != virtAddr>>shift {
			continue
		}

		for j := i; j > 0; j-- {
			tlb.entries[j] = tlb.entries[j-1]
		}

		tlb.entries[0] = node
		return node.pte, shift
	}

	tlb.MissCount++
	return 0, 0
}

func (tlb *TLB) Append(virtAddr, shift, pte int) {
	for i := len(tlb.entries) - 1; i > 0; i-- {
		tlb.entries[i] = tlb.entries[i-1]
	}

	tlb.entries[0].vas = virtAddr&(-1<<shift) | shift
	tlb.entries[0].pte = pte
}
