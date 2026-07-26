package rv

import (
	"encoding/binary"
	cp "github.com/temnok/rv/cpu"
	"github.com/temnok/rv/isa"
	"github.com/temnok/rv/ram"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/uart"
	"io"
)

func BootLinux(in io.Reader, out io.Writer, timeout int) *state.CPU {
	var (
		ramBaseAddr  = 0x8000_0000
		opensbiAddr  = ramBaseAddr
		kernelAddr   = 0x8020_0000
		dtbAddr      = kernelAddr - 0x2000
		dynInfoAddr  = kernelAddr - 0x40
		diskBaseAddr = 0x8800_0000

		cpu  = cp.New(512*1024*1024, opensbiAddr)
		uart = uart.New(cpu)
	)

	cpu.Devices = uart.Access

	dir := "build/output/"

	ram.PopulateFromFile(cpu.RAM, opensbiAddr, dir+"/opensbi")
	ram.PopulateFromFile(cpu.RAM, dtbAddr, dir+"/devtree")

	buf := make([]byte, 8*6)
	binary.Encode(buf, binary.LittleEndian, []uint64{0x4942534f, 2, uint64(kernelAddr), 1, 0, 0})
	ram.Populate(cpu.RAM, dynInfoAddr, buf)

	ram.PopulateFromFile(cpu.RAM, kernelAddr, dir+"/kernel")
	ram.PopulateFromFile(cpu.RAM, diskBaseAddr, dir+"/ramdisk")

	cpu.X[isa.A1] = dtbAddr
	cpu.X[isa.A2] = dynInfoAddr

	inChan := make(chan byte)
	go func() {
		b := []byte{0}
		for {
			if n, _ := in.Read(b); n > 0 {
				inChan <- b[0]
			}
		}
	}()

steps:
	for step := 0; timeout == 0 || step < timeout; step++ {
		if uart.AcceptsInput() {
			select {
			case b := <-inChan:
				if b == 'D'-'@' {
					break steps
				}

				uart.SetInput(b)
			default:
			}
		}

		if uart.HasOutput() {
			out.Write([]byte{uart.GetOutput()})
		}

		cp.Step(cpu)
	}

	//debug.Dump(cpu)

	return cpu
}
