package main

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/alfiejones/panda-ci/pkg/utils"

	pb "github.com/alfiejones/panda-ci/proto/go/v1"
)

func checkDenoFile(ctx context.Context, config *pb.WorkflowRunnerInitConfig) error {
	defer utils.MeasureTime(time.Now(), "checking deno workflow")

	cmd := exec.CommandContext(ctx, "deno", "check", "--node-modules-dir=none", config.File)

	cmd.Dir = "/home/pandaci/repo"

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error().Err(err).Strs("out", strings.Split(string(output), "\n")).Msg("workflow file has errors")
		return fmt.Errorf("Failed to check deno file, see workflow logs for more info")
	}

	// We only use debug here as we won't want to be storing these logs
	log.Debug().Strs("out", strings.Split(string(output), "\n")).Msg("workflow logs")

	return nil
}

func verifyDeno(ctx context.Context) error {
	log.Info().Msg("verifying deno install")
	cmd := exec.CommandContext(ctx, "deno", "--version")

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error().Err(err).Strs("out", strings.Split(string(output), "\n")).Msg("failed to get deno version. Check that deno is installed")
		return err
	}

	log.Info().Strs("info", strings.Split(string(output), "\n")).Msg("deno info")

	return nil
}

func runDenoFile(ctx context.Context, config *pb.WorkflowRunnerInitConfig) error {
	defer utils.MeasureTime(time.Now(), "running deno workflow")

	// cmd := exec.CommandContext(ctx, "deno", "run", "--v8-flags=--max-old-space-size=2048,--optimize-for-size", "--allow-all", "--node-modules-dir=none", config.File)
	cmd := exec.CommandContext(ctx, "deno", "run", "--allow-all", "--node-modules-dir=none", config.File)
	cmd.Dir = "/home/pandaci/repo"

	output, err := cmd.CombinedOutput()
	if err != nil {
		// The workflow failing is probably due to bad code. Either in our library or the user
		log.Error().Err(err).Strs("out", strings.Split(string(output), "\n")).Msg("deno workflow file failed")
		return err
	}

	// We only use debug here as we won't want to be storing these logs
	log.Debug().Strs("out", strings.Split(string(output), "\n")).Msg("workflow logs")
	log.Info().Msg("workflow finished")

	return nil
}

func runWorkflow(ctx context.Context, config *pb.WorkflowRunnerInitConfig) error {
	log.Info().Msg("running workflow")
	if err := verifyDeno(ctx); err != nil {
		return err
	}

	if err := checkDenoFile(ctx, config); err != nil {
		return err
	}

	if err := runDenoFile(ctx, config); err != nil {
		return err
	}

	return nil
}
