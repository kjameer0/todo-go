package main

import (
	"fmt"
	"time"
)

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
	a.InsertionOrder = []string{}
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
