package rv

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"golang.org/x/term"
	"io"
	"os"
	"strings"
)

func BootLinux(xlen int) {
	state := check1(term.MakeRaw(0))
	defer func() {
		check(term.Restore(0, state))
	}()

	bootLinux(xlen, os.Stdin, os.Stdout, 0)
}

func bootLinux(xlen int, in io.Reader, out io.Writer, timeout int) {
	var (
		cpu   CPU
		ram   RAM
		clint CLINT
		plic  PLIC
		uart  UART
	)

	ramBaseAddr := 0x8000_0000
	dtbAddr := ramBaseAddr + 0x0200_0000

	cpu.Init(xlen, Bus{&ram, &clint, &plic, &uart}, ramBaseAddr, 11, dtbAddr)

	ram.Init(&cpu, ramBaseAddr, 128*1024*1024)
	clint.Init(&cpu, 0x0200_0000)
	plic.Init(&cpu, 0x0C00_0000)

	terminal := newTerminal(in, out)
	uart.Init(&plic, 0x0300_0000, 1, terminal.callback)

	path := fmt.Sprintf("buildroot/output/rv%v", cpu.XLen)
	ram.Load(ramBaseAddr, readFile(path+".kernel.gz", ""))
	ram.Load(dtbAddr, readFile(path+".dtb", ""))

	for step := 0; !terminal.Closed; step++ {
		cpu.Step()

		if timeout > 0 && step > timeout {
			break
		}
	}
}

func readFile(path, checksum string) []byte {
	content := check1(os.ReadFile(path))

	if strings.HasSuffix(path, ".gz") {
		r := check1(gzip.NewReader(bytes.NewReader(content)))
		content = check1(io.ReadAll(r))
	}

	if cs := fmt.Sprintf("%x", sha256.Sum256(content)); checksum != "" && checksum != cs {
		panic(path + " checksum check failed, expected " + cs)
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
