package main

import (
	"fmt"
	"os"

	"golang.org/x/term"
	"todo.com/ansi"
)

func main() {
	// Enable the alternate screen buffer
	fmt.Print("\033[?1049h")
	// Restore normal buffer on exit

	// Put terminal in raw mode
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error enabling raw mode:", err)
		return
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	fmt.Print(ansi.Home)
	fmt.Println("Enter text (Press Ctrl+D to stop):")
	fmt.Print("\r")
	var input []byte
	buf := make([]byte, 1)

	for {
		_, err := os.Stdin.Read(buf)
		if err != nil {
			break
		}
		fmt.Print(string(buf))
		// If Ctrl+D is pressed, exit loop
		if buf[0] == 4 {
			fmt.Println("\r")
			fmt.Print(ansi.AltBufferOff)
			break
		} else if buf[0] == 13 {
			fmt.Println("\r")
		}
		input = append(input, buf[0])
	}
	fmt.Print("Input as byte slice:", string(input))
}
