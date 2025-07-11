package main

import (
	rv "github.com/temnok/gorv"
	"os"
)

func main() {
	rv.RunKernel(64, os.Stdin, os.Stdout)
}
