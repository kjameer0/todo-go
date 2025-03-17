package main

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func main() {
	// disable input buffering
	term.MakeRaw(int(os.Stdin.Fd()))
	var b []byte = make([]byte, 3)
	for {
		os.Stdin.Read(b)
		fmt.Println(b)
	}
}
