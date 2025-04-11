package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"todo.com/ansi"
	navmenu "todo.com/nav-menu"

	"github.com/aidarkhanov/nanoid"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"golang.org/x/term"
)

// TODO: allow user to adjust config
// TODO: create flow to allow user to submit a start date for a task, default today
// TODO: allow user to set due date for task
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
func makeList(a *app, items []string, l *tview.List) {
	head := l
	for i, item := range items {
		head = head.AddItem(a.Tasks[item].String(), "", rune(i+'a'), nil)
	}
	head.AddItem("Exit", "", rune('q'), nil)
}

func handleOption(a *app, option stringWrapper) error {
	switch option {
	case CHECK_TASKS:
		app := tview.NewApplication()
		table := tview.NewTable().
			SetBorders(true)
		word := 0
		cols, rows := math.Ceil(math.Sqrt(float64(len(a.Tasks)))), math.Ceil(math.Sqrt(float64(len(a.Tasks))))
		if len(a.Tasks) == 0 {
			table.SetCell(0, 0,
				tview.NewTableCell("No tasks in list").
					SetTextColor(tcell.ColorWhite).
					SetAlign(tview.AlignCenter).SetSelectable(false))
		}
		for r := 0; r < int(rows); r++ {
			for c := 0; c < int(cols); c++ {
				color := tcell.ColorWhite
				if c < 1 || r < 1 {
					color = tcell.ColorYellow
				}
				text := ""
				if word < len(a.Tasks) {
					text = a.Tasks[a.InsertionOrder[word]].String()
				}
				table.SetCell(r, c,
					tview.NewTableCell(text).
						SetTextColor(color).
						SetAlign(tview.AlignCenter))
				word = (word + 1)
			}
		}
		table.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEscape {
				app.Stop()
			}
			if key == tcell.KeyEnter {
				app.Stop()
				// table.SetSelectable(true, true)
			}
		}).SetSelectedFunc(func(row int, column int) {
			table.GetCell(row, column).SetTextColor(tcell.ColorRed)
			table.SetSelectable(false, false)
		})
		if err := app.SetRoot(table, true).EnableMouse(true).Run(); err != nil {
			panic(err)
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
		var selection string
		app := tview.NewApplication()
		list := tview.NewList().
			AddItem(string(options[0]), "", 'a', func() {
				app.Stop()
				selection = string(options[0])
			}).
			AddItem(string(options[1]), "", 'b', func() {
				app.Stop()
				selection = string(options[1])
			}).
			AddItem(string(options[2]), "", 'c', func() {
				app.Stop()
				selection = string(options[2])
			}).
			AddItem(string(options[3]), "", 'd', func() {
				app.Stop()
				selection = string(options[3])
			}).
			AddItem(string(options[4]), "", 'e', func() {
				app.Stop()
				selection = string(options[4])
			}).
			AddItem("Quit", "", 'q', func() {
				app.Stop()
				selection = string(options[5])
			})
		if err := app.SetRoot(list, true).EnableMouse(true).Run(); err != nil {
			panic(err)
		}
		err = handleOption(a, stringWrapper(selection))
		if err != nil {
			if err == io.EOF {
				fmt.Println("Cancelled adding a task")
				continue
			}
		}
	}
}
