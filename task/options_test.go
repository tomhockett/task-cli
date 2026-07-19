package task_test

import (
	"testing"

	"github.com/tomhockett/task-cli/task"
)

func TestInMemoryStore_AddWithPriority(t *testing.T) {
	store := task.NewInMemoryTaskStore()
	high := task.PriorityHigh

	tsk, err := store.Add("Urgent task", task.AddOptions{Priority: &high})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tsk.Priority != task.PriorityHigh {
		t.Errorf("got priority %v, want PriorityHigh", tsk.Priority)
	}
}

func TestInMemoryStore_AddWithTags(t *testing.T) {
	store := task.NewInMemoryTaskStore()

	tsk, err := store.Add("Work task", task.AddOptions{Tags: []string{"work", "urgent"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tsk.Tags) != 2 {
		t.Fatalf("got %d tags, want 2", len(tsk.Tags))
	}
	if tsk.Tags[0] != "work" {
		t.Errorf("got tag %q, want %q", tsk.Tags[0], "work")
	}
	if tsk.Tags[1] != "urgent" {
		t.Errorf("got tag %q, want %q", tsk.Tags[1], "urgent")
	}
}

func TestSQLiteStore_AddWithPriority(t *testing.T) {
	store := newTestStore(t)
	high := task.PriorityHigh

	tsk, err := store.Add("Urgent task", task.AddOptions{Priority: &high})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tsk.Priority != task.PriorityHigh {
		t.Errorf("got priority %v, want PriorityHigh", tsk.Priority)
	}

	// Verify it persisted correctly
	tasks, err := store.List(task.ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error from List: %v", err)
	}
	if len(tasks) == 0 {
		t.Fatal("expected at least 1 task, got 0")
	}
	if tasks[0].Priority != task.PriorityHigh {
		t.Errorf("persisted priority: got %v, want PriorityHigh", tasks[0].Priority)
	}
}

func TestSQLiteStore_AddWithTags(t *testing.T) {
	store := newTestStore(t)

	tsk, err := store.Add("Work task", task.AddOptions{Tags: []string{"work", "urgent"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tsk.Tags) != 2 {
		t.Fatalf("got %d tags, want 2", len(tsk.Tags))
	}

	// Verify tags persisted correctly
	tasks, err := store.List(task.ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error from List: %v", err)
	}
	if len(tasks) == 0 {
		t.Fatal("expected at least 1 task, got 0")
	}
	if len(tasks[0].Tags) != 2 {
		t.Fatalf("persisted tags: got %d, want 2", len(tasks[0].Tags))
	}
	if tasks[0].Tags[0] != "work" {
		t.Errorf("persisted tag: got %q, want %q", tasks[0].Tags[0], "work")
	}
}
