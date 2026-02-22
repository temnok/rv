package main

import (
	"github.com/temnok/rv"
	"github.com/temnok/rv/ins"
)

func main() {
	rv.BootLinux("biko/output")

	ins.PrintFreqStats()
}
