package types

import (
	"time"
)

type JobType string

type RunnerType string

type Job struct {
	ID         string `json:"id"`
	WorkflowID string `json:"workflowID"`
	Name       string `json:"name"`

	RunnerType RunnerType `json:"runnerType"`

	RunnerData map[string]any `json:"runnerData"`

	StartedAt  time.Time `json:"startedAt"`
	FinishedAt time.Time `json:"finishedAt"`

	GrpcAddress string `json:"-"`
}

type JobRunnerDataLocal struct {
	Image       string `json:"image"`
	ContainerID string `json:"containerID"`
}

const (
	RunnerTypeLocal RunnerType = "local"
	RunnerTypeFly   RunnerType = "fly"
)
