package rv

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestBoot(t *testing.T) {
	testBoot(t, "biko.gz", 200_000_000)
}

func TestBootGo(t *testing.T) {
	testBoot(t, "biko-go.gz", 5_000_000_000)
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

	bootLinux(kernelPath, inR, outW, timeout)

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
