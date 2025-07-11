package main

import (
	rv "github.com/temnok/gorv"
	"os"
)

func main() {
	rv.RunKernel(32, os.Stdin, os.Stdout)
}
