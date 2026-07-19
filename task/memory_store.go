package task

import (
	"time"
)

// InMemoryTaskStore holds tasks in a slice — great for testing, no I/O needed.
// Like using a plain Ruby array instead of hitting the database.
type InMemoryTaskStore struct {
	tasks  []Task
	nextID int
}

// Compile-time check that InMemoryTaskStore implements TaskStore.
// If InMemoryTaskStore is missing any interface methods, this line will fail to compile.
var _ TaskStore = (*InMemoryTaskStore)(nil)

// NewInMemoryTaskStore is a constructor function — Go's version of .new.
// Returns a pointer so all method calls share the same data.
func NewInMemoryTaskStore() *InMemoryTaskStore {
	return &InMemoryTaskStore{nextID: 1}
}

func (s *InMemoryTaskStore) Add(title string, opts AddOptions) (Task, error) {
	var priority Priority
	if opts.Priority != nil {
		priority = *opts.Priority
	} else {
		priority = PriorityMedium
	}
	t := Task{
		ID:        s.nextID,
		Title:     title,
		Status:    StatusTodo,
		Priority:  priority,
		Tags:      opts.Tags,
		CreatedAt: time.Now(),
	}
	s.tasks = append(s.tasks, t)
	s.nextID++
	return t, nil
}

func (s *InMemoryTaskStore) List(opts ListOptions) ([]Task, error) {
	var tasks []Task
	for _, t := range s.tasks {
		if opts.matches(t) {
			tasks = append(tasks, t)
		}
	}
	return tasks, nil
}

func (s *InMemoryTaskStore) Complete(id int) error {
	for index, t := range s.tasks {
		if t.ID == id {
			s.tasks[index].Status = StatusDone
			now := time.Now()
			s.tasks[index].CompletedAt = &now
			return nil
		}
	}
	return ErrTaskNotFound
}

func (s *InMemoryTaskStore) Delete(id int) error {
	for index, t := range s.tasks {
		if t.ID == id {
			s.tasks = append(s.tasks[:index], s.tasks[index+1:]...)
			return nil
		}
	}
	return ErrTaskNotFound
}
