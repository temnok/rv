package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	rv "github.com/temnok/gorv"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"os"
)

func main() {
	state := check1(terminal.MakeRaw(0))
	defer terminal.Restore(0, state)

	var (
		cpu   rv.CPU
		ram   rv.RAM
		clint rv.CLINT
		plic  rv.PLIC
		uart  rv.UART
	)

	const (
		ramBaseAddr = int32(-1 << 31)
		dtbAddr     = ramBaseAddr + 0x200_0000
	)

	cpu.Init(rv.Bus{&ram, &clint, &plic, &uart}, ramBaseAddr, []int32{
		11: dtbAddr,
	})

	ram.Init(ramBaseAddr, 128*1024*1024)
	clint.Init(&cpu, 0x200_0000)
	plic.Init(&cpu, 0xC00_0000)
	uart.Init(&plic, 0x300_0000, terminalCallback)

	ram.Load(ramBaseAddr, ungzip(check1(os.ReadFile("linux/bin/fw_payload.bin.gz"))))
	ram.Load(dtbAddr, check1(os.ReadFile("linux/bin/rv.dtb")))

	for step := 0; ; step++ {
		cpu.Step()

		if step%10_000_000 == 0 {
			fmt.Printf("*** Steps:%vM, irq:%v, traps:%v, CLINT:%v, PLIC:%v, UART:%v\r\n",
				step/1_000_000,
				cpu.InterruptCount,
				cpu.TrapCount,
				clint.AccessCount,
				plic.AccessCount,
				uart.AccessCount,
			)
		}
	}
}

func terminalCallback(ch *byte, write bool) bool {
	if write {
		buf := []byte{*ch}
		os.Stdout.Write(buf)
		return true
	}

	return false
}

func ungzip(data []byte) []byte {
	r := check1(gzip.NewReader(bytes.NewReader(data)))
	return check1(io.ReadAll(r))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func check1[A any](a A, err error) A {
	check(err)
	return a
}
