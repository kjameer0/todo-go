package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	ki "todo.com/keypressinterface"

	"github.com/aidarkhanov/nanoid"
	"golang.org/x/term"
)

// TODO: do not allow empty tasks to be added

const up = "\033[A"
const down = "\033[B"
const left = "\033[D"
const enter = 13

const CHECK_TASKS = "Check tasks"
const UPDATE_TASK = "Update tasks"
const ADD_A_TASK = "Add a task"
const DELETE_A_TASK = "Delete a Task"
const QUIT = "Quit"

type task struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}

func (t *task) String() string {
	var completed string
	if !t.Completed {
		completed = "❌"
	} else {
		completed = "✅"
	}
	return fmt.Sprintf("%s %s", t.Name, completed)
}

type app struct {
	Tasks                 map[string]*task `json:"tasks"`
	InsertionOrder        []string         `json:"insertionOrder"`
	originalTerminalState *term.State
	t                     *term.Terminal
	saveLocation          string
}
type saveData struct {
	Tasks          map[string]*task `json:"tasks"`
	InsertionOrder []string         `json:"insertionOrder"`
}

func newApp() *app {
	tasks := make(map[string]*task, 100)
	return &app{Tasks: tasks}
}

func newTask(name string, completed bool) *task {
	var taskId string
	taskId, err := nanoid.Generate(nanoid.DefaultAlphabet, 20)
	if err != nil {
		log.Fatal(err)
	}
	t := &task{Id: taskId, Name: name, Completed: completed}
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

func saveToFile(a *app) {
	s := saveData{}
	s.Tasks = a.Tasks
	s.InsertionOrder = a.InsertionOrder
	taskJson, err := json.Marshal(s)
	if err != nil {
		log.Fatal("failed to convert tasks to JSON")
	}
	err = os.WriteFile(a.saveLocation, taskJson, 0644)
	if err != nil {
		log.Fatal("failed to write to file ", a.saveLocation)
	}
}

func readTasksFromFile(a *app) {
	data, err := os.ReadFile(a.saveLocation)
	s := saveData{}
	if err != nil {
		log.Fatal("failed to read from save location", err)
	}
	json.Unmarshal(data, &s)
	a.InsertionOrder = s.InsertionOrder
	if len(s.Tasks) > 0 {
		a.Tasks = s.Tasks
	}
}

// task functions
func addTask(a *app, taskText string) {
	addedTask := newTask(taskText, false)
	a.Tasks[addedTask.Id] = addedTask
	a.InsertionOrder = append(a.InsertionOrder, addedTask.Id)
	saveToFile(a)
}
func removeTask(a *app, taskId string) bool {
	if _, ok := a.Tasks[taskId]; !ok {
		return false
	}
	delete(a.Tasks, taskId)
	for idx, id := range a.InsertionOrder {
		if id == taskId {
			a.InsertionOrder[idx] = ""
			break
		}
	}
	saveToFile(a)
	return true
}
func listTasks(a *app) {
	fmt.Println("Tasks:")

	for _, taskId := range a.InsertionOrder {
		if taskId == "" {
			continue
		}
		curTask := a.Tasks[taskId]
		var completed string
		if !curTask.Completed {
			completed = "❌"
		} else {
			completed = "✅"
		}
		fmt.Printf("\t%s %s id: %s\n", curTask.Name, completed, curTask.Id)
	}
}
func updateTask(a *app, t *task) {
	t.Completed = !t.Completed
	saveToFile(a)
}
func handleOption(a *app, options []string, selected int) {
	// switch on different options
	chosenOption := options[selected]
	term.Restore(int(os.Stdin.Fd()), a.originalTerminalState)
	defer term.MakeRaw(int(os.Stdin.Fd()))

	switch chosenOption {
	case CHECK_TASKS:
		listTasks(a)
	case UPDATE_TASK:
		// TODO: figure out update logic
		items := []*task{}
		fmt.Println(a.InsertionOrder)
		for _, item := range a.InsertionOrder {
			taskItem := a.Tasks[item]
			items = append(items, taskItem)
		}
		m, err := ki.NewMatrixMenu(items, int(os.Stdin.Fd()))
		if err != nil {
			log.Fatal("failed to generate menu")
		}
		selection, err := m.RenderInterface()
		if err != nil {
			log.Fatal("menu failed in update task")
		}
		updateTask(a, selection)
		if err != nil {
			log.Fatal("failed to render menu for updating tasks")
		}
		fmt.Println(string(a.t.Escape.Green) + "Task Updated" + string(a.t.Escape.Reset))
	case ADD_A_TASK:
		r := bufio.NewReader(os.Stdin)
		fmt.Print("Enter task: ")
		userTask, err := r.ReadString('\n')
		if err != nil {
			log.Fatal("Error reading task to add")
		}
		userTask = strings.TrimSpace(userTask)
		addTask(a, userTask)
		fmt.Println(string(a.t.Escape.Green) + "Task Added" + string(a.t.Escape.Reset))
	case DELETE_A_TASK:
		items := []*task{}
		for _, item := range a.Tasks {
			items = append(items, item)
		}
		fmt.Print("Choose task to delete: \n")
		menu, err := ki.NewMatrixMenu(items, int(os.Stdin.Fd()))
		if err != nil {
			log.Fatal("failed to create delete menu")
		}
		selection, err := menu.RenderInterface()
		if err != nil {
			log.Fatal("failed to delete")
		}
		wasRemoved := removeTask(a, selection.Id)
		if wasRemoved {
			fmt.Println(string(a.t.Escape.Green) + "Task Deleted" + string(a.t.Escape.Reset))
		} else {
			fmt.Println(string(a.t.Escape.Yellow) + "No task matching the provided id" + string(a.t.Escape.Reset))
		}
	case QUIT:
		exitCleanup(a)
	}
}

func main() {
	a := newApp()
	a.saveLocation = "./tasks.json"
	readTasksFromFile(a)
	fmt.Println("Welcome to Task Checker, what up?")

	options := []string{CHECK_TASKS, UPDATE_TASK, ADD_A_TASK, DELETE_A_TASK, QUIT}
	selected := 0

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	t := term.NewTerminal(os.Stdin, "")
	if err != nil {
		log.Fatal("error in making raw")
	}
	a.originalTerminalState = oldState
	a.t = t

	// read three bytes at once from stdin to capture arrow key presses
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
