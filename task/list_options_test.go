package task_test

import (
	"testing"

	"github.com/tomhockett/task-cli/task"
)

func TestInMemoryStore_ListByStatus(t *testing.T) {
	store := task.NewInMemoryTaskStore()
	store.Add("Task 1", task.AddOptions{})
	store.Add("Task 2", task.AddOptions{})
	store.Complete(1)

	done := task.StatusDone
	tasks, err := store.List(task.ListOptions{Status: &done})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("got %d tasks, want 1", len(tasks))
	}
	if tasks[0].Status != task.StatusDone {
		t.Errorf("got status %v, want StatusDone", tasks[0].Status)
	}
}

func TestInMemoryStore_ListByTag(t *testing.T) {
	store := task.NewInMemoryTaskStore()
	store.Add("Work task", task.AddOptions{Tags: []string{"work"}})
	store.Add("Home task", task.AddOptions{Tags: []string{"home"}})
	store.Add("Work urgent", task.AddOptions{Tags: []string{"work", "urgent"}})

	tasks, err := store.List(task.ListOptions{Tag: "work"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tasks) != 2 {
		t.Fatalf("got %d tasks with tag 'work', want 2", len(tasks))
	}
}

func TestSQLiteStore_ListByStatus(t *testing.T) {
	store := newTestStore(t)
	store.Add("Task 1", task.AddOptions{})
	store.Add("Task 2", task.AddOptions{})
	store.Complete(1)

	done := task.StatusDone
	tasks, err := store.List(task.ListOptions{Status: &done})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("got %d tasks, want 1", len(tasks))
	}
	if tasks[0].Status != task.StatusDone {
		t.Errorf("got status %v, want StatusDone", tasks[0].Status)
	}
}

func TestSQLiteStore_ListByTag(t *testing.T) {
	store := newTestStore(t)
	store.Add("Work task", task.AddOptions{Tags: []string{"work"}})
	store.Add("Home task", task.AddOptions{Tags: []string{"home"}})
	store.Add("Work urgent", task.AddOptions{Tags: []string{"work", "urgent"}})

	tasks, err := store.List(task.ListOptions{Tag: "work"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tasks) != 2 {
		t.Fatalf("got %d tasks with tag 'work', want 2", len(tasks))
	}
}
