package task

import (
	"encoding/json"
	"fmt"
	"time"
)

// Task is the core domain struct — like an AR model but with explicit fields.
// CompletedAt is a pointer so it can be nil (Go's way of expressing "optional").
// The `json:"..."` struct tags control the JSON field names, the same way
// as_json(only: [...]) shapes a Rails model's JSON. "omitempty" drops the field
// entirely when it's nil or empty.
type Task struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Status      Status     `json:"status"`
	Priority    Priority   `json:"priority"`
	Tags        []string   `json:"tags,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// Status represents the state of a task — like an enum in Rails.
// Go uses iota to auto-increment integer constants within a const block.
type Status int

const (
	StatusTodo Status = iota
	StatusDone
)

func (s Status) String() string {
	switch s {
	case StatusDone:
		return "done"
	default:
		return "todo"
	}
}

// ParseStatus converts a user-supplied string into a Status.
// Both the CLI flags and JSON decoding go through here, so exactly one place
// knows which strings are valid and what to say when they aren't.
func ParseStatus(s string) (Status, error) {
	switch s {
	case "todo":
		return StatusTodo, nil
	case "done":
		return StatusDone, nil
	default:
		return 0, fmt.Errorf("%q is not a valid status (want: todo, done)", s)
	}
}

// MarshalJSON makes Status serialize as "todo"/"done" instead of 0/1.
// Implementing json.Marshaler is Go's version of overriding as_json on a model.
func (s Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// UnmarshalJSON is the other half of the pair, so JSON output can be decoded
// back into a Task. It needs a pointer receiver because it mutates the value.
func (s *Status) UnmarshalJSON(data []byte) error {
	var name string
	if err := json.Unmarshal(data, &name); err != nil {
		return fmt.Errorf("decoding status: %w", err)
	}
	parsed, err := ParseStatus(name)
	if err != nil {
		return err
	}
	*s = parsed
	return nil
}

// Priority represents how urgent a task is.
type Priority int

const (
	PriorityLow Priority = iota
	PriorityMedium
	PriorityHigh
)

func (p Priority) String() string {
	switch p {
	case PriorityMedium:
		return "medium"
	case PriorityHigh:
		return "high"
	default:
		return "low"
	}
}

// ParsePriority converts a user-supplied string into a Priority.
func ParsePriority(s string) (Priority, error) {
	switch s {
	case "low":
		return PriorityLow, nil
	case "medium":
		return PriorityMedium, nil
	case "high":
		return PriorityHigh, nil
	default:
		return 0, fmt.Errorf("%q is not a valid priority (want: low, medium, high)", s)
	}
}

func (p Priority) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

func (p *Priority) UnmarshalJSON(data []byte) error {
	var name string
	if err := json.Unmarshal(data, &name); err != nil {
		return fmt.Errorf("decoding priority: %w", err)
	}
	parsed, err := ParsePriority(name)
	if err != nil {
		return err
	}
	*p = parsed
	return nil
}
