package rv

import (
	"fmt"
	"github.com/deadsy/rvda"
	"strings"
)

var (
	debugTrapCount = 0
	debugTrace     [][4]uint
)

func (cpu *CPU) debugStep() bool {
	pc := cpu.PC
	oldRegs := cpu.Reg

	opcode := cpu.Step()

	entry := [4]uint{uint(pc), uint(opcode)}
	for i, val := range cpu.Reg {
		if val != oldRegs[i] {
			entry[2], entry[3] = uint(i), uint(val)
			break
		}
	}

	debugTrace = append(debugTrace, entry)
	if n := 100; len(debugTrace) == n+1 {
		copy(debugTrace[:n], debugTrace[1:])
		debugTrace = debugTrace[:n]
	}

	if cpu.isTrapped() /*|| cycleCount == 1_000_000*/ {
		debugTrapCount++

		//if trapCount == 8 {
		//	fmt.Printf("Cycle: %v, trap: %v\r\n\r\n", cycleCount, trapCount)
		//
		//	debugDump(cpu)
		//
		//	return false
		//}
	}

	return true
}

func debugDump(cpu *CPU) {
	isa, _ := rvda.New(uint(cpu.XLen), rvda.RV64gc)

	for _, entry := range debugTrace {
		fmt.Printf("%v\r\n", disassemble(isa, entry))
	}

	fmt.Printf("\r\nmcause:%x, mepc:%x, mtvec:%x, mstatus:%x\r\n",
		uint(cpu.CSR.Mcause), uint(cpu.CSR.Mepc), uint(cpu.CSR.Mtvec), uint(cpu.CSR.Mstatus))
	fmt.Printf("scause:%x, sepc:%x, stvec:%x, satp:%x\r\n",
		uint(cpu.CSR.Scause), uint(cpu.CSR.Sepc), uint(cpu.CSR.Stvec), uint(cpu.CSR.Satp))
	fmt.Printf("priv:%v, medeleg:%x\r\n\r\n", cpu.Priv, uint(cpu.CSR.Medeleg))

	for i := range 16 {
		fmt.Printf("% 5v:%16x      % 5v:%16x\r\n",
			regNames[i], uint(cpu.Reg[i]), regNames[16+i], uint(cpu.Reg[16+i]))
	}
}

func disassemble(isa *rvda.ISA, entry [4]uint) string {
	addr, code, reg, regVal := entry[0], entry[1], entry[2], entry[3]

	line := isa.Disassemble(addr, code).String()
	parts := strings.Split(line, "\t")
	ops := strings.Split(parts[1], " ")
	for len(ops) < 2 {
		ops = append(ops, "")
	}

	line = fmt.Sprintf("%-30v %-6v %-16v", parts[0], ops[0], ops[1])

	if entry[2] != 0 {
		line += fmt.Sprintf("// %v=%x", regNames[reg], regVal)
	}

	return line
}

var regNames = []string{
	"zero", "ra", "sp", "gp", "tp", "t0", "t1", "t2",
	"s0", "s1", "a0", "a1", "a2", "a3", "a4", "a5",
	"a6", "a7", "s2", "s3", "s4", "s5", "s6", "s7",
	"s8", "s9", "s10", "s11", "t3", "t4", "t5", "t6",
}
