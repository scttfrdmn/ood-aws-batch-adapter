// Package ood defines the OOD job spec types shared across adapters.
package ood

// JobSpec is the OOD job submission payload.
type JobSpec struct {
	Script      string            `json:"script"`
	JobName     string            `json:"job_name"`
	Queue       string            `json:"queue,omitempty"`
	Walltime    string            `json:"walltime,omitempty"`
	NumCores    int               `json:"num_cores,omitempty"`
	MemoryGB    int               `json:"memory_gb,omitempty"`
	Env         map[string]string `json:"env,omitempty"`
	NativeSpecs []string          `json:"native_specs,omitempty"`
}

// JobStatus maps adapter-specific states to OOD status strings.
type JobStatus struct {
	ID       string `json:"id"`
	Status   string `json:"status"` // queued, running, completed, failed, cancelled
	ExitCode int    `json:"exit_code,omitempty"`
	Message  string `json:"message,omitempty"`
}

// OOD status constants
const (
	StatusQueued    = "queued"
	StatusRunning   = "running"
	StatusCompleted = "completed"
	StatusFailed    = "failed"
	StatusCancelled = "cancelled"
	StatusUnknown   = "undetermined"
)
