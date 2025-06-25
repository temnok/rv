package main

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	rv "github.com/temnok/gorv"
	"golang.org/x/term"
	"io"
	"os"
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

	kernel := ungzip(check1(os.ReadFile("linux/bin/fw_payload.bin.gz")))
	if sha256sum(kernel) != "0380799b9aa3abe5ea38e0d0ab90f5c613d3bc42952522b62a491c615b7bbcdd" {
		panic("kernel sha check fail")
	}

	ram.Load(ramBaseAddr, kernel)
	ram.Load(dtbAddr, check1(os.ReadFile("linux/bin/rv.dtb")))

	for step := 0; !terminal.Closed; step++ {
		cpu.Step()

		//if step%10_000_000 == 0 {
		//	fmt.Printf("*** Steps:%vM, irq:%v, traps:%v, CLINT:%v, CLINTirq:%v, PLIC:%v, PLICirq:%v, UART:%v\r\n",
		//		step/1_000_000,
		//		cpu.InterruptCount,
		//		cpu.TrapCount,
		//		clint.AccessCount,
		//		clint.InterruptCount,
		//		plic.AccessCount,
		//		plic.InterruptCount,
		//		uart1.AccessCount,
		//	)
		//}
	}
}

func ungzip(data []byte) []byte {
	r := check1(gzip.NewReader(bytes.NewReader(data)))
	return check1(io.ReadAll(r))
}

func sha256sum(data []byte) string {
	sha := sha256.Sum256(data)
	return hex.EncodeToString(sha[:])
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
