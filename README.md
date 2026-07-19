# task-cli

A simple task manager for your terminal, written in Go and backed by SQLite. Add tasks with priorities and tags, list and filter them, mark them done, and delete them — your tasks persist between runs.

## Requirements

- [Go](https://go.dev/dl/) 1.25 or newer

That's it. The SQLite driver ([modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite)) is pure Go, so there's no cgo toolchain or system SQLite installation needed.

## Installation

Install directly with `go install`:

```bash
go install github.com/tomhockett/task-cli/cmd/task@latest
```

Or clone and build from source:

```bash
git clone https://github.com/tomhockett/task-cli.git
cd task-cli
go build -o task ./cmd/task
```

## Setup

None required. On first run, `task` creates its database automatically at:

```
~/.task-cli/tasks.db
```

To start fresh, delete that file.

## Usage

```
task <command> [flags] [args]
```

### Commands

| Command | Description | Example |
|---------|-------------|---------|
| `add <title>` | Add a new task | `task add Buy groceries` |
| `list` | List tasks | `task list` |
| `done <id>` | Mark a task as done | `task done 1` |
| `delete <id>` | Delete a task | `task delete 1` |

### `add` flags

| Flag | Values | Default | Description |
|------|--------|---------|-------------|
| `--priority` | `low`, `medium`, `high` | `medium` | Set the task's priority |
| `--tag` | any string | — | Tag the task (repeatable) |

```bash
task add --priority high --tag work --tag urgent Finish the quarterly report
```

Flags come before the title; everything after the flags is joined into the title, so quotes are optional.

### `list` flags

| Flag | Values | Default | Description |
|------|--------|---------|-------------|
| `--status` | `todo`, `done` | show all | Filter by status |
| `--tag` | any string | show all | Filter to tasks with an exact matching tag |

```bash
task list --status done
task list --tag work
```

### Example session

```bash
$ task add Buy groceries
Added task 1
$ task add --priority high --tag work Prepare demo
Added task 2
$ task list
1, Buy groceries, todo
2, Prepare demo, todo
$ task done 1
$ task list --status done
1, Buy groceries, done
```

## Development

The codebase is fully testable without touching your real database — the CLI writes to an injected `io.Writer`, stores implement a common `TaskStore` interface (in-memory and SQLite), and SQLite tests run against temp files via `t.TempDir()`.

```bash
go test ./...        # run all tests
go test -v ./...     # verbose
go vet ./...         # static checks
go run ./cmd/task    # run without installing
```

### Project layout

```
cmd/task/   main entry point — wires the SQLite store into the CLI
cli/        command dispatch, flag parsing, and output formatting
task/       domain types (Task, Status, Priority) and the TaskStore
            interface with in-memory and SQLite implementations
```
