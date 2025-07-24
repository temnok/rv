package rv

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInstructions32(t *testing.T) {
	testInstructions(t, 32)
}

func TestInstructions64(t *testing.T) {
	testInstructions(t, 64)
}

func testInstructions(t *testing.T, xlen int) {
	matches, err := filepath.Glob(fmt.Sprintf("tests/pass/rv%v*", xlen))
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range matches {
		i := strings.LastIndexByte(file, '/')
		testName := file[i+1:]

		t.Run(testName, func(t *testing.T) {
			runTest(t, xlen, file)
		})
	}
}

func runTest(t *testing.T, xlen int, file string) {
	program, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}

	cpu := &CPU{}
	ram := &RAM{}

	ramBaseAddr := 0x8000_0000

	cpu.Init(xlen, Bus{ram}, ramBaseAddr)
	ram.Init(cpu, ramBaseAddr, 64*1024)
	ram.Load(ramBaseAddr, program)

	instrCounts := make([]int, len(program))
	trapCount := 0
	var lastPCs []uint
	var lastTraps [][2]uint

	for startPC := cpu.PC; ; {
		if i := cpu.PC - startPC; i >= 0 && i < len(instrCounts) {
			instrCounts[i]++
		}

		lastPCs = append(lastPCs, uint(cpu.PC))
		if n := 10; len(lastPCs) == n+1 {
			copy(lastPCs[:n], lastPCs[1:])
			lastPCs = lastPCs[:n]
		}

		if !cpu.Step() {
			debugDump(cpu)
			break
		}

		if cpu.isTrapped() {
			trapCount++

			lastTraps = append(lastTraps, [2]uint{uint(cpu.PC), uint(cpu.Updated.TrapXcause)})
			if n := 10; len(lastTraps) == n+1 {
				copy(lastTraps[:n], lastTraps[1:])
				lastTraps = lastTraps[:n]
			}
		}

		if cpu.CSR.Cycle == 100_000 {
			var addresses []uint
			for i, c := range instrCounts {
				if c > 10_000 {
					addresses = append(addresses, uint(ramBaseAddr+i))
				}
			}

			if n := 10; len(lastTraps) > n {
				lastTraps = lastTraps[len(lastTraps)-n:]
			}

			t.Errorf("timeout: trapCount=%v, priv=%v, mcause=%08x, x31=%08x\n"+
				"last PCs: %x\nlast traps: %x\nloop at addresses: %x\n",
				trapCount, cpu.Priv, cpu.CSR.Mcause, uint(cpu.X[31]), lastPCs, lastTraps, addresses)
			break
		}

		if cpu.isTrapped() {
			if cause := cpu.Updated.TrapXcause; cause == ExceptionEnvironmentCallFromUMode ||
				cause == ExceptionEnvironmentCallFromSMode ||
				cause == ExceptionEnvironmentCallFromMMode {

				if cpu.X[3] == 1 && cpu.X[10] == 0 {
					//fmt.Printf("cycles: %v\n", cpu.CSR.Cycle)
				} else {
					debugDump(cpu)

					t.Errorf("cycles: %v\nlast PCs: %x\nlast traps: %x\n"+
						"priv=%v, pc=%08x\n"+
						"mstatus=%08x, xepc=%08x\n"+
						"xcause=%08x, xtval=%08x\n"+
						"a0=%08x\n",
						cpu.CSR.Cycle, lastPCs, lastTraps,
						cpu.Updated.TrapPriv, uint(cpu.Updated.TrapPC),
						uint(cpu.Updated.TrapMstatus), uint(cpu.Updated.TrapXepc),
						uint(cpu.Updated.TrapXcause), uint(cpu.Updated.TrapXtval),
						uint(cpu.X[10]),
					)
				}

				break
			}
		}
	}
}
