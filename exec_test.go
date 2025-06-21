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

		program, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}

		cpu := &CPU{
			mem: make([]byte, 0x10000),
			pc:  ramBaseAddr,
		}

		copy(cpu.mem, program)

		count := 0
		for !cpu.faulted() {
			cpu.execInstr()
			count++
		}

		if !(cpu.eCall && cpu.x[3] == 1 && cpu.x[10] == 0) {
			t.Fatalf("Test %v: count=%v, instrIllegal=%v, memAccessFault=%v, pc=%08x, x3=%08x, x10=%08x\n",
				testName, count, cpu.instrIllegal, cpu.memAccessFault,
				uint32(cpu.pc), uint32(cpu.x[3]), uint32(cpu.x[10]))
		}

		fmt.Printf("%v: OK, count=%v\n", testName, count)
	}
}
