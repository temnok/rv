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
	cpu.init(64 * 1024)

	copy(cpu.mem, program)

	for count := 1; ; count++ {
		if count == 100_000 {
			t.Errorf("timeout")
			return
		}

		cpu.step()

		if cpu.trapped && (cpu.csr.mcause == ExceptionEnvironmentCallFromUMode ||
			cpu.csr.mcause == ExceptionEnvironmentCallFromSMode ||
			cpu.csr.mcause == ExceptionEnvironmentCallFromMMode) {

			if cpu.x[3] == 1 && cpu.x[10] == 0 {
				fmt.Printf("instructions: %v\n", count)
			} else {
				t.Errorf("instructions: %v, cause: %v, pc: %08x, mepc: %08x, x3: %08x, x10: %08x\n",
					count, cpu.csr.mcause, uint32(cpu.pc), uint32(cpu.csr.mepc),
					uint32(cpu.x[3]), uint32(cpu.x[10]))
			}

			break
		}
	}
}
