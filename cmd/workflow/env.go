package main

import (
	"fmt"
	"os"

	pb "github.com/alfiejones/panda-ci/proto/go/v1"
)

func setEnv(workflowID string, config *pb.WorkflowRunnerInitConfig, port int) error {
	envs := [][2]string{
		{"NO_COLOR", "true"},
		{"CI", "true"},
		{"PANDACI_WORKFLOW_ID", workflowID},
		{"PANDACI_WORKFLOW_JWT", config.WorkflowMeta.WorkflowJwt},
		{"PANDACI_WORKFLOW_GRPC_ADDRESS", fmt.Sprintf("http://localhost:%d", port)},
		{"PANDACI_COMMIT_SHA", config.WorkflowMeta.Repo.Sha},
		{"PANDACI_BRANCH", config.WorkflowMeta.Repo.Branch},
		{"PANDACI_REPO_URL", config.WorkflowMeta.Repo.Url}, // TODO - change this to GIT_URL and the repo url should be the html url
	}

	if config.WorkflowMeta.PrNumber != nil {
		envs = append(envs, [2]string{"PANDACI_PR_NUMBER", fmt.Sprintf("%d", config.WorkflowMeta.GetPrNumber())})
	}

	for _, env := range envs {
		os.Setenv(env[0], env[1])
	}

	for _, env := range config.Env {
		os.Setenv(env.Key, env.Value)
	}

	return nil
}
