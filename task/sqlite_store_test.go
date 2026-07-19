package task_test

import (
	"path/filepath"
	"testing"

	"github.com/tomhockett/task-cli/task"
)

// newTestStore creates a SQLiteStore backed by a temporary file.
// t.TempDir() is automatically cleaned up after the test — like DatabaseCleaner in Rails.
func newTestStore(t *testing.T) *task.SQLiteStore {
	t.Helper()
	path := filepath.Join(t.TempDir(), "test.db")
	store, err := task.NewSQLiteStore(path)
	if err != nil {
		t.Fatalf("failed to create SQLiteStore: %v", err)
	}
	return store
}

func TestSQLiteStore_Open(t *testing.T) {
	store := newTestStore(t)
	if store == nil {
		t.Fatal("expected a non-nil store")
	}
}

func TestSQLiteStore_Add(t *testing.T) {
	store := newTestStore(t)

	tsk, err := store.Add("Buy groceries", task.AddOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tsk.ID != 1 {
		t.Errorf("got ID %d, want 1", tsk.ID)
	}
	if tsk.Title != "Buy groceries" {
		t.Errorf("got title %q, want %q", tsk.Title, "Buy groceries")
	}
	if tsk.Status != task.StatusTodo {
		t.Errorf("got status %v, want StatusTodo", tsk.Status)
	}
}

func TestSQLiteStore_List(t *testing.T) {
	store := newTestStore(t)

	store.Add("Buy groceries", task.AddOptions{})
	store.Add("Walk the dog", task.AddOptions{})

	tasks, err := store.List(task.ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tasks) != 2 {
		t.Fatalf("got %d tasks, want 2", len(tasks))
	}
	if tasks[0].Title != "Buy groceries" {
		t.Errorf("got title %q, want %q", tasks[0].Title, "Buy groceries")
	}
	if tasks[1].Title != "Walk the dog" {
		t.Errorf("got title %q, want %q", tasks[1].Title, "Walk the dog")
	}
}

func TestSQLiteStore_Complete(t *testing.T) {
	store := newTestStore(t)
	store.Add("Buy groceries", task.AddOptions{})

	err := store.Complete(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tasks, _ := store.List(task.ListOptions{})
	if tasks[0].Status != task.StatusDone {
		t.Errorf("got status %v, want StatusDone", tasks[0].Status)
	}
	if tasks[0].CompletedAt == nil {
		t.Error("expected CompletedAt to be set after completing")
	}
}

func TestSQLiteStore_Complete_NotFound(t *testing.T) {
	store := newTestStore(t)

	err := store.Complete(999)
	if err != task.ErrTaskNotFound {
		t.Errorf("got error %v, want ErrTaskNotFound", err)
	}
}

func TestSQLiteStore_Delete(t *testing.T) {
	store := newTestStore(t)
	store.Add("Buy groceries", task.AddOptions{})
	store.Add("Walk the dog", task.AddOptions{})

	err := store.Delete(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tasks, _ := store.List(task.ListOptions{})
	if len(tasks) != 1 {
		t.Fatalf("got %d tasks after delete, want 1", len(tasks))
	}
	if tasks[0].ID != 2 {
		t.Errorf("remaining task ID: got %d, want 2", tasks[0].ID)
	}
}

func TestSQLiteStore_Delete_NotFound(t *testing.T) {
	store := newTestStore(t)

	err := store.Delete(999)
	if err != task.ErrTaskNotFound {
		t.Errorf("got error %v, want ErrTaskNotFound", err)
	}
}
