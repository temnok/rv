package main

import (
	"github.com/temnok/rv"
)

func main() {
	rv.BootLinux("biko/output/rv.dtb", "biko/output/kernel.gz")
	//rv.BootLinux("biko/output/biko.gz")
	//ins.PrintFreqStats()
}
