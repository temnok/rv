package main

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	rv "github.com/temnok/gorv"
	"golang.org/x/term"
	"io"
	"os"
	"strings"
)

func main() {
	state := check1(term.MakeRaw(0))
	defer term.Restore(0, state)

	var (
		cpu          rv.CPU
		ram          rv.RAM
		clint        rv.CLINT
		plic         rv.PLIC
		uart1, uart2 rv.UART
	)

	const (
		ramBaseAddr = rv.Xint(-1 << 31)
		dtbAddr     = ramBaseAddr + 0x200_0000
	)

	cpu.Init(rv.Bus{&ram, &clint, &plic, &uart1, &uart2}, ramBaseAddr, []rv.Xint{
		11: dtbAddr,
	})

	ram.Init(ramBaseAddr, 128*1024*1024)
	clint.Init(&cpu, 0x0200_0000)
	plic.Init(&cpu, 0x0C00_0000)

	terminal := newTerminal()
	uart1.Init(&plic, 0x0300_0000, 1, terminal.callback)
	uart2.Init(&plic, 0x0600_0000, 2, nil)

	ram.Load(ramBaseAddr, readFile("linux/buildroot/rv32imac/bin/fw_payload.bin.gz",
		""))
	ram.Load(dtbAddr, readFile("linux/buildroot/rv32imac/bin/rv.dtb",
		"46f0e3d9ddba89ec2e8251409bca4daca38715a847954a526feb2be9065347e7"))

	for step := 0; !terminal.Closed; step++ {
		cpu.Step()
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
