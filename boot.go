package rv

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"github.com/temnok/rv/clint"
	cp "github.com/temnok/rv/cpu"
	"github.com/temnok/rv/isa"
	"github.com/temnok/rv/plic"
	"github.com/temnok/rv/ram"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/uart"
	"golang.org/x/term"
	"io"
	"os"
	"strings"
)

func BootLinux(kernelPath string) {
	termState, err := term.MakeRaw(0)
	if err != nil {
		panic(err)
	}

	defer func() {
		term.Restore(0, termState)
	}()

	bootLinux(kernelPath, os.Stdin, os.Stdout, 0)
}

func bootLinux(kernelPath string, in io.Reader, out io.Writer, timeout int) *state.CPU {
	var (
		ramBaseAddr  = 0x8000_0000
		opensbiAddr  = ramBaseAddr
		kernelAddr   = 0x8020_0000
		dtbAddr      = kernelAddr - 0x2000
		dynInfoAddr  = kernelAddr - 0x40
		diskBaseAddr = 0x8400_0000

		cpu   = cp.New(opensbiAddr)
		ram   = ram.New(ramBaseAddr, 128*1024*1024)
		clint = clint.New(cpu, 0x100_0000)
		plic  = plic.New(cpu, 0x200_0000)
		uart  = uart.New(plic, 0x300_0000, 1, in, out)
	)

	cpu.Bus = state.Bus{ram, clint, plic, uart}

	dir := kernelPath[:strings.LastIndexByte(kernelPath, '/')+1]

	ram.Load(opensbiAddr, readFile(dir+"opensbi.gz"))
	ram.Load(dtbAddr, readFile(dir+"rv.dtb"))

	buf := make([]byte, 8*6)
	binary.Encode(buf, binary.LittleEndian, []uint64{0x4942534f, 2, uint64(kernelAddr), 1, 0, 0})
	ram.Load(dynInfoAddr, buf)

	ram.Load(kernelAddr, readFile(kernelPath))
	ram.Load(diskBaseAddr, readFile(dir+"ramdisk.img"))

	cpu.X[isa.A1] = dtbAddr
	cpu.X[isa.A2] = dynInfoAddr

	for step := 0; uart.Input() != 'D'-'@'; step++ {
		ok := cp.Step(cpu)

		if !ok || (timeout > 0 && step > timeout) {
			break
		}
	}

	return cpu
}

func readFile(path string) []byte {
	content := check1(os.ReadFile(path))
	if !strings.HasSuffix(path, ".gz") {
		return content
	}

	r := check1(gzip.NewReader(bytes.NewReader(content)))
	return check1(io.ReadAll(r))
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
