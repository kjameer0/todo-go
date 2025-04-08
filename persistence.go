package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

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
		if errors.Is(err, os.ErrNotExist) {
			return
		}
		log.Fatal("failed to read from save location", err)
	}
	json.Unmarshal(data, &s)
	a.InsertionOrder = s.InsertionOrder
	if len(s.Tasks) > 0 {
		a.Tasks = s.Tasks
	}
}
