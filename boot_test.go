package rv

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestBoot(t *testing.T) {
	testBoot(t, "biko.gz", 50_000_000)
}

// TODO: enable when faster
func xTestBootGo(t *testing.T) {
	testBoot(t, "biko-go.gz", 1_000_000_000)
}

func testBoot(t *testing.T, kernelFileName string, timeout int) {
	t.Parallel()

	kernelPath := "biko/output/" + kernelFileName
	stopString := "user@rv"
	inR, inW := io.Pipe()

	outW := &matchWriter{
		t:          t,
		stopString: stopString,
		input:      inW,
	}

	cpu, err := bootLinux(kernelPath, inR, outW, timeout)
	if err != nil {
		t.Fatalf("Boot error: %v", err)
	}

	if !outW.success {
		t.Fatalf("Expected '%v' stop string, got:\n%v", stopString, string(outW.output))
	}

	fmt.Printf("Instr/CInstr: %v/%v, traps: %v, TLB lookups/misses: %v/%v\n",
		cpu.InstrCount, cpu.CInstrCount, cpu.TrapCount, cpu.TLB.LookupCount, cpu.TLB.MissCount)
}

type matchWriter struct {
	t          *testing.T
	stopString string
	input      io.Writer
	output     []byte
	success    bool
	write      func(p []byte) (n int, err error)
}

func (m *matchWriter) Write(p []byte) (n int, err error) {
	m.output = append(m.output, p...)
	if bytes.HasSuffix(m.output, []byte(m.stopString)) {
		m.success = true
		fmt.Fprintf(m.input, "\x03")
	}

	return len(p), nil
}
