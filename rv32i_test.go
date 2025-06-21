package rv

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRV32(t *testing.T) {
	matches, err := filepath.Glob("tests/pass/*")
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range matches {
		i := strings.LastIndexByte(file, '/')
		testName := file[i+1:]

		content, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}

		cpu := &CPU{
			ram: make([]byte, 0x10000),
			pc:  ramBase,
		}

		copy(cpu.ram, content)
		count := 0
		for !cpu.faulted() {
			cpu.executeInstr()
			count++
		}

		if !(cpu.eCall && cpu.x[3] == 1 && cpu.x[10] == 0) {
			t.Fatalf("Test %v: count=%v, instrIllegal=%v, accessFault=%v, pc=%08x, x3=%08x, x10=%08x\n",
				testName, count, cpu.instrIllegal, cpu.accessFault,
				uint32(cpu.pc), uint32(cpu.x[3]), uint32(cpu.x[10]))
		}

		fmt.Printf("%v: OK, count=%v\n", testName, count)
	}
}
