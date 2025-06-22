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

	cpu := &CPU{
		mem: make([]byte, 0x10000),
		pc:  ramBaseAddr,
	}

	copy(cpu.mem, program)

	for count := 1; ; count++ {
		if count == 10_000 {
			t.Errorf("timeout")
			return
		}

		cpu.execInstr()

		if cpu.trapped && (cpu.csr.mcause == ExceptionEnvironmentCallFromUMode ||
			cpu.csr.mcause == ExceptionEnvironmentCallFromSMode ||
			cpu.csr.mcause == ExceptionEnvironmentCallFromMMode) {

			if cpu.x[3] == 1 && cpu.x[10] == 0 {
				fmt.Printf("count=%v\n", count)
			} else {
				t.Errorf("count=%v, cause=%v, pc=%08x, x3=%08x, x10=%08x\n",
					count, cpu.csr.mcause, uint32(cpu.pc), uint32(cpu.x[3]), uint32(cpu.x[10]))
			}

			return
		}
	}

}
