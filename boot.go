package rv

import (
	"bytes"
	"compress/gzip"
	"github.com/temnok/rv/clint"
	cp "github.com/temnok/rv/cpu"
	"github.com/temnok/rv/debug"
	"github.com/temnok/rv/isa"
	"github.com/temnok/rv/plic"
	"github.com/temnok/rv/ram"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/terminal"
	"github.com/temnok/rv/uart"
	"golang.org/x/term"
	"io"
	"os"
	"strings"
)

func BootLinux(dtbPath, kernelPath string) {
	state := check1(term.MakeRaw(0))
	defer func() {
		check(term.Restore(0, state))
	}()

	bootLinux(dtbPath, kernelPath, os.Stdin, os.Stdout, 0)
}

func bootLinux(dtbPath, kernelPath string, in io.Reader, out io.Writer, timeout int) *state.CPU {
	var (
		ramBaseAddr = 0x8000_0000
		kernelAddr  = 0x8020_0000
		cpu         = cp.New(kernelAddr)
		ram         = ram.New(cpu, ramBaseAddr, 512*1024*1024)
		clint       = clint.New(cpu, 0x0200_0000)
		plic        = plic.New(cpu, 0x0C00_0000)
		terminal    = terminal.New(in, out)
		uart        = uart.New(plic, 0x0300_0000, 1, terminal.Callback)
	)

	cpu.Priv = 1
	cpu.CSR.Mideleg = 0x0000000000000222
	cpu.CSR.Medeleg = 0x000000000004b109

	cpu.Bus = state.Bus{ram, clint, plic, uart}

	ram.Load(ramBaseAddr, readFile(dtbPath))
	cpu.X[isa.A1] = ramBaseAddr

	ram.Load(kernelAddr, readFile(kernelPath))

	for step := 0; !terminal.Closed; step++ {
		//if step == 336000 {
		//if trap.IsEntered(cpu) {
		//	debug.Dump(cpu)
		//	break
		//}

		ok := debug.Step(cpu)
		//ok := cp.Step(cpu)

		if !ok || (timeout > 0 && step > timeout) {
			break
		}
	}

	return cpu
}

func readFile(path string) []byte {
	content := check1(os.ReadFile(path))

	if strings.HasSuffix(path, ".gz") {
		r := check1(gzip.NewReader(bytes.NewReader(content)))
		content = check1(io.ReadAll(r))
	}

	return content
}

func check1[A any](a A, err error) A {
	check(err)
	return a
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
