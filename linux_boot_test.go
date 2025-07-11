package rv

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestRunKernel32(t *testing.T) {
	testRunKernel(t, 32)
}

func TestRunKernel64(t *testing.T) {
	testRunKernel(t, 64)
}

func testRunKernel(t *testing.T, xlen int) {
	t.Parallel()

	inR, inW := io.Pipe()

	outW := &matchWriter{
		t:     t,
		input: inW,
	}

	runKernel(xlen, inR, outW, 200_000_000)

	if !outW.success {
		t.Fatalf("Expected 'buildroot login: ' prompt, got:\n%v", string(outW.output))
	}
}

type matchWriter struct {
	t       *testing.T
	input   io.Writer
	output  []byte
	success bool
	write   func(p []byte) (n int, err error)
}

func (m *matchWriter) Write(p []byte) (n int, err error) {
	m.output = append(m.output, p...)
	if bytes.HasSuffix(m.output, []byte("buildroot login: ")) {
		m.success = true
		fmt.Fprintf(m.input, "\x03")
	}

	return len(p), nil
}
