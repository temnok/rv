package rv

import (
	"fmt"
	"github.com/deadsy/rvda"
	"strings"
)

var (
	cycleCount = 0
	trapCount  = 0
	trace      [][4]uint
)

func (cpu *CPU) debugStep() {
	cycleCount++

	pc := cpu.pc
	oldRegs := cpu.x

	opcode := cpu.step()

	entry := [4]uint{uint(pc), uint(opcode)}
	for i, val := range cpu.x {
		if val != oldRegs[i] {
			entry[2], entry[3] = uint(i), uint(val)
			break
		}
	}

	trace = append(trace, entry)
	if n := 1000; len(trace) == n+1 {
		copy(trace[:n], trace[1:])
		trace = trace[:n]
	}

	if cpu.isTrapped /*|| cycleCount == 1_000_000*/ {
		trapCount++

		if trapCount == 1 {
			fmt.Printf("Cycle: %v, trap: %v\r\n\r\n", cycleCount, trapCount)

			debugDump(cpu)

			panic("stop")
		}
	}
}

var isa, _ = rvda.New(64, rvda.RV64gc)

func debugDump(cpu *CPU) {
	for _, entry := range trace {
		line := isa.Disassemble(entry[0], entry[1]).String()
		i := strings.LastIndexByte(line, ' ')
		line = line[:i] + "\t" + line[i+1:]

		if entry[2] != 0 {
			line += fmt.Sprintf("\t\t// %v:%x", regNames[entry[2]], entry[3])
		}

		fmt.Printf("%v\r\n", line)
	}

	fmt.Printf("\r\nmcause:%x, mepc:%x, mstatus:%x, satp:%x\r\n\r\n",
		cpu.csr.mcause, cpu.csr.mepc, cpu.csr.mstatus, cpu.csr.satp)

	for i := range 16 {
		fmt.Printf("% 5v:%16x      % 5v:%16x\r\n", regNames[i], cpu.x[i], regNames[16+i], cpu.x[16+i])
	}
}

var regNames = []string{
	"zero", "ra", "sp", "gp", "tp", "t0", "t1", "t2",
	"s0", "s1", "a0", "a1", "a2", "a3", "a4", "a5",
	"a6", "a7", "s2", "s3", "s4", "s5", "s6", "s7",
	"s8", "s9", "s10", "s11", "t3", "t4", "t5", "t6",
}
