package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"todo.com/ansi"
	navmenu "todo.com/nav-menu"

	"github.com/aidarkhanov/nanoid"
	"golang.org/x/term"
)

// TODO: allow user to adjust config
// TODO: create flow to allow user to submit a start date for a task, default today
// TODO: allow user to set due date for task
// TODO: allow user to filter by today and not today??
type taskDate struct {
	t time.Time
}

func (t taskDate) String() string {
	return monthDayYear(t.t)
}

type stringWrapper string

func (s stringWrapper) String() string {
	return string(s)
}

const CHECK_TASKS stringWrapper = "Check tasks"
const UPDATE_TASK stringWrapper = "Update tasks"
const ADD_A_TASK stringWrapper = "Add a task"
const DELETE_A_TASK stringWrapper = "Delete a specific task"
const DELETE_ALL_TASKS stringWrapper = "Delete -all- tasks"
const QUIT stringWrapper = "Quit"

var options = []stringWrapper{CHECK_TASKS, UPDATE_TASK, ADD_A_TASK, DELETE_A_TASK, DELETE_ALL_TASKS, QUIT}

type task struct {
	Id             string    `json:"id"`
	Name           string    `json:"name"`
	Completed      bool      `json:"completed"`
	CompletionDate time.Time `json:"completionDate"`
	BeginDate      time.Time `json:"beginDate"`
}

func (t *task) String() string {
	var completed string
	if !t.Completed {
		completed = "❌"
	} else {
		completed = "✅"
	}
	var completionDate string
	if t.CompletionDate.IsZero() {
		completionDate = ""
	} else {
		completionDate = monthDayYear(t.CompletionDate)
	}
	return fmt.Sprintf("%s %s %s", t.Name, completed, completionDate)
}

type app struct {
	Tasks                 map[string]*task `json:"tasks"`
	InsertionOrder        []string         `json:"insertionOrder"`
	originalTerminalState *term.State
	t                     *term.Terminal
	saveLocation          string
	configPath            string
	config                *config
}
type saveData struct {
	Tasks          map[string]*task `json:"tasks"`
	InsertionOrder []string         `json:"insertionOrder"`
}

func newApp() *app {
	tasks := make(map[string]*task, 100)
	return &app{Tasks: tasks}
}

func newTask(name string, completed bool, beginDate time.Time) *task {
	if name == "" {
		log.Fatal("a task must have a name")
	}
	var taskId string
	taskId, err := nanoid.Generate(nanoid.DefaultAlphabet, 20)
	if err != nil {
		log.Fatal(err)
	}
	t := &task{Id: taskId, Name: name, Completed: completed, BeginDate: beginDate}
	return t
}

func clearLines(nLines int) {
	for i := 0; i < nLines; i++ {
		fmt.Print("\033[F\033[K") // Move cursor up and clear line
	}
}
func exitCleanup(a *app) {
	os.Exit(0)
}

func (a *app) listInsertionOrder() []*task {
	items := []*task{}
	for _, item := range a.InsertionOrder {
		taskItem, ok := a.Tasks[item]
		if ok {
			items = append(items, taskItem)
		}
	}
	return items
}

// task functions
func addTask(a *app, taskText string, beginTime time.Time) {
	addedTask := newTask(taskText, false, beginTime)
	a.Tasks[addedTask.Id] = addedTask
	a.InsertionOrder = append(a.InsertionOrder, addedTask.Id)
	saveToFile(a)
}
func removeTask(a *app, taskId string) bool {
	if _, ok := a.Tasks[taskId]; !ok {
		return false
	}
	delete(a.Tasks, taskId)
	//remove deleted id from insertion order
	filteredInsertionOrder := []string{}
	for _, id := range a.InsertionOrder {
		if id == taskId {
			continue
		}
		filteredInsertionOrder = append(filteredInsertionOrder, id)
	}
	a.InsertionOrder = filteredInsertionOrder
	saveToFile(a)
	return true
}
func removeAllTasks(a *app) {
	clear(a.InsertionOrder)
	clear(a.Tasks)
	saveToFile(a)
}
func listTasks(a *app) {
	for _, taskId := range a.InsertionOrder {
		if taskId == "" {
			continue
		}
		curTask := a.Tasks[taskId]
		//show a task if it not complete or if show complete and task
		if !a.config.ShowComplete && curTask.Completed {
			continue
		}
		if time.Now().Compare(curTask.BeginDate) == -1 {
			continue
		}
		var completed string
		if !curTask.Completed {
			completed = "❌"
		} else {
			completed = "✅"
		}
		t := monthDayYear(curTask.CompletionDate)
		if curTask.CompletionDate.IsZero() {
			t = ""
		}
		fmt.Printf("\t%s %s %s\n", curTask.Name, completed, t)
	}
}
func updateTask(a *app, t *task) {
	t.Completed = !t.Completed
	if t.Completed {
		t.CompletionDate = time.Now()
	}
	saveToFile(a)
}
func handleOption(a *app, option stringWrapper) error {
	switch option {
	case CHECK_TASKS:
		if len(a.Tasks) == 0 {
			fmt.Println("No tasks currently")
		} else {
			fmt.Println("Tasks:")
			listTasks(a)
		}
	case UPDATE_TASK:
		if len(a.Tasks) == 0 {
			fmt.Println("No tasks to update")
			break
		}
		taskItems := a.listInsertionOrder()
		m := navmenu.NewMenu(taskItems, int(os.Stdin.Fd()))
		selection, err := m.Render()
		if err != nil {
			log.Fatal("menu failed in update task")
		}
		updateTask(a, selection)
		if err != nil {
			log.Fatal("failed to render menu for updating tasks")
		}
		fmt.Println(ansi.Green + "Task Updated" + ansi.Reset)
	case ADD_A_TASK:
		r := bufio.NewReader(os.Stdin)
		fmt.Print("Enter task: ")
		userTask, err := r.ReadString('\n')
		fmt.Print("\r", ansi.ClearLine+"\n")
		if err != nil {
			return err
		}
		userTask = strings.TrimSpace(userTask)
		fmt.Print("When do you want this task to appear on your list:\n")
		dates := []taskDate{}
		curTime := time.Now()
		for i := 0; i < 14; i++ {
			dates = append(dates, taskDate{t: addDayToDate(curTime, i)})
		}
		m := navmenu.NewMenu(dates, int(os.Stdin.Fd()))
		beginDate, err := m.Render()
		if err != nil {
			return err
		}
		addTask(a, userTask, beginDate.t)
		fmt.Println(ansi.Green + "Task Added" + ansi.Reset)
	case DELETE_A_TASK:
		if len(a.Tasks) == 0 {
			fmt.Println("No tasks to delete")
			break
		}
		items := a.listInsertionOrder()
		fmt.Print("Choose task to delete: \n")
		m := navmenu.NewMenu(items, int(os.Stdin.Fd()))
		selection, err := m.Render()
		if err != nil {
			log.Fatal("failed to delete")
		}
		wasRemoved := removeTask(a, selection.Id)
		if wasRemoved {
			fmt.Println(string(ansi.Green) + "Task Deleted" + string(ansi.Reset))
		} else {
			fmt.Println(string(ansi.Yellow) + "No task matching the provided id" + string(ansi.Reset))
		}
	case DELETE_ALL_TASKS:
		if len(a.Tasks) == 0 {
			fmt.Println("No tasks to delete")
			break
		}
		y := "yes"
		n := "no"
		yOrN := []stringWrapper{stringWrapper(y), stringWrapper(n)}

		m := navmenu.NewMenu(yOrN, int(os.Stdin.Fd()))

		fmt.Println("Are you sure you want to delete all of your tasks?:")
		selection, err := m.Render()
		if err != nil {
			log.Fatal("failed to select yes or no")
		}

		if string(selection) == y {
			removeAllTasks(a)
			fmt.Println(ansi.Green + "All Tasks Deleted" + ansi.Reset)
		} else if string(selection) == n {
			fmt.Println(ansi.Green + "Deletion cancelled" + ansi.Reset)
		}
	case QUIT:
		exitCleanup(a)
	}
	return nil
}

func main() {
	a := newApp()
	a.saveLocation = "./tasks.json"
	a.configPath = "./config.json"
	c, err := a.loadConfig()
	if err != nil {
		log.Fatal(err)
	}
	a.config = c
	readTasksFromFile(a)
	fmt.Println("Welcome to Task Checker, what up?")

	for {
		m := navmenu.NewMenu(options, int(os.Stdin.Fd()))
		option, err := m.Render()
		if err != nil {
			fmt.Println(err)
			continue
		}
		err = handleOption(a, option)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Cancelled adding a task")
				continue
			}
		}
	}
}
