package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aidarkhanov/nanoid"
	ki "todo.com/keypressinterface"
)

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
func newTask(name string, completed bool) *task {
	var taskId string
	taskId, err := nanoid.Generate(nanoid.DefaultAlphabet, 20)
	if err != nil {
		log.Fatal(err)
	}
	t := &task{Id: taskId, Name: name, Completed: completed}
	return t
}

func main() {
	tasks := []*task{
		newTask("grocery shopping", false),
		newTask("email justin", false),
		newTask("meal prepping", false),
		newTask("laundry", false),
		newTask("work on project", false),
		newTask("exercise", false),
		newTask("read", false),
		newTask("write", false),
		newTask("walk the dog", false),
		newTask("pay bills", false),
		newTask("home repairs", false),
	}

	menu, err := ki.NewMatrixMenu(tasks, int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal(err.Error())
	}
	selection, err := menu.RenderInterface()
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(selection)
}
