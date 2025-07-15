package rv

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"errors"
	"fmt"
	"golang.org/x/term"
	"io"
	"os"
	"strings"
)

func BootLinux(xlen int, dir string) {
	state := check1(term.MakeRaw(0))
	defer func() {
		check(term.Restore(0, state))
	}()

	bootLinux(xlen, dir, os.Stdin, os.Stdout, 0)
}

func bootLinux(xlen int, dir string, in io.Reader, out io.Writer, timeout int) {
	var (
		cpu   CPU
		ram   RAM
		clint CLINT
		plic  PLIC
		uart  UART
	)

	ramBaseAddr := 0x8000_0000
	dtbReg, dtbAddr := 11, ramBaseAddr+0x0200_0000

	path := fmt.Sprintf("%v/rv%v", dir, xlen)
	kernelPath := path + ".kernel"
	if !existsFile(kernelPath) {
		kernelPath += ".gz"
	}

	if dtbPath := path + ".dtb"; !existsFile(dtbPath) {
		dtbReg = 0
		dtbAddr = 0
	}

	cpu.Init(xlen, Bus{&ram, &clint, &plic, &uart}, ramBaseAddr, dtbReg, dtbAddr)
	ram.Init(&cpu, ramBaseAddr, 128*1024*1024)
	clint.Init(&cpu, 0x0200_0000)
	plic.Init(&cpu, 0x0C00_0000)

	terminal := newTerminal(in, out)
	uart.Init(&plic, 0x0300_0000, 1, terminal.callback)

	ram.Load(ramBaseAddr, readFile(kernelPath, ""))

	if dtbAddr != 0 {
		ram.Load(dtbAddr, readFile(path+".dtb", ""))
	}

	if initrdPath := path + ".initramfs.cpio.gz"; existsFile(initrdPath) {
		ram.Load(0x87000000, readFile(initrdPath, ""))
	}

	for step := 0; !terminal.Closed; step++ {
		ok := cpu.Step()

		if !ok || (timeout > 0 && step > timeout) {
			break
		}
	}
}

func existsFile(path string) bool {
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false // File does not exist
	}
	return err == nil
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
