package cli_test

import (
	"bytes"
	"testing"

	"github.com/tomhockett/task-cli/cli"
	"github.com/tomhockett/task-cli/task"
)

// newTestCLI wires up a CLI with an in-memory store and a buffer to capture output.
// Like injecting a StringIO in a Rails service test — no real I/O needed.
func newTestCLI() (*cli.CLI, *bytes.Buffer) {
	store := task.NewInMemoryTaskStore()
	buf := &bytes.Buffer{}
	return cli.NewCLI(store, buf), buf
}

func TestCLI_Add(t *testing.T) {
	c, buf := newTestCLI()

	err := c.Run([]string{"add", "Buy groceries"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Output should confirm the task was created
	got := buf.String()
	want := "Added task 1\n"
	if got != want {
		t.Errorf("got output %q, want %q", got, want)
	}
}

func TestCLI_Add_MissingTitle(t *testing.T) {
	c, _ := newTestCLI()

	err := c.Run([]string{"add"})
	if err == nil {
		t.Fatal("expected an error but got nil")
	}
}

func TestCLI_Add_WithPriority(t *testing.T) {
	c, _ := newTestCLI()

	err := c.Run([]string{"add", "--priority", "high", "Urgent task"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCLI_Add_WithInvalidPriority(t *testing.T) {
	c, _ := newTestCLI()

	err := c.Run([]string{"add", "--priority", "nonsense", "Task"})
	if err == nil {
		t.Fatal("expected an error for invalid priority but got nil")
	}
}

func TestCLI_Add_WithTags(t *testing.T) {
	c, _ := newTestCLI()

	err := c.Run([]string{"add", "--tag", "work", "--tag", "urgent", "Work task"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCLI_List(t *testing.T) {
	c, buf := newTestCLI()

	c.Run([]string{"add", "Buy groceries"})
	c.Run([]string{"add", "Walk the dog"})

	buf.Reset() // clear add output, only care about list output
	err := c.Run([]string{"list"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	if !contains(got, "Buy groceries") {
		t.Errorf("output missing %q:\n%s", "Buy groceries", got)
	}
	if !contains(got, "Walk the dog") {
		t.Errorf("output missing %q:\n%s", "Walk the dog", got)
	}
	if !contains(got, "1") {
		t.Errorf("output missing task ID 1:\n%s", got)
	}
}

func TestCLI_List_Empty(t *testing.T) {
	c, buf := newTestCLI()

	err := c.Run([]string{"list"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	want := "No tasks\n"
	if got != want {
		t.Errorf("got output %q, want %q", got, want)
	}
}

func TestCLI_List_FilterByStatus(t *testing.T) {
	c, buf := newTestCLI()

	c.Run([]string{"add", "Buy groceries"})
	c.Run([]string{"add", "Walk the dog"})
	c.Run([]string{"done", "1"})

	buf.Reset()
	err := c.Run([]string{"list", "--status", "done"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	if !contains(got, "Buy groceries") {
		t.Errorf("output missing completed task %q:\n%s", "Buy groceries", got)
	}
	if contains(got, "Walk the dog") {
		t.Errorf("output should not include todo task %q:\n%s", "Walk the dog", got)
	}
}

func TestCLI_List_FilterByTag(t *testing.T) {
	c, buf := newTestCLI()

	c.Run([]string{"add", "--tag", "work", "Finish report"})
	c.Run([]string{"add", "--tag", "home", "Mow the lawn"})

	buf.Reset()
	err := c.Run([]string{"list", "--tag", "work"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	if !contains(got, "Finish report") {
		t.Errorf("output missing tagged task %q:\n%s", "Finish report", got)
	}
	if contains(got, "Mow the lawn") {
		t.Errorf("output should not include untagged task %q:\n%s", "Mow the lawn", got)
	}
}

func TestCLI_List_InvalidStatus(t *testing.T) {
	c, _ := newTestCLI()

	err := c.Run([]string{"list", "--status", "nonsense"})
	if err == nil {
		t.Fatal("expected an error for invalid status but got nil")
	}
}

func TestCLI_Done(t *testing.T) {
	store := task.NewInMemoryTaskStore()
	buf := &bytes.Buffer{}
	c := cli.NewCLI(store, buf)

	c.Run([]string{"add", "Buy groceries"})

	err := c.Run([]string{"done", "1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify state directly via the store
	tasks, _ := store.List(task.ListOptions{})
	if tasks[0].Status != task.StatusDone {
		t.Errorf("got status %v, want StatusDone", tasks[0].Status)
	}
}

func TestCLI_Done_InvalidID(t *testing.T) {
	c, _ := newTestCLI()

	err := c.Run([]string{"done", "abc"})
	if err == nil {
		t.Fatal("expected an error for non-numeric ID but got nil")
	}
}

func TestCLI_Done_NotFound(t *testing.T) {
	c, _ := newTestCLI()

	err := c.Run([]string{"done", "999"})
	if err == nil {
		t.Fatal("expected an error for missing task but got nil")
	}
}

func TestCLI_Done_MissingID(t *testing.T) {
	c, _ := newTestCLI()

	err := c.Run([]string{"done"})
	if err == nil {
		t.Fatal("expected an error when no ID provided but got nil")
	}
}

func TestCLI_Delete(t *testing.T) {
	store := task.NewInMemoryTaskStore()
	buf := &bytes.Buffer{}
	c := cli.NewCLI(store, buf)

	c.Run([]string{"add", "Buy groceries"})

	err := c.Run([]string{"delete", "1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tasks, _ := store.List(task.ListOptions{})
	if len(tasks) != 0 {
		t.Errorf("got %d tasks after delete, want 0", len(tasks))
	}
}

func TestCLI_Delete_InvalidID(t *testing.T) {
	c, _ := newTestCLI()

	err := c.Run([]string{"delete", "abc"})
	if err == nil {
		t.Fatal("expected an error for non-numeric ID but got nil")
	}
}

func TestCLI_Delete_NotFound(t *testing.T) {
	c, _ := newTestCLI()

	err := c.Run([]string{"delete", "999"})
	if err == nil {
		t.Fatal("expected an error for missing task but got nil")
	}
}

func TestCLI_Delete_MissingID(t *testing.T) {
	c, _ := newTestCLI()

	err := c.Run([]string{"delete"})
	if err == nil {
		t.Fatal("expected an error when no ID provided but got nil")
	}
}

func TestCLI_NoArgs(t *testing.T) {
	c, _ := newTestCLI()

	err := c.Run([]string{})
	if err == nil {
		t.Fatal("expected an error for no args but got nil")
	}
}

func TestCLI_UnknownCommand(t *testing.T) {
	c, _ := newTestCLI()

	err := c.Run([]string{"frobnicate"})
	if err == nil {
		t.Fatal("expected an error for unknown command but got nil")
	}
}

// contains is a small helper to avoid importing strings in the test file.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsAt(s, substr))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
