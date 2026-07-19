// Package task defines the TaskStore interface and the sentinel error ErrTaskNotFound.
package task

import "errors"

// Sentinel errors — like custom exception classes in Ruby.
// Callers check with errors.Is(err, ErrTaskNotFound).
var ErrTaskNotFound = errors.New("task not found")

// TaskStore is the interface both InMemoryTaskStore and SQLiteStore will implement.
// In Rails terms, this is the contract that ActiveRecord provides implicitly.
// In Go, you define it explicitly so you can swap implementations.
type TaskStore interface {
	Add(title string, opts AddOptions) (Task, error)
	List(opts ListOptions) ([]Task, error)
	Complete(id int) error
	Delete(id int) error
}
