package cli_test

import (
	"errors"
	"testing"

	"github.com/tomhockett/task-cli/cli"
	"github.com/tomhockett/task-cli/task"
)

// assertErrorContains checks the message a user would actually read.
func assertErrorContains(t *testing.T, err error, want string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected an error containing %q but got nil", want)
	}
	if !contains(err.Error(), want) {
		t.Errorf("got error %q, want it to contain %q", err.Error(), want)
	}
}

func TestCLI_Done_InvalidID_Message(t *testing.T) {
	c, _ := newTestCLI()

	err := c.Run([]string{"done", "abc"})
	assertErrorContains(t, err, `"abc" is not a valid task ID`)
}

func TestCLI_Done_ZeroID(t *testing.T) {
	c, _ := newTestCLI()

	err := c.Run([]string{"done", "0"})
	assertErrorContains(t, err, "starting at 1")
	if !errors.Is(err, cli.ErrInvalidID) {
		t.Errorf("got error %v, want it to wrap cli.ErrInvalidID", err)
	}
}

func TestCLI_Done_InvalidID_HidesStrconvDetail(t *testing.T) {
	c, _ := newTestCLI()

	err := c.Run([]string{"done", "abc"})
	if contains(err.Error(), "strconv") {
		t.Errorf("got error %q, want it to hide strconv internals", err)
	}
	if !errors.Is(err, cli.ErrInvalidID) {
		t.Errorf("got error %v, want it to wrap cli.ErrInvalidID", err)
	}
}

func TestCLI_Done_NotFound_MessageAndUnwrap(t *testing.T) {
	c, _ := newTestCLI()

	err := c.Run([]string{"done", "999"})
	assertErrorContains(t, err, "no task with ID 999")

	// Wrapping with %w means the sentinel is still reachable.
	if !errors.Is(err, task.ErrTaskNotFound) {
		t.Errorf("got error %v, want it to wrap task.ErrTaskNotFound", err)
	}
}

func TestCLI_Delete_NotFound_MessageAndUnwrap(t *testing.T) {
	c, _ := newTestCLI()

	err := c.Run([]string{"delete", "999"})
	assertErrorContains(t, err, "no task with ID 999")

	if !errors.Is(err, task.ErrTaskNotFound) {
		t.Errorf("got error %v, want it to wrap task.ErrTaskNotFound", err)
	}
}

func TestCLI_MissingID_MentionsUsage(t *testing.T) {
	c, _ := newTestCLI()

	assertErrorContains(t, c.Run([]string{"done"}), "usage: task done 1")
	assertErrorContains(t, c.Run([]string{"delete"}), "usage: task delete 1")
}

func TestCLI_Add_MissingTitle_MentionsUsage(t *testing.T) {
	c, _ := newTestCLI()

	assertErrorContains(t, c.Run([]string{"add"}), "usage: task add")
}

func TestCLI_Add_InvalidPriority_ListsChoices(t *testing.T) {
	c, _ := newTestCLI()

	err := c.Run([]string{"add", "--priority", "nonsense", "Task"})
	assertErrorContains(t, err, `"nonsense" is not a valid priority (want: low, medium, high)`)
}

func TestCLI_List_InvalidStatus_ListsChoices(t *testing.T) {
	c, _ := newTestCLI()

	err := c.Run([]string{"list", "--status", "nonsense"})
	assertErrorContains(t, err, `"nonsense" is not a valid status (want: todo, done)`)
}

func TestCLI_List_InvalidFormat_ListsChoices(t *testing.T) {
	c, _ := newTestCLI()

	err := c.Run([]string{"list", "--format", "xml"})
	assertErrorContains(t, err, `"xml" is not a valid format (want: table, json)`)
}

func TestCLI_UnknownCommand_ListsCommands(t *testing.T) {
	c, _ := newTestCLI()

	err := c.Run([]string{"frobnicate"})
	assertErrorContains(t, err, `unknown command "frobnicate"`)
	assertErrorContains(t, err, "add, list, done, delete")
}

func TestCLI_Done_ReportsSuccess(t *testing.T) {
	c, buf := newTestCLI()

	c.Run([]string{"add", "Buy groceries"})
	buf.Reset()

	if err := c.Run([]string{"done", "1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got, want := buf.String(), "Completed task 1\n"; got != want {
		t.Errorf("got output %q, want %q", got, want)
	}
}

func TestCLI_Delete_ReportsSuccess(t *testing.T) {
	c, buf := newTestCLI()

	c.Run([]string{"add", "Buy groceries"})
	buf.Reset()

	if err := c.Run([]string{"delete", "1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got, want := buf.String(), "Deleted task 1\n"; got != want {
		t.Errorf("got output %q, want %q", got, want)
	}
}
