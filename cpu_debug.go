package rv

import (
	"fmt"
	"github.com/deadsy/rvda"
	"math"
	"strings"
)

var (
	debugTrapCount = 0
	debugTrace     [][]int
)

func (cpu *CPU) debugStep() bool {
	opcode := cpu.innerStep()

	entry := []int{cpu.PC, opcode}

	switch {
	case cpu.Updated.XReg > 0:
		entry = append(entry, cpu.Updated.XVal)
	case cpu.Updated.FReg >= 0:
		entry = append(entry, 0, cpu.Updated.FVal)
	case cpu.Updated.CReg >= 0:
		entry = append(entry, cpu.Updated.CVal)
	}

	debugTrace = append(debugTrace, entry)
	if n := 100; len(debugTrace) == n+1 {
		copy(debugTrace[:n], debugTrace[1:])
		debugTrace = debugTrace[:n]
	}

	//if cpu.isTrapped() {
	//	debugDump(cpu)
	//	return false
	//}

	return true
}

func debugDump(cpu *CPU) {
	isa, _ := rvda.New(uint(cpu.XLen), rvda.RV64gc)

	for _, entry := range debugTrace {
		fmt.Printf("%v\r\n", disassemble(isa, entry))
	}

	fmt.Printf("\r\nCycle: %v, trap: %v\r\n", cpu.CSR.Cycle, debugTrapCount)

	up := &cpu.Updated
	if cpu.Updated.TrapEnter {
		fmt.Printf("\r\nold priv:%x, priv:%x, pc:%x, mstatus:%x\r\n",
			cpu.Priv, uint(up.TrapPriv), uint(up.PC), uint(up.TrapMstatus))
		fmt.Printf("xepc:%x, xcause:%x, xtval:%x\r\n",
			uint(up.TrapXepc), uint(up.TrapXcause), uint(up.TrapXtval))
		fmt.Printf("mtvec:%x, stvec:%x\r\n",
			uint(cpu.CSR.Mtvec), uint(cpu.CSR.Stvec))
	}

	//for i := range 16 {
	//	fmt.Printf("% 5v:%16x      % 5v:%16x\r\n",
	//		regNames[i], uint(cpu.Reg[i]), regNames[16+i], uint(cpu.Reg[16+i]))
	//}
}

func disassemble(isa *rvda.ISA, entry []int) string {
	addr, code := entry[0], entry[1]

	line := isa.Disassemble(uint(addr), uint(code)).String()
	parts := strings.Split(line, "\t")
	ops := strings.Split(parts[1], " ")
	for len(ops) < 2 {
		ops = append(ops, "")
	}

	line = fmt.Sprintf("%-30v %-7v %-16v", parts[0], ops[0], ops[1])

	if len(entry) == 3 {
		line += fmt.Sprintf("// %x", uint(entry[2]))

		if fmt.Sprintf("%x", uint(entry[2])) != fmt.Sprint(entry[2]) {
			line += fmt.Sprintf(" (%v)", entry[2])
		}
	} else if len(entry) == 4 {
		line += fmt.Sprintf("// %x (f32=%v, f64=%v)", uint(entry[3]),
			math.Float32frombits(uint32(entry[3])), math.Float64frombits(uint64(entry[3])))
	}

	return line
}
