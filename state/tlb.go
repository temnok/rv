package state

type TLB struct {
	entries [8]tlbEntry
}

type tlbEntry struct {
	vas, pte int
}

const (
	pteG = 5
)

func (tlb *TLB) Flush(leaveGlobals bool) {
	if !leaveGlobals {
		tlb.entries[0].vas = 0
		return
	}

	i := 0
	for _, entry := range tlb.entries {
		if entry.vas == 0 {
			break
		}

		if entry.pte&(1<<pteG) != 0 {
			tlb.entries[i] = entry
			i++
		}
	}

	if i < len(tlb.entries) {
		tlb.entries[i].vas = 0
	}
}

func (tlb *TLB) Lookup(virtAddr int) (int, int) {
	for i, entry := range tlb.entries {
		if entry.vas == 0 {
			break
		}

		shift := entry.vas & 63
		if entry.vas>>shift != virtAddr>>shift {
			continue
		}

		for j := i; j > 0; j-- {
			tlb.entries[j] = tlb.entries[j-1]
		}

		tlb.entries[0] = entry
		return entry.pte, shift
	}

	return 0, 0
}

func (tlb *TLB) Append(virtAddr, shift, pte int) {
	for i := len(tlb.entries) - 1; i > 0; i-- {
		tlb.entries[i] = tlb.entries[i-1]
	}

	tlb.entries[0].vas = virtAddr&(-1<<shift) | shift
	tlb.entries[0].pte = pte
}
