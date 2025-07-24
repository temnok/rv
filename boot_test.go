package rv

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestRunKernel32(t *testing.T) {
	testRunKernel(t, 32, "biko/output-32", "user@rv")
}

func TestRunKernel64(t *testing.T) {
	testRunKernel(t, 64, "biko/output-64", "user@rv")
}

func testRunKernel(t *testing.T, xlen int, kernelPath, stopString string) {
	t.Parallel()

	inR, inW := io.Pipe()

	outW := &matchWriter{
		t:          t,
		stopString: stopString,
		input:      inW,
	}

	bootLinux(xlen, kernelPath, inR, outW, 200_000_000)

	if !outW.success {
		t.Fatalf("Expected '%v' stop string, got:\n%v", stopString, string(outW.output))
	}
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
