package types

import (
	pb "github.com/alfiejones/panda-ci/proto/go/v1"
)

type WorkflowFile struct {
	Hash    string
	Content string
	Path    string
}

type WorkflowRunAlertType string

const (
	WorkflowRunAlertTypeError   WorkflowRunAlertType = "error"
	WorkflowRunAlertTypeWarning WorkflowRunAlertType = "warning"
	WorkflowRunAlertTypeInfo    WorkflowRunAlertType = "info"
)

type WorkflowRunAlert struct {
	Type    WorkflowRunAlertType `json:"type" validate:"required,oneof=error warning info"`
	Title   string               `json:"title" validate:"required,max=256"`
	Message string               `json:"message" validate:"required,max=512"`
}

type RunnerWorkflowConfigDeno struct {
	Version string `json:"version"`
}

type RunnerWorkflowConfig struct {
	WorkflowMeta        pb.WorkflowMeta          `json:"workflowMeta"`
	OrchestratorAddress string                   `json:"orchestratorAddress"`
	CloneConfig         GitCloneOptions          `json:"cloneConfig"`
	File                string                   `json:"file"`
	WorkflowJWT         string                   `json:"workflowJWT"`
	DenoConfig          RunnerWorkflowConfigDeno `json:"denoConfig"`
	JWTPublicKey        string                   `json:"jwtPublicKey"`
}

type RunnerJobConfig struct {
	OrchestratorAddress string          `json:"orchestratorAddress"`
	RunnerAddress       string          `json:"runnerAddress"`
	CloneConfig         GitCloneOptions `json:"cloneConfig"`
	WorkflowJWT         string          `json:"workflowJWT"`
	JWTPublicKey        string          `json:"jwtPublicKey"`
	Host                string          `json:"host"`
	JobID               string          `json:"jobID"`
}

type WorkflowConfig struct {
	Config WorkflowRawConfig
}

type WorkflowRawConfigOnPush struct {
	Branches       *[]string `json:"branches" regex:"branches"`
	BranchesIgnore []string  `json:"branchesIgnore" regex:"branchesIgnore"`
}

type WorkflowRawConfigOnPullRequest struct {
	Events               *[]string `json:"events" regex:"events"`
	TargetBranches       *[]string `json:"targetBranches" regex:"targetBranches"`
	TargetBranchesIgnore []string  `json:"targetBranchesIgnore" regex:"targetBranchesIgnore"`
}

type WorkflowRawConfigOn struct {
	Push *WorkflowRawConfigOnPush        `json:"push" regex:"push"`
	Pr   *WorkflowRawConfigOnPullRequest `json:"pr" regex:"pr"`
}

type WorkflowRawConfig struct {
	Name string               `json:"name" regex:"name"`
	On   *WorkflowRawConfigOn `json:"on" regex:"on"`
}

type WorkflowDefintion struct {
	FileHash           string
	Committer          Committer
	GitTitle           string
	RunWorkflowRequest *pb.RunnerServiceStartWorkflowRequest
}

const (
	StepRunTypeExec = "exec"
)

type StepRunType string

type RunTrigger string

const (
	RunTriggerPush                   RunTrigger = "push"
	RunTriggerPullRequestOpened      RunTrigger = "pull_request-opened"
	RunTriggerPullRequestSynchronize RunTrigger = "pull_request-synchronize"
	RunTriggerPullRequestClosed      RunTrigger = "pull_request-closed"
	RunTriggerManual                 RunTrigger = "manual"
)

func RawPullRequestEventToRunTrigger(event string) RunTrigger {
	eventMap := map[string]RunTrigger{
		"opened":      RunTriggerPullRequestOpened,
		"synchronize": RunTriggerPullRequestSynchronize,
		"closed":      RunTriggerPullRequestClosed,
		"reopened":    RunTriggerPullRequestOpened,
	}

	return eventMap[event]
}

func RunTriggerToProto(trigger RunTrigger) pb.Trigger {
	switch trigger {
	case RunTriggerPush:
		return pb.Trigger_TRIGGER_PUSH
	case RunTriggerPullRequestOpened:
		return pb.Trigger_TRIGGER_PR_OPENED
	case RunTriggerPullRequestSynchronize:
		return pb.Trigger_TRIGGER_PR_SYNC
	case RunTriggerPullRequestClosed:
		return pb.Trigger_TRIGGER_PR_CLOSED
	case RunTriggerManual:
		return pb.Trigger_TRIGGER_MANUAL
	default:
		return pb.Trigger_TRIGGER_UNSPECIFIED
	}
}

func RunTriggerFromProto(trigger pb.Trigger) RunTrigger {
	switch trigger {
	case pb.Trigger_TRIGGER_PUSH:
		return RunTriggerPush
	case pb.Trigger_TRIGGER_PR_OPENED:
		return RunTriggerPullRequestOpened
	case pb.Trigger_TRIGGER_PR_SYNC:
		return RunTriggerPullRequestSynchronize
	case pb.Trigger_TRIGGER_PR_CLOSED:
		return RunTriggerPullRequestClosed
	case
		pb.Trigger_TRIGGER_MANUAL:
		return RunTriggerManual
	default:
		return ""
	}
}
