package rv

import (
	"encoding/binary"
	cp "github.com/temnok/rv/cpu"
	"github.com/temnok/rv/isa"
	"github.com/temnok/rv/plic"
	"github.com/temnok/rv/ram"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/uart"
	"golang.org/x/term"
	"io"
	"os"
)

func BootLinux() {
	termState, err := term.MakeRaw(0)
	if err != nil {
		panic(err)
	}

	defer func() {
		term.Restore(0, termState)
	}()

	bootLinux(os.Stdin, os.Stdout, 0)
}

func bootLinux(in io.Reader, out io.Writer, timeout int) *state.CPU {
	var (
		ramBaseAddr  = 0x8000_0000
		opensbiAddr  = ramBaseAddr
		kernelAddr   = 0x8020_0000
		dtbAddr      = kernelAddr - 0x2000
		dynInfoAddr  = kernelAddr - 0x40
		diskBaseAddr = 0x8800_0000

		cpu  = cp.New(opensbiAddr)
		ram  = ram.New(512 * 1024 * 1024)
		plic = plic.New(cpu)
		uart = uart.New(plic, 1, in, out)
	)

	cpu.RAM = ram.Access
	cpu.Devices = []state.Device{2: plic.Access, 3: uart.Access}

	dir := "build/output/"

	ram.LoadFile(opensbiAddr, dir+"/opensbi")
	ram.LoadFile(dtbAddr, dir+"/devtree")

	buf := make([]byte, 8*6)
	binary.Encode(buf, binary.LittleEndian, []uint64{0x4942534f, 2, uint64(kernelAddr), 1, 0, 0})
	ram.Load(dynInfoAddr, buf)

	ram.LoadFile(kernelAddr, dir+"/kernel")
	ram.LoadFile(diskBaseAddr, dir+"/ramdisk")

	cpu.X[isa.A1] = dtbAddr
	cpu.X[isa.A2] = dynInfoAddr

	for step := 0; uart.Input() != 'D'-'@'; step++ {
		cp.Step(cpu)

		//if !debug.Step(cpu) {
		//	break
		//}

		if timeout != 0 && step > timeout {
			break
		}
	}

	//debug.Dump(cpu)

	return cpu
}
