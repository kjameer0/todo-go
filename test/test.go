package main

import "fmt"

func main() {
	// disable input buffering
	var s string
	fmt.Scanln(&s)
	fmt.Println(s, "string")
}
