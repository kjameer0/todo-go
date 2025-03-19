package main

import (
	"log"
	"os"

	ki "todo.com/keypressinterface"
)


func main() {
	tasks := []string{"grocery shopping", "meal prepping", "laundry", "work on project", "exercise", "read", "write", "walk the dog", "call family", "meditate", "pay bills", "clean the house", "check emails", "water plants", "car maintenance", "home repairs"}

	menu, err := ki.NewMatrixMenu(tasks, 3, 9, int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal(err.Error())
	}
	menu.RenderInterface()
}
