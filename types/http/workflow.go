package typesHTTP

import (
	"time"

	"github.com/pandaci-com/pandaci/types"
)

type Committer struct {
	Name   string `json:"name,omitempty"`
	Email  string `json:"email"`
	Avatar string `json:"avatar,omitempty"`
}

type WorkflowRun struct {
	ID     string `json:"id"`
	Number int    `json:"number"`

	CreatedAt  time.Time  `json:"createdAt"`
	FinishedAt *time.Time `json:"finishedAt,omitempty"`

	ProjectID string `json:"projectID"`

	Status     types.RunStatus      `json:"status"`
	Conclusion *types.RunConclusion `json:"conclusion,omitempty"`

	Name string `json:"name"`

	GitSha    string  `json:"gitSha"`
	GitBranch string  `json:"gitBranch"`
	GitTitle  *string `json:"gitTitle,omitempty"`

	PrNumber *int32 `json:"prNumber,omitempty"`

	PrURL     *string `json:"prURL,omitempty"`
	CommitURL *string `json:"commitURL,omitempty"`

	Trigger types.RunTrigger `json:"trigger"`

	OutputURL string `json:"outputURL,omitempty"`

	Committer Committer `json:"committer,omitempty"`

	Alerts []types.WorkflowRunAlert `json:"alerts,omitempty"`

	Jobs []JobRun `json:"jobs,omitempty"`
}

type JobRun struct {
	ID            string `json:"id"`
	Number        int    `json:"number"`
	WorkflowRunID string `json:"workflowRunId"`

	Name string `json:"name"`

	CreatedAt  time.Time  `json:"createdAt"`
	FinishedAt *time.Time `json:"finishedAt,omitempty"`

	Status     types.RunStatus      `json:"status"`
	Conclusion *types.RunConclusion `json:"conclusion,omitempty"`

	Runner string `json:"runner"`

	Tasks []TaskRun `json:"tasks"`
}

type TaskRun struct {
	ID string `json:"id"`

	CreatedAt  time.Time  `json:"createdAt"`
	FinishedAt *time.Time `json:"finishedAt,omitempty"`

	Name string `json:"name"`

	Status     types.RunStatus      `json:"status"`
	Conclusion *types.RunConclusion `json:"conclusion,omitempty"`

	DockerImage *string `json:"dockerImage,omitempty"`

	Steps []StepRun `json:"steps"`
}

type StepRun struct {
	ID string `json:"id"`

	Type types.StepRunType `json:"type"`

	CreatedAt  time.Time  `json:"createdAt"`
	FinishedAt *time.Time `json:"finishedAt,omitempty"`

	Status     types.RunStatus      `json:"status"`
	Conclusion *types.RunConclusion `json:"conclusion,omitempty"`

	OutputURL string `json:"outputURL,omitempty"`

	Name string `json:"name"`
}

func (s *StepRun) GetStep() *StepRun {
	return s
}

func (s *StepRun) GetTask() *TaskRun {
	return nil
}

type LogStream struct {
	URL           string `json:"url"`
	Authorization string `json:"authorization"`
}
