package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/tomhockett/task-cli/cli"
	"github.com/tomhockett/task-cli/task"
)

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

// run is everything main does, minus the process-level concerns (reading
// os.Args, writing to stderr, exiting). Returning an error and taking an
// io.Writer is what makes the end-to-end test able to drive the real CLI
// against a real SQLite file.
func run(args []string, out io.Writer) error {
	path, err := databasePath()
	if err != nil {
		return err
	}

	store, err := task.NewSQLiteStore(path)
	if err != nil {
		return err
	}
	defer store.Close()

	return cli.NewCLI(store, out).Run(args)
}

// databasePath resolves ~/.task-cli/tasks.db, creating the directory if needed.
// os.UserHomeDir reads $HOME, so a test can point the whole app at a temp
// directory with t.Setenv("HOME", t.TempDir()).
func databasePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("finding home directory: %w", err)
	}

	dir := filepath.Join(home, ".task-cli")
	// 0755 is Unix file permissions (rwxr-xr-x)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("creating data directory %s: %w", dir, err)
	}

	return filepath.Join(dir, "tasks.db"), nil
}
