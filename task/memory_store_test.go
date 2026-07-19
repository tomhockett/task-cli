package task_test

import (
	"testing"

	"github.com/tomhockett/task-cli/task"
)

func TestInMemoryStore_AddAndList(t *testing.T) {
	store := task.NewInMemoryTaskStore()

	// Add first task
	t1, err := store.Add("Buy groceries", task.AddOptions{})
	assertNoError(t, err)
	assertEqual(t, t1.ID, 1)
	assertEqual(t, t1.Title, "Buy groceries")
	assertEqual(t, t1.Status, task.StatusTodo)

	// Add second task — ID auto-increments
	t2, err := store.Add("Walk the dog", task.AddOptions{})
	assertNoError(t, err)
	assertEqual(t, t2.ID, 2)
	assertEqual(t, t2.Title, "Walk the dog")

	// List returns both
	tasks, err := store.List(task.ListOptions{})
	assertNoError(t, err)
	assertEqual(t, len(tasks), 2)
}

func TestInMemoryStore_ListEmpty(t *testing.T) {
	store := task.NewInMemoryTaskStore()

	tasks, err := store.List(task.ListOptions{})
	assertNoError(t, err)
	assertEqual(t, len(tasks), 0)
}

func TestInMemoryStore_Complete(t *testing.T) {
	store := task.NewInMemoryTaskStore()
	store.Add("Buy groceries", task.AddOptions{})

	err := store.Complete(1)
	assertNoError(t, err)

	tasks, _ := store.List(task.ListOptions{})
	assertEqual(t, tasks[0].Status, task.StatusDone)
	if tasks[0].CompletedAt == nil {
		t.Error("expected CompletedAt to be set after completing")
	}
}

func TestInMemoryStore_CompleteNotFound(t *testing.T) {
	store := task.NewInMemoryTaskStore()

	err := store.Complete(999)
	assertErrorIs(t, err, task.ErrTaskNotFound)
}

func TestInMemoryStore_Delete(t *testing.T) {
	store := task.NewInMemoryTaskStore()
	store.Add("Buy groceries", task.AddOptions{})
	store.Add("Walk the dog", task.AddOptions{})

	err := store.Delete(1)
	assertNoError(t, err)

	tasks, _ := store.List(task.ListOptions{})
	assertEqual(t, len(tasks), 1)
	assertEqual(t, tasks[0].ID, 2)
}

func TestInMemoryStore_DeleteNotFound(t *testing.T) {
	store := task.NewInMemoryTaskStore()

	err := store.Delete(999)
	assertErrorIs(t, err, task.ErrTaskNotFound)
}
