package rv

import (
	"bytes"
	"compress/gzip"
	"github.com/temnok/rv/clint"
	cp "github.com/temnok/rv/cpu"
	"github.com/temnok/rv/plic"
	"github.com/temnok/rv/ram"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/terminal"
	"github.com/temnok/rv/uart"
	"io"
	"os"
	"strings"
)

func BootLinux(dir string) (err error) {
	terminal.WithRaw(func(in io.Reader, out io.Writer) {
		_, err = bootLinux(dir, in, out, 0)
	})

	return
}

func bootLinux(kernelPath string, in io.Reader, out io.Writer, timeout int) (*state.CPU, error) {
	var (
		ramBaseAddr = 0x8000_0000
		cpu         = cp.New(ramBaseAddr)
		ram         = ram.New(cpu, ramBaseAddr, 512*1024*1024)
		clint       = clint.New(cpu, 0x0200_0000)
		plic        = plic.New(cpu, 0x0C00_0000)
		terminal    = terminal.New(in, out)
		uart        = uart.New(plic, 0x0300_0000, 1, terminal.GetChar, terminal.PutChar)
	)

	cpu.Bus = state.Bus{ram, clint, plic, uart}

	kernelBytes, err := readFile(kernelPath)
	if err != nil {
		return nil, err
	}

	ram.Load(ramBaseAddr, kernelBytes)

	for step := 0; !terminal.Closed; step++ {
		ok := cp.Step(cpu)

		if !ok || (timeout > 0 && step > timeout) {
			break
		}

		uart.IO()
	}

	return cpu, nil
}

func readFile(path string) ([]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if !strings.HasSuffix(path, ".gz") {
		return content, nil
	}

	r, err := gzip.NewReader(bytes.NewReader(content))
	if err != nil {
		return nil, err
	}

	return io.ReadAll(r)
}
