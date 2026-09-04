package cli_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/tomhockett/task-cli/cli"
	"github.com/tomhockett/task-cli/task"
)

func TestFormatTaskJSON(t *testing.T) {
	now := time.Now()
	completed := now.Add(time.Hour)
	tasks := []task.Task{
		{ID: 1, Title: "Buy groceries", Status: task.StatusTodo, Priority: task.PriorityMedium, Tags: []string{"home"}, CreatedAt: now},
		{ID: 2, Title: "Walk the dog", Status: task.StatusDone, Priority: task.PriorityHigh, CreatedAt: now, CompletedAt: &completed},
	}

	got, err := cli.FormatTaskJSON(tasks)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// The real assertion: the output decodes back into tasks.
	var decoded []task.Task
	if err := json.Unmarshal([]byte(got), &decoded); err != nil {
		t.Fatalf("output is not valid JSON: %v\n%s", err, got)
	}

	if len(decoded) != 2 {
		t.Fatalf("got %d tasks, want 2", len(decoded))
	}
	if decoded[0].Title != "Buy groceries" {
		t.Errorf("got title %q, want %q", decoded[0].Title, "Buy groceries")
	}
	if decoded[1].Status != task.StatusDone {
		t.Errorf("got status %v, want done", decoded[1].Status)
	}
	if decoded[1].Priority != task.PriorityHigh {
		t.Errorf("got priority %v, want high", decoded[1].Priority)
	}
	if decoded[1].CompletedAt == nil {
		t.Error("expected CompletedAt to survive the round trip")
	}

	// Enums should be human-readable strings, not raw integers.
	if !contains(got, `"status": "done"`) {
		t.Errorf("expected status rendered as a string:\n%s", got)
	}
	if !contains(got, `"priority": "high"`) {
		t.Errorf("expected priority rendered as a string:\n%s", got)
	}
}

func TestFormatTaskJSON_Empty(t *testing.T) {
	got, err := cli.FormatTaskJSON(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "[]\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestCLI_List_FormatJSON(t *testing.T) {
	c, buf := newTestCLI()

	c.Run([]string{"add", "--priority", "high", "--tag", "work", "Finish report"})

	buf.Reset()
	if err := c.Run([]string{"list", "--format", "json"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var decoded []task.Task
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("output is not valid JSON: %v\n%s", err, buf.String())
	}
	if len(decoded) != 1 {
		t.Fatalf("got %d tasks, want 1", len(decoded))
	}
	if decoded[0].Title != "Finish report" {
		t.Errorf("got title %q, want %q", decoded[0].Title, "Finish report")
	}
	if decoded[0].Priority != task.PriorityHigh {
		t.Errorf("got priority %v, want high", decoded[0].Priority)
	}
	if len(decoded[0].Tags) != 1 || decoded[0].Tags[0] != "work" {
		t.Errorf("got tags %v, want [work]", decoded[0].Tags)
	}
}

func TestCLI_List_FormatJSON_Empty(t *testing.T) {
	c, buf := newTestCLI()

	if err := c.Run([]string{"list", "--format", "json"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "[]\n"
	if got := buf.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestCLI_List_InvalidFormat(t *testing.T) {
	c, _ := newTestCLI()

	err := c.Run([]string{"list", "--format", "xml"})
	if err == nil {
		t.Fatal("expected an error for an unknown format but got nil")
	}
}
