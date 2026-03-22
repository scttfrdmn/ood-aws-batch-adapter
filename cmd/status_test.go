package cmd

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/batch/types"
)

func TestBatchStateToOod(t *testing.T) {
	tests := []struct {
		state    types.JobStatus
		expected string
	}{
		{types.JobStatusSubmitted, "queued"},
		{types.JobStatusPending, "queued"},
		{types.JobStatusRunnable, "queued"},
		{types.JobStatusStarting, "running"},
		{types.JobStatusRunning, "running"},
		{types.JobStatusSucceeded, "completed"},
		{types.JobStatusFailed, "failed"},
		{types.JobStatus("UNKNOWN_STATE"), "undetermined"},
	}

	for _, tt := range tests {
		t.Run(string(tt.state), func(t *testing.T) {
			got := batchStateToOod(tt.state)
			if got != tt.expected {
				t.Errorf("batchStateToOod(%q) = %q, want %q", tt.state, got, tt.expected)
			}
		})
	}
}
