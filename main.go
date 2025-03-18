package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aidarkhanov/nanoid"
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

type task struct {
	id   string
	name string
}
type app struct {
	tasks                 map[string]*task
	insertionOrder        []string
	originalTerminalState *term.State
	t                     *term.Terminal
}

func createTaskMap(a *app, tasks []string) map[string]*task {
	taskMap := make(map[string]*task)
	for _, taskText := range tasks {
		nano, _ := nanoid.Generate(nanoid.DefaultAlphabet, 5)
		taskMap[nano] = newTask(nano, taskText)
		a.insertionOrder = append(a.insertionOrder, nano)
	}
	return taskMap
}
func newApp() *app {
	return &app{}
}

func newTask(id string, name string) *task {
	t := &task{id: id, name: name}
	return t
}
func clearLines(nLines int) {
	for i := 0; i < nLines; i++ {
		fmt.Print("\033[F\033[K") // Move cursor up and clear line
	}
}
func exitCleanup(a *app) {
	term.Restore(int(os.Stdin.Fd()), a.originalTerminalState)
	os.Exit(0)
}

// task functions
func addTask(a *app, taskText string) {
	var taskId string
	i:= 0
	for _, ok := a.tasks[taskId]; ok; {
		newId, err := nanoid.Generate(nanoid.DefaultAlphabet, 5)
		if err != nil {
			log.Fatal("problem generating nanoid when adding task")
		}
		taskId = newId
		i+=1
		if i > 500 {
			log.Fatal("Too many attempts to generate id")
		}
	}
	a.tasks[taskId] = &task{id: taskId, name: taskText}
	a.insertionOrder = append(a.insertionOrder, taskId)
}
func removeTask(a *app, taskId string) {
	delete(a.tasks, taskId)
	for idx, id := range a.insertionOrder {
		if id == taskId {
			a.insertionOrder[idx] = ""
			return
		}
	}
}
func listTasks(a *app) {
	fmt.Println("Tasks:")
	for _, taskId := range a.insertionOrder {
		if taskId == "" {
			continue
		}
		curTask := a.tasks[taskId]
		fmt.Printf("\t%s id: %s\n", curTask.name, curTask.id)
	}
}
func handleOption(a *app, options []string, selected int) {
	// switch on different options
	chosenOption := options[selected]
	term.Restore(int(os.Stdin.Fd()), a.originalTerminalState)
	defer term.MakeRaw(int(os.Stdin.Fd()))

	switch chosenOption {
	case CHECK_TASKS:
		listTasks(a)
	case ADD_A_TASK:
		fmt.Println("Adding task")
	case DELETE_A_TASK:
		// TODO: come up with interface for deleting
		fmt.Println("deleting task")
	case QUIT:
		exitCleanup(a)
	}
}

func main() {
	a := newApp()
	a.tasks = createTaskMap(a, []string{"check phone", "look at phone", "put phone down"})
	fmt.Println("Welcome to Task Checker, what up?")

	options := []string{"Check tasks", "Add a task", "Delete a Task", "Quit"}
	selected := 0

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	t := term.NewTerminal(os.Stdin, "")
	if err != nil {
		log.Fatal("error in making raw")
	}
	a.originalTerminalState = oldState
	a.t = t

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
			exitCleanup(a)
		}
		clearLines(len(options))

		if (userInput) == up {
			selected = (selected - 1)
			if selected == -1 {
				selected = len(options) - 1
			}
		} else if userInput == down {
			selected = (selected + 1) % len(options)
		} else if buf[0] == enter {
			handleOption(a, options, selected)
		}
		fmt.Print("\r")
	}
}
