package main

import (
	"fmt"
	"time"
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
	fmt.Printf("\033[%dA", nLines)
}
func main() {
	fmt.Println("Welcome to Task Checker, what up?")
	options := []string{"Check tasks", "Add a task", "Delete a Task"}
	selected := 0
	for{
		for idx, t := range options {
			arrow := "  "
			if idx == selected {
				arrow = "> "
			}
			fmt.Printf("%s%s\n",arrow,t)
		}
		selected = selected + 1
		if selected == 3 {
			break
		}
		fmt.Print("\n")
		time.Sleep(time.Second * 1)
		clearLines(len(options)+1)
	}

}
