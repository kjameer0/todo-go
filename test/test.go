// package main

// import (
// 	"encoding/json"
// 	"fmt"

// 	"github.com/aidarkhanov/nanoid"
// 	"golang.org/x/term"
// )

// type app struct {
// 	tasks                 map[string]*task
// 	insertionOrder        []string
// 	originalTerminalState *term.State
// 	t                     *term.Terminal
// }
// type response1 struct {
// 	Page   int
// 	Fruits []string
// }
// type response2 struct {
// 	Page   int `json: "page"`
// 	Fruits int `json: "fruits"`
// }

// func newApp() *app {
// 	return &app{}
// }
// func newTask(id string, name string) *task {
// 	t := &task{Id: id, Name: name}
// 	return t
// }

// func createTaskMap(a *app, tasks []string) map[string]*task {
// 	taskMap := make(map[string]*task)
// 	for _, taskText := range tasks {
// 		nano, _ := nanoid.Generate(nanoid.DefaultAlphabet, 5)
// 		taskMap[nano] = newTask(nano, taskText)
// 		a.insertionOrder = append(a.insertionOrder, nano)
// 	}
// 	return taskMap
// }
// func main() {
// 	a := newApp()
// 	a.tasks = createTaskMap(a, []string{"check phone", "look at phone", "put phone down"})
// 	// disable input buffering
// 	bolB, _ := json.Marshal(true)
// 	fmt.Println(string(bolB))
// 	intB, _ := json.Marshal(1)
// 	fmt.Println(string(intB))
// 	sl := []string{"hi", "there"}
// 	slJson, _ := json.Marshal(sl)
// 	fmt.Println(string(slJson))
// 	taskJson, _ := json.Marshal(a.tasks)
// 	fmt.Println(string(taskJson))
// }

package main

import (
	"fmt"
	"io"
	"log"
	"os"

	navmenu "todo.com/nav-menu"
)

type task struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}

func (t task) String() string {
	var completed string
	if !t.Completed {
		completed = "❌"
	} else {
		completed = "✅"
	}
	return fmt.Sprintf("%s %s", t.Name, completed)
}

func main() {
	// Create a base menu
	menu := []task{
		{Id: "1", Name: "Say hi", Completed: false},
		{Id: "2", Name: "Check emails", Completed: false},
		{Id: "3", Name: "Write code", Completed: true},
		{Id: "4", Name: "Review PRs #1", Completed: false},
		{Id: "5", Name: "Review PRs #2", Completed: false},
	}

	// Call NewMenu with the menu
	m := navmenu.NewMenu(menu, int(os.Stdin.Fd()))
	t, err := m.Render()
	if err != nil {
		if err != io.EOF {
			log.Fatal("failure in rendering menu: ", err)
		}
	}
	fmt.Println(t)

}
