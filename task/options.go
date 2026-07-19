package task

import "slices"

// AddOptions allows specifying optional parameters when creating a task.
// The pointer-for-optional pattern: if Priority is nil, use the default (PriorityMedium).
type AddOptions struct {
	Priority *Priority
	Tags     []string
}

// ListOptions allows filtering tasks when listing.
// Pointer fields use the pointer-for-optional pattern: nil means "no filter".
type ListOptions struct {
	Status *Status
	Tag    string // filter to tasks with this tag (empty string means no filter)
}

// matches reports whether t passes the filters in opts.
// Like an AR scope that's a no-op when the param is absent.
func (opts ListOptions) matches(t Task) bool {
	if opts.Status != nil && t.Status != *opts.Status {
		return false
	}
	if opts.Tag != "" && !slices.Contains(t.Tags, opts.Tag) {
		return false
	}
	return true
}
