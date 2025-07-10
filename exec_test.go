package rv

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInstructions(t *testing.T) {
	matches, err := filepath.Glob("tests/pass/*")
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

	cpu := &CPU{}
	ram := &RAM{}

	ramBaseAddrUint := uint(0x8000_0000)
	ramBaseAddr := int(ramBaseAddrUint)

	cpu.Init(Bus{ram}, ramBaseAddr, nil)
	ram.Init(ramBaseAddr, 64*1024)
	ram.Load(ramBaseAddr, program)

	instrCounts := make([]int, len(program))
	trapCount := 0
	var lastPCs []uint
	var lastTraps [][2]uint

	for {
		if i := cpu.pc - ramBaseAddr; i >= 0 && i < len(instrCounts) {
			instrCounts[i]++
		}

		lastPCs = append(lastPCs, uint(cpu.pc))
		if n := 10; len(lastPCs) == n+1 {
			copy(lastPCs[:n], lastPCs[1:])
			lastPCs = lastPCs[:n]
		}

		prevPC := cpu.pc
		cpu.Step()

		if cpu.isTrapped {
			trapCount++

			lastTraps = append(lastTraps, [2]uint{uint(prevPC), uint(cpu.csr.mcause)})
			if n := 10; len(lastTraps) == n+1 {
				copy(lastTraps[:n], lastTraps[1:])
				lastTraps = lastTraps[:n]
			}
		}

		if cpu.csr.cycle == 100_000 {
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
				trapCount, cpu.priv, cpu.csr.mcause, uint(cpu.x[31]), lastPCs, lastTraps, addresses)
			break
		}

		if cpu.isTrapped {
			if cause := cpu.csr.mcause; cause == ExceptionEnvironmentCallFromUMode ||
				cause == ExceptionEnvironmentCallFromSMode ||
				cause == ExceptionEnvironmentCallFromMMode {

				if cpu.x[3] == 1 && cpu.x[10] == 0 {
					fmt.Printf("cycles: %v\n", cpu.csr.cycle)
				} else {
					t.Errorf("cycles: %v\nlast PCs: %x\nlast traps: %x\n"+
						"priv=%v, cause=%v,  mepc=%08x, mstatus=%08x, "+
						"gp=%08x, a0=%08x, t0=%08x, t2=%08x, t6=%08x\n",
						cpu.csr.cycle, lastPCs, lastTraps,
						cpu.priv, cause, uint(cpu.csr.mepc), uint(cpu.csr.mstatus),
						uint(cpu.x[3]), uint(cpu.x[10]),
						uint(cpu.x[5]), uint(cpu.x[7]), uint(cpu.x[31]))
				}

				break
			}
		}
	}
}
