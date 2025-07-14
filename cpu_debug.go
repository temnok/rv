package rv

import (
	"fmt"
	"github.com/deadsy/rvda"
	"strings"
)

var (
	debugTrapCount = 0
	debugTrace     [][]int
)

func (cpu *CPU) debugStep() bool {
	pc := cpu.Updated.PC
	opcode := cpu.innerStep()

	entry := []int{pc, opcode}

	if cpu.Updated.RegIndex != 0 {
		entry = append(entry, cpu.Updated.RegValue)
	}

	if cpu.Updated.CSRIndex != 0 {
		entry = append(entry, cpu.Updated.CSRValue)
	}

	debugTrace = append(debugTrace, entry)
	if n := 100; len(debugTrace) == n+1 {
		copy(debugTrace[:n], debugTrace[1:])
		debugTrace = debugTrace[:n]
	}

	if cpu.isTrapped() {
		//if cpu.CSR.Cycle == 10 {
		debugTrapCount++

		if debugTrapCount == 1 {
			fmt.Printf("Cycle: %v, trap: %v\r\n\r\n", cpu.CSR.Cycle, debugTrapCount)

			debugDump(cpu)

			return false
		}
	}

	return true
}

func debugDump(cpu *CPU) {
	isa, _ := rvda.New(uint(cpu.XLen), rvda.RV64gc)

	for _, entry := range debugTrace {
		fmt.Printf("%v\r\n", disassemble(isa, entry))
	}

	up := &cpu.Updated
	if cpu.Updated.TrapEnter {
		fmt.Printf("\r\npriv:%x, pc:%x, mstatus:%x\r\n",
			uint(up.TrapPriv), uint(up.PC), uint(up.TrapMstatus))
		fmt.Printf("xepc:%x, xcause:%x, xtval:%x\r\n",
			uint(up.TrapXepc), uint(up.TrapXcause), uint(up.TrapXtval))
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

	line = fmt.Sprintf("%-30v %-6v %-16v", parts[0], ops[0], ops[1])

	if len(entry) > 2 {
		line += fmt.Sprintf("// %x", entry[2])
	}

	return line
}
