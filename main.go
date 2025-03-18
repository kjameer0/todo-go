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
// TODO:  figure out how to select between different options and do something
const up = "\033[A"
const down = "\033[B"
const left = "\033[D"
const enter = 13

const CHECK_TASKS = "Check tasks"
const ADD_A_TASK = "Add a task"
const DELETE_A_TASK = "Delete a Task"
const QUIT = "Quit"

type app struct {
	tasks []task
}

func newApp(tasks []task) *app {
	return &app{tasks: tasks}
}

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

// task functions
func addTask(tasks *[]task, task task) {
	*tasks = append(*tasks, task)
}
func removeTask(tasks *[]task, taskId int) {
	filteredTasks := []task{}
	for _, task := range *tasks {
		if task.id != taskId {
			filteredTasks = append(filteredTasks, task)
		}
	}
	*tasks = filteredTasks
}
func handleOption(options []string, selected int) {
	// switch on different options
	chosenOption := options[selected]
	switch chosenOption {
	case CHECK_TASKS:
		fmt.Println("checking tasks")
	case ADD_A_TASK:
		fmt.Println("Adding task")
	case DELETE_A_TASK:
		fmt.Println("deleting task")
	case QUIT:
		os.Exit(0)
	}

}

func main() {
	a := newApp([]task{})
	print(a)
	fmt.Println("Welcome to Task Checker, what up?")

	options := []string{"Check tasks", "Add a task", "Delete a Task", "Quit"}
	selected := 0

	term.NewTerminal(os.Stdin, "")
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	defer term.Restore(int(os.Stdin.Fd()), oldState)
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
			term.Restore(int(os.Stdin.Fd()), oldState)
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
