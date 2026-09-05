package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/tomhockett/task-cli/task"
)

type CLI struct {
	store task.TaskStore
	out   io.Writer
}

func NewCLI(store task.TaskStore, out io.Writer) *CLI {
	return &CLI{
		store: store,
		out:   out,
	}
}

func (c *CLI) Run(s []string) error {
	if len(s) == 0 {
		return fmt.Errorf("usage: task <command> [args]\n\ncommands: add, list, done, delete")
	}
	taskCommand := s[0]

	switch taskCommand {
	case "add":
		return c.runAdd(s[1:])
	case "list":
		return c.runList(s[1:])
	case "done":
		return c.runDone(s[1:])
	case "delete":
		return c.runDelete(s[1:])
	default:
		return fmt.Errorf("unknown command %q (want: add, list, done, delete)", taskCommand)
	}
}

func (c *CLI) runAdd(args []string) error {
	addFlags := flag.NewFlagSet("add", flag.ContinueOnError)
	priority := addFlags.String("priority", "medium", "Priority: low, medium, high")
	var tags []string
	addFlags.Func("tag", "Add a tag (repeatable)", func(s string) error {
		tags = append(tags, s)
		return nil
	})

	// Parse flags. %w keeps flag's own error inspectable (errors.Is(err, flag.ErrHelp))
	// while still adding context about which command failed.
	if err := addFlags.Parse(args); err != nil {
		return fmt.Errorf("add: %w", err)
	}

	// Remaining args after flags are the title
	rest := addFlags.Args()
	if len(rest) == 0 {
		return errors.New(`missing task title (usage: task add "Buy groceries")`)
	}
	title := strings.Join(rest, " ") // join multi-word titles

	taskPriority, err := task.ParsePriority(*priority)
	if err != nil {
		return fmt.Errorf("add: %w", err)
	}

	opts := task.AddOptions{
		Tags:     tags,
		Priority: &taskPriority,
	}
	t, err := c.store.Add(title, opts)
	if err != nil {
		return fmt.Errorf("adding task: %w", err)
	}
	fmt.Fprintf(c.out, "Added task %d\n", t.ID)
	return nil
}

func (c *CLI) runList(args []string) error {
	listFlags := flag.NewFlagSet("list", flag.ContinueOnError)
	status := listFlags.String("status", "", "Filter by status: todo, done")
	tag := listFlags.String("tag", "", "Filter by tag")
	format := listFlags.String("format", "table", "Output format: table, json")
	if err := listFlags.Parse(args); err != nil {
		return fmt.Errorf("list: %w", err)
	}

	var opts task.ListOptions
	if *status != "" {
		taskStatus, err := task.ParseStatus(*status)
		if err != nil {
			return fmt.Errorf("list: %w", err)
		}
		opts.Status = &taskStatus
	}
	opts.Tag = *tag

	tasks, err := c.store.List(opts)
	if err != nil {
		return fmt.Errorf("listing tasks: %w", err)
	}
	output, err := formatTasks(tasks, *format)
	if err != nil {
		return fmt.Errorf("list: %w", err)
	}
	fmt.Fprint(c.out, output)
	return nil
}

func (c *CLI) runDone(args []string) error {
	id, err := c.commandID("done", args)
	if err != nil {
		return err
	}
	if err := c.store.Complete(id); err != nil {
		return wrapStoreError("completing", id, err)
	}
	fmt.Fprintf(c.out, "Completed task %d\n", id)
	return nil
}

func (c *CLI) runDelete(args []string) error {
	id, err := c.commandID("delete", args)
	if err != nil {
		return err
	}
	if err := c.store.Delete(id); err != nil {
		return wrapStoreError("deleting", id, err)
	}
	fmt.Fprintf(c.out, "Deleted task %d\n", id)
	return nil
}

// commandID pulls the single task ID argument that done and delete both take.
func (c *CLI) commandID(command string, args []string) (int, error) {
	if len(args) == 0 {
		return 0, fmt.Errorf("missing task ID (usage: task %s 1)", command)
	}
	return parseID(args[0])
}

// wrapStoreError turns a bare ErrTaskNotFound into a message naming the ID the
// user actually typed. Wrapping with %w means errors.Is(err, task.ErrTaskNotFound)
// still works for callers — the context is added, not substituted.
func wrapStoreError(verb string, id int, err error) error {
	if errors.Is(err, task.ErrTaskNotFound) {
		return fmt.Errorf("no task with ID %d: %w", id, err)
	}
	return fmt.Errorf("%s task %d: %w", verb, id, err)
}

// formatTasks picks the renderer for --format. Keeping the switch here means
// runList doesn't have to know how either format is produced.
func formatTasks(tasks []task.Task, format string) (string, error) {
	switch format {
	case "table":
		return FormatTaskTable(tasks), nil
	case "json":
		return FormatTaskJSON(tasks)
	default:
		return "", fmt.Errorf("%q is not a valid format (want: table, json)", format)
	}
}

// ErrInvalidID is a sentinel so callers can test for bad input with errors.Is
// instead of matching on message text.
var ErrInvalidID = errors.New("task IDs are whole numbers starting at 1")

func parseID(s string) (int, error) {
	// strconv's own error ("strconv.Atoi: parsing \"abc\": invalid syntax")
	// leaks implementation detail at users, so it's swapped for ErrInvalidID.
	// Wrapping with %w still leaves the sentinel reachable via errors.Is.
	id, err := strconv.Atoi(s)
	if err != nil || id < 1 {
		return 0, fmt.Errorf("%q is not a valid task ID: %w", s, ErrInvalidID)
	}
	return id, nil
}
