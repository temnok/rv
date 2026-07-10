package state

type TLB struct {
	nodes                  [8]tlbNode
	LookupCount, MissCount int
}

type tlbNode struct {
	vas, pte int
}

func (tlb *TLB) Flush() {
	tlb.nodes[0].vas = 0
}

func (tlb *TLB) Lookup(virtAddr int) (int, int) {
	tlb.LookupCount++

	for i, node := range tlb.nodes {
		if node.vas == 0 {
			break
		}

		shift := node.vas & 0xfff
		if node.vas>>shift != virtAddr>>shift {
			continue
		}

		for j := i; j > 0; j-- {
			tlb.nodes[j] = tlb.nodes[j-1]
		}

		tlb.nodes[0] = node
		return node.pte, shift
	}

	tlb.MissCount++
	return 0, 0
}

func (tlb *TLB) Append(virtAddr, shift, pte int) {
	for i := len(tlb.nodes) - 1; i > 0; i-- {
		tlb.nodes[i] = tlb.nodes[i-1]
	}

	tlb.nodes[0].vas = virtAddr&(-1<<shift) | shift
	tlb.nodes[0].pte = pte
}
