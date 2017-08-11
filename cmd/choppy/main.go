package main

import (
	"os"

	"github.com/fogleman/choppy/chopsui"
)

func main() {
	args := os.Args[1:]
	if len(args) > 0 {
		chopsui.Run(args[0])
	} else {
		chopsui.Run("")
	}
}
