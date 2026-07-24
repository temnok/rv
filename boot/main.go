package main

import (
	"github.com/temnok/rv"
	"golang.org/x/term"
	"os"
)

func main() {
	term.MakeRaw(0)

	rv.BootLinux(os.Stdin, os.Stdout, 0)
}
