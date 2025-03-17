package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/term"
)

// TODO:  create struct
// TODO:  add tasks
// TODO:  delete task
type task struct {
	id   int
	name string
}

func newTask(id int, name string) *task {
	t := &task{id: id, name: name}
	return t
}
func clearLines(nLines int) {
	for i := 0; i < nLines; i++ {
		fmt.Print("\033[F\033[K") // Move cursor up and clear line
	}
}

const up = "\033[A"
const down = "\033[B"
const left = "\033[D"
const enter = 13

func main() {
	fmt.Println("Welcome to Task Checker, what up?")
	options := []string{"Check tasks", "Add a task", "Delete a Task"}
	selected := 0
	term.NewTerminal(os.Stdin, "")
	_, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal("error in making raw")
	}
	buf := make([]byte, 3)
	for {
		for idx, option := range options {
			arrow := " "
			if idx == selected {
				arrow = ">"
			}
			fmt.Printf("%v %v\n\r", arrow, option)
		}
		_, err = os.Stdin.Read(buf)
		if err != nil {
			log.Fatal("Something went wrong reading input")
		}
		userInput := string(buf)
		if buf[0] == 3 {
			os.Exit(0)
		}
		if (userInput) == up {
			selected = (selected - 1)
			if selected == -1 {
				selected = len(options) - 1
			}
			clearLines(len(options))
		} else if userInput == down {
			selected = (selected + 1) % len(options)
			clearLines(len(options))
		} else if buf[0] == enter {
			userSelection := options[selected]
			fmt.Print("\r")
			clearLines(len(options))
			fmt.Print(userSelection, "\n")
		}
		fmt.Print("\r")
	}
}

// so how would i approach this
//i render and wait for input
// only accept up or down
