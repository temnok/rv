package main

import (
	"github.com/temnok/rv"
)

func main() {
	rv.BootLinux("build/output/kernel.gz")
	//ins.PrintFreqStats()
}
