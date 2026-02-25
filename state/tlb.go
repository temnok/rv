package state

type TLB struct {
	nodes                  *tlbNode
	LookupCount, MissCount int
}

type tlbNode struct {
	va, shift, pte int
	next           *tlbNode
}

const tlbSize = 16

func (tlb *TLB) Flush() {
	tlb.nodes = nil
}

func (tlb *TLB) Lookup(virtAddr int) (int, int) {
	tlb.LookupCount++

	for p, n, i := (*tlbNode)(nil), tlb.nodes, 0; n != nil; p, n = n, n.next {
		if n.va == virtAddr&(-1<<n.shift) {
			if p != nil {
				p.next, n.next = n.next, tlb.nodes
				tlb.nodes = n
			}
			return n.pte, n.shift
		}

		if i++; i == tlbSize {
			n.next = nil
			break
		}
	}

	tlb.MissCount++
	return 0, 0
}

func (tlb *TLB) Append(virtAddr, shift, pte int) {
	tlb.nodes = &tlbNode{
		va:    virtAddr & (-1 << shift),
		shift: shift,
		pte:   pte,
		next:  tlb.nodes,
	}
}
