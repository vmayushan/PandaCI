package main

import (
	"os"

	pb "github.com/alfiejones/panda-ci/proto/go/v1"
)

func setEnv(workflowID string, config *pb.JobRunnerInitConfig, port int) error {
	envs := [][2]string{
		{"CI", "true"},
	}

	for _, env := range envs {
		os.Setenv(env[0], env[1])
	}

	return nil
}
