package cli

import (
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
		return fmt.Errorf("usage: task <command> [args]")
	}
	taskCommand := s[0]

	switch taskCommand {
	case "add":
		addFlags := flag.NewFlagSet("add", flag.ContinueOnError)
		priority := addFlags.String("priority", "medium", "Priority: low, medium, high")
		var tags []string
		addFlags.Func("tag", "Add a tag (repeatable)", func(s string) error {
			tags = append(tags, s)
			return nil
		})

		// Parse flags
		if err := addFlags.Parse(s[1:]); err != nil {
			return err
		}

		// Remaining args after flags are the title
		args := addFlags.Args()
		if len(args) == 0 {
			return fmt.Errorf("missing task title")
		}
		title := strings.Join(args, " ") // join multi-word titles
		taskPriority, err := task.ParsePriority(*priority)
		if err != nil {
			return err
		}

		opts := task.AddOptions{
			Tags:     tags,
			Priority: &taskPriority,
		}
		t, err := c.store.Add(title, opts)
		if err != nil {
			return err
		}
		fmt.Fprintf(c.out, "Added task %d\n", t.ID)
	case "list":
		listFlags := flag.NewFlagSet("list", flag.ContinueOnError)
		status := listFlags.String("status", "", "Filter by status: todo, done")
		tag := listFlags.String("tag", "", "Filter by tag")
		format := listFlags.String("format", "table", "Output format: table, json")
		if err := listFlags.Parse(s[1:]); err != nil {
			return err
		}

		var opts task.ListOptions
		if *status != "" {
			taskStatus, err := task.ParseStatus(*status)
			if err != nil {
				return err
			}
			opts.Status = &taskStatus
		}
		opts.Tag = *tag

		tasks, err := c.store.List(opts)
		if err != nil {
			return err
		}
		output, err := formatTasks(tasks, *format)
		if err != nil {
			return err
		}
		fmt.Fprint(c.out, output)
	case "done":
		if len(s) < 2 {
			return fmt.Errorf("missing task ID")
		}
		id, err := c.parseID(s[1])
		if err != nil {
			return err
		}
		err = c.store.Complete(id)
		if err != nil {
			return err
		}
	case "delete":
		if len(s) < 2 {
			return fmt.Errorf("missing task ID")
		}
		id, err := c.parseID(s[1])
		if err != nil {
			return err
		}
		err = c.store.Delete(id)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown command %q", taskCommand)
	}
	return nil
}

// formatTasks picks the renderer for --format. Keeping the switch here means
// Run doesn't have to know how either format is produced.
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

func (c *CLI) parseID(s string) (int, error) {
	id, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid task ID: %q", s)
	}
	return id, nil
}
