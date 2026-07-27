package tlb

import "github.com/temnok/rv/state"

func Lookup(cpu *state.CPU, va int) (int, int) {
	for i, entry := range cpu.TLB {
		if entry.Vas == 0 {
			break
		}

		if shift := entry.Vas & 63; entry.Vas>>shift == va>>shift {
			if i > 0 {
				cpu.Update.Targets |= state.UpdateTLB

				for j := i; j > 0; j-- {
					cpu.Update.TLB[j] = cpu.TLB[j-1]
				}

				cpu.Update.TLB[0] = entry

			}

			return entry.Pte, shift
		}
	}

	return 0, 0
}
