package main

import (
	"encoding/json"
	"fmt"

	"github.com/aidarkhanov/nanoid"
	"golang.org/x/term"
)

type task struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
type app struct {
	tasks                 map[string]*task
	insertionOrder        []string
	originalTerminalState *term.State
	t                     *term.Terminal
}
type response1 struct {
	Page   int
	Fruits []string
}
type response2 struct {
	Page   int `json: "page"`
	Fruits int `json: "fruits"`
}

func newApp() *app {
	return &app{}
}
func newTask(id string, name string) *task {
	t := &task{Id: id, Name: name}
	return t
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
func main() {
	a := newApp()
	a.tasks = createTaskMap(a, []string{"check phone", "look at phone", "put phone down"})
	// disable input buffering
	bolB, _ := json.Marshal(true)
	fmt.Println(string(bolB))
	intB, _ := json.Marshal(1)
	fmt.Println(string(intB))
	sl := []string{"hi", "there"}
	slJson, _ := json.Marshal(sl)
	fmt.Println(string(slJson))
	taskJson, _ := json.Marshal(a.tasks)
	fmt.Println(string(taskJson))
}
