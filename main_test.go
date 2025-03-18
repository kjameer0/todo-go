package main

import (
	"reflect"
	"testing"
)

func Test_removeTask(t *testing.T) {
	type args struct {
		tasks  *[]task
		taskId int
	}
	tests := []struct {
		name     string
		args     args
		expected []task
	}{
		{
			name: "standard",
			args: args{
				tasks: &[]task{
					{id: 1, name: "Task 1"},
					{id: 2, name: "Task 2"},
				},
				taskId: 1,
			},
			expected: []task{
				{id: 2, name: "Task 2"},
			},
		},
		{
			name: "non-existent task",
			args: args{
				tasks: &[]task{
					{id: 1, name: "Task 1"},
					{id: 2, name: "Task 2"},
				},
				taskId: 3,
			},
			expected: []task{
				{id: 1, name: "Task 1"},
				{id: 2, name: "Task 2"},
			}, // No change expected
		},
		{
			name: "empty list",
			args: args{
				tasks:  &[]task{},
				taskId: 1,
			},
			expected: []task{}, // No change, still empty
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			removeTask(tt.args.tasks, tt.args.taskId)
			if !reflect.DeepEqual(*tt.args.tasks, tt.expected) {
				t.Errorf("removeTask() got = %v, want = %v", *tt.args.tasks, tt.expected)
			}
		})
	}
}
