package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tomhockett/task-cli/task"
)

func FormatTaskTable(tasks []task.Task) string {
	if len(tasks) == 0 {
		return "No tasks\n"
	}
	var sb strings.Builder
	for _, t := range tasks {
		fmt.Fprintf(&sb, "%d, %s, %s\n", t.ID, t.Title, t.Status)
	}
	return sb.String()
}

// FormatTaskJSON renders tasks as an indented JSON array, ready to pipe into jq.
// A nil slice marshals to "null", which is awkward for consumers, so an empty
// result is normalized to "[]" first.
func FormatTaskJSON(tasks []task.Task) (string, error) {
	if tasks == nil {
		tasks = []task.Task{}
	}
	encoded, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return "", fmt.Errorf("encoding tasks as JSON: %w", err)
	}
	return string(encoded) + "\n", nil
}
