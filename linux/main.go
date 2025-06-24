package main

import (
	rv "github.com/temnok/gorv"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
)

func main() {
	state, err := terminal.MakeRaw(0)
	if err != nil {
		log.Fatalln("Could not make stdin a raw terminal: ", err.Error())
	}
	defer terminal.Restore(0, state)

	var (
		cpu   rv.CPU
		ram   rv.RAM
		clint rv.CLINT
		plic  rv.PLIC
		uart  rv.UART
	)

	const ramBaseAddr = -1 << 31
	cpu.Init(ramBaseAddr, rv.Bus{&ram, &clint, &plic, &uart})
	ram.Init(ramBaseAddr, 128*1024*1024, nil)
	clint.Init(&cpu, 0x200_0000)
	plic.Init(&cpu, 0xC00_0000)
	uart.Init(&plic, 0x300_0000, terminalCallback)
	
}

func terminalCallback(ch *byte, write bool) bool {
	if write {
		buf := []byte{*ch}
		os.Stdout.Write(buf)
	}

	return true
}
