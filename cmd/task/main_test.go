package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/tomhockett/task-cli/task"
)

// runTask invokes the real CLI against the real SQLite store, exactly as the
// binary does, and returns whatever it printed. Each call opens and closes the
// database, so it also proves data survives between "runs" of the program.
func runTask(t *testing.T, args ...string) string {
	t.Helper()
	var buf bytes.Buffer
	if err := run(args, &buf); err != nil {
		t.Fatalf("task %s: unexpected error: %v", strings.Join(args, " "), err)
	}
	return buf.String()
}

// runTaskExpectingError is the same, but for the failure paths.
func runTaskExpectingError(t *testing.T, args ...string) error {
	t.Helper()
	var buf bytes.Buffer
	err := run(args, &buf)
	if err == nil {
		t.Fatalf("task %s: expected an error but got nil", strings.Join(args, " "))
	}
	return err
}

// TestEndToEnd walks the whole feature set in one flow:
// add -> list -> done -> filter -> json -> delete.
func TestEndToEnd(t *testing.T) {
	// Point ~/.task-cli at a temp directory so the test never touches the
	// developer's real task list. t.Setenv restores $HOME when the test ends.
	t.Setenv("HOME", t.TempDir())

	// --- add ---
	if got, want := runTask(t, "add", "--priority", "high", "--tag", "work", "Finish the report"), "Added task 1\n"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
	if got, want := runTask(t, "add", "--tag", "home", "Mow the lawn"), "Added task 2\n"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}

	// --- list --- (a separate process would see both tasks; so does a separate run)
	listed := runTask(t, "list")
	for _, want := range []string{"Finish the report", "Mow the lawn", "todo"} {
		if !strings.Contains(listed, want) {
			t.Errorf("list output missing %q:\n%s", want, listed)
		}
	}

	// --- done ---
	if got, want := runTask(t, "done", "1"), "Completed task 1\n"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}

	// --- filter by status ---
	done := runTask(t, "list", "--status", "done")
	if !strings.Contains(done, "Finish the report") {
		t.Errorf("--status done should list the completed task:\n%s", done)
	}
	if strings.Contains(done, "Mow the lawn") {
		t.Errorf("--status done should not list the todo task:\n%s", done)
	}

	// --- filter by tag ---
	tagged := runTask(t, "list", "--tag", "home")
	if !strings.Contains(tagged, "Mow the lawn") {
		t.Errorf("--tag home should list the tagged task:\n%s", tagged)
	}
	if strings.Contains(tagged, "Finish the report") {
		t.Errorf("--tag home should not list the work task:\n%s", tagged)
	}

	// --- json ---
	var tasks []task.Task
	if err := json.Unmarshal([]byte(runTask(t, "list", "--format", "json")), &tasks); err != nil {
		t.Fatalf("--format json did not produce decodable JSON: %v", err)
	}
	if len(tasks) != 2 {
		t.Fatalf("got %d tasks from JSON, want 2", len(tasks))
	}
	first := tasks[0]
	if first.ID != 1 || first.Title != "Finish the report" {
		t.Errorf("got task %+v, want ID 1 titled %q", first, "Finish the report")
	}
	if first.Status != task.StatusDone {
		t.Errorf("got status %v, want done", first.Status)
	}
	if first.Priority != task.PriorityHigh {
		t.Errorf("got priority %v, want high", first.Priority)
	}
	if len(first.Tags) != 1 || first.Tags[0] != "work" {
		t.Errorf("got tags %v, want [work]", first.Tags)
	}
	if first.CompletedAt == nil {
		t.Error("expected CompletedAt to be set on the completed task")
	}
	if first.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}

	// --- delete ---
	if got, want := runTask(t, "delete", "2"), "Deleted task 2\n"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}

	remaining := runTask(t, "list")
	if strings.Contains(remaining, "Mow the lawn") {
		t.Errorf("deleted task still listed:\n%s", remaining)
	}
	if !strings.Contains(remaining, "Finish the report") {
		t.Errorf("surviving task missing from list:\n%s", remaining)
	}

	// --- errors reach the user intact through the real wiring ---
	err := runTaskExpectingError(t, "done", "999")
	if !strings.Contains(err.Error(), "no task with ID 999") {
		t.Errorf("got error %q, want it to name the missing ID", err)
	}
	if !errors.Is(err, task.ErrTaskNotFound) {
		t.Errorf("got error %v, want it to wrap task.ErrTaskNotFound", err)
	}
}

// TestEndToEnd_EmptyDatabase covers the first-run experience: no database file
// exists yet, and list should still work.
func TestEndToEnd_EmptyDatabase(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	if got, want := runTask(t, "list"), "No tasks\n"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
	if got, want := runTask(t, "list", "--format", "json"), "[]\n"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
