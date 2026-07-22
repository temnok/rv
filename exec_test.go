package rv

import (
	cp "github.com/temnok/rv/cpu"
	"github.com/temnok/rv/debug"
	"github.com/temnok/rv/ram"
	"github.com/temnok/rv/trap"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInstructions(t *testing.T) {
	matches, err := filepath.Glob("tests/pass/rv64*")
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range matches {
		i := strings.LastIndexByte(file, '/')
		testName := file[i+1:]

		t.Run(testName, func(t *testing.T) {
			runTest(t, file)
		})
	}
}

func runTest(t *testing.T, file string) {
	program, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}

	ramBaseAddr := 0x8000_0000
	cpu := cp.New(ramBaseAddr)
	ram := ram.New(64 * 1024)
	ram.Load(ramBaseAddr, program)

	cpu.RAM = ram.Access

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

		isTrap := cp.Step(cpu) < 0

		if isTrap {
			trapCount++

			lastTraps = append(lastTraps, [2]uint{uint(cpu.PC), uint(cpu.Update.TrapXcause)})
			if n := 10; len(lastTraps) == n+1 {
				copy(lastTraps[:n], lastTraps[1:])
				lastTraps = lastTraps[:n]
			}
		}

		if cpu.CSR.Mcycle == 100_000 {
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

		if isTrap {
			if cause := cpu.Update.TrapXcause; cause == trap.EnvironmentCallFromUMode ||
				cause == trap.EnvironmentCallFromSMode ||
				cause == trap.EnvironmentCallFromMMode {

				if cpu.X[3] == 1 && cpu.X[10] == 0 {
					//fmt.Printf("cycles: %v\n", cpu.CSR.Mcycle)
				} else {
					debug.Dump(cpu)

					t.Errorf("cycles: %v\nlast PCs: %x\nlast traps: %x\n"+
						"priv=%v, pc=%08x\n"+
						"mstatus=%08x\n"+
						"xcause=%08x, xtval=%08x\n"+
						"a0=%08x\n",
						cpu.CSR.Mcycle, lastPCs, lastTraps,
						cpu.Update.Priv, uint(cpu.Update.PC),
						uint(cpu.Update.Mstatus),
						uint(cpu.Update.TrapXcause), uint(cpu.Update.TrapXtval),
						uint(cpu.X[10]),
					)
				}

				break
			}
		}
	}
}
