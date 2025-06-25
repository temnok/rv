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
		ramBaseAddr = int32(-1 << 31)
		dtbAddr     = ramBaseAddr + 0x200_0000
	)

	cpu.Init(rv.Bus{&ram, &clint, &plic, &uart1, &uart2}, ramBaseAddr, []int32{
		11: dtbAddr,
	})

	ram.Init(ramBaseAddr, 128*1024*1024)
	clint.Init(&cpu, 0x200_0000)
	plic.Init(&cpu, 0xC00_0000)

	terminal := newTerminal()
	uart1.Init(&plic, 0x300_0000, 1, terminal.callback)
	uart2.Init(&plic, 0x600_0000, 2, nil)

	ram.Load(ramBaseAddr, readFile("linux/bin/fw_payload.bin.gz",
		"0380799b9aa3abe5ea38e0d0ab90f5c613d3bc42952522b62a491c615b7bbcdd"))
	ram.Load(dtbAddr, readFile("linux/bin/rv.dtb",
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

	if checksum != "" && checksum != fmt.Sprintf("%x", sha256.Sum256(content)) {
		panic("file checksum check failed")
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
