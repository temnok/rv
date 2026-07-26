package rv

import (
	"encoding/binary"
	cp "github.com/temnok/rv/cpu"
	"github.com/temnok/rv/isa"
	"github.com/temnok/rv/ram"
	"github.com/temnok/rv/state"
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

		cpu = cp.New(512*1024*1024, opensbiAddr)
	)

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

	ctrlD := false

	cpu.UARTInput = func() (byte, bool) {
		select {
		case b := <-inChan:
			if b == 'D'-'@' {
				ctrlD = true
			}

			return b, true
		default:
			return 0, false
		}
	}

	cpu.UARTOutput = func(b byte) bool {
		out.Write([]byte{b})
		return true
	}

	for step := 0; !ctrlD && (timeout == 0 || step < timeout); step++ {
		cp.Step(cpu)
	}

	//debug.Dump(cpu)

	return cpu
}
