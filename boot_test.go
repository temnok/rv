package rv

import (
	"bytes"
	"fmt"
	"github.com/temnok/rv/debug"
	"io"
	"testing"
)

func TestBoot(t *testing.T) {
	inR, inW := io.Pipe()

	outW := &matchWriter{
		t:          t,
		stopString: "user@rv",
		input:      inW,
	}

	cpu := BootLinux(inR, outW, 110_000_000)

	if !outW.success {
		debug.Dump(cpu)

		t.Fatalf("Expected '%v' stop string, got:\n%v", outW.stopString, string(outW.output))
	}

	fmt.Printf("TLB lookups/misses: %v/%v\n", cpu.TLB.LookupCount, cpu.TLB.MissCount)
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
		fmt.Fprintf(m.input, "\x04")
	}

	return len(p), nil
}
