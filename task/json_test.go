package task_test

import (
	"encoding/json"
	"testing"

	"github.com/tomhockett/task-cli/task"
)

func TestParseStatus(t *testing.T) {
	for _, want := range []task.Status{task.StatusTodo, task.StatusDone} {
		got, err := task.ParseStatus(want.String())
		assertNoError(t, err)
		assertEqual(t, got, want)
	}

	_, err := task.ParseStatus("nonsense")
	assertError(t, err)
}

func TestParsePriority(t *testing.T) {
	for _, want := range []task.Priority{task.PriorityLow, task.PriorityMedium, task.PriorityHigh} {
		got, err := task.ParsePriority(want.String())
		assertNoError(t, err)
		assertEqual(t, got, want)
	}

	_, err := task.ParsePriority("nonsense")
	assertError(t, err)
}

func TestStatusJSONRoundTrip(t *testing.T) {
	encoded, err := json.Marshal(task.StatusDone)
	assertNoError(t, err)
	assertEqual(t, string(encoded), `"done"`)

	var decoded task.Status
	assertNoError(t, json.Unmarshal(encoded, &decoded))
	assertEqual(t, decoded, task.StatusDone)

	assertError(t, json.Unmarshal([]byte(`"nonsense"`), &decoded))
}

func TestPriorityJSONRoundTrip(t *testing.T) {
	encoded, err := json.Marshal(task.PriorityHigh)
	assertNoError(t, err)
	assertEqual(t, string(encoded), `"high"`)

	var decoded task.Priority
	assertNoError(t, json.Unmarshal(encoded, &decoded))
	assertEqual(t, decoded, task.PriorityHigh)

	assertError(t, json.Unmarshal([]byte(`"nonsense"`), &decoded))
}
