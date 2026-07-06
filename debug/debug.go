package debug

import (
	"fmt"
	"github.com/deadsy/rvda"
	"github.com/temnok/rv/arch"
	cp "github.com/temnok/rv/cpu"
	"github.com/temnok/rv/state"
	"math"
	"strings"
)

var (
	debugTrapCount = 0
	debugTrace     [][]int
)

func Step(cpu *state.CPU) bool {
	opcode := cp.InnerStep(cpu)

	entry := []int{cpu.PC, opcode}

	switch {
	case cpu.Update.XReg > 0:
		entry = append(entry, cpu.Update.XVal)
	case cpu.Update.FReg >= 0:
		entry = append(entry, 0, cpu.Update.FVal)
	case cpu.Update.CReg >= 0:
		entry = append(entry, cpu.Update.CVal)
	}

	debugTrace = append(debugTrace, entry)
	if n := 100; len(debugTrace) == n+1 {
		copy(debugTrace[:n], debugTrace[1:])
		debugTrace = debugTrace[:n]
	}

	//if trap.IsEntered(cpu.CPU) {
	//	Dump(cpu)
	//	return false
	//}

	return true
}

func Dump(cpu *state.CPU) {
	isa, _ := rvda.New(uint(arch.XLen), rvda.RV64gc)

	for _, entry := range debugTrace {
		fmt.Printf("%v\r\n", disassemble(isa, entry))
	}

	fmt.Printf("\r\nCycle: %v, Trap: %v\r\n", cpu.CSR.Cycle, debugTrapCount)

	up := &cpu.Update
	//if cpu.Update.TrapEnter {
	fmt.Printf("\r\nold priv:%x, priv:%x, pc:%x, mstatus:%x\r\n",
		cpu.Priv, uint(up.TrapPriv), uint(up.PC), uint(up.TrapMstatus))
	fmt.Printf("xepc:%x, xcause:%x, xtval:%x\r\n",
		uint(up.TrapXepc), uint(up.TrapXcause), uint(up.TrapXtval))
	fmt.Printf("mtvec:%x, stvec:%x\r\n",
		uint(cpu.CSR.Mtvec), uint(cpu.CSR.Stvec))
	//}

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
