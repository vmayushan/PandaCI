package scannerGit

import (
	"context"
	"fmt"
	"slices"
	"time"

	scannerShared "github.com/alfiejones/panda-ci/app/scanner/shared"
	"github.com/alfiejones/panda-ci/pkg/encryption"
	"github.com/alfiejones/panda-ci/pkg/utils/env"
	"github.com/alfiejones/panda-ci/types"
	typesDB "github.com/alfiejones/panda-ci/types/database"
	"github.com/rs/zerolog/log"

	"github.com/gobwas/glob"

	pb "github.com/alfiejones/panda-ci/proto/go/v1"
)

func shouldRun(ctx context.Context, config types.WorkflowRawConfig, triggerEvent types.TriggerEvent) (bool, error) {

	log.Info().Interface("config", config).Interface("triggerEvent", triggerEvent).Msg("Checking if we should run")

	if triggerEvent.Trigger == types.RunTriggerManual {
		return true, nil
	}

	if triggerEvent.Trigger == types.RunTriggerPush {

		if config.On == nil {
			config.On = &types.WorkflowRawConfigOn{}
		}

		if config.On.Push == nil {
			config.On.Push = &types.WorkflowRawConfigOnPush{}
		}

		if config.On.Push.Branches == nil {
			config.On.Push.Branches = &[]string{"*"}
		}

		for _, ignore := range config.On.Push.BranchesIgnore {
			g, err := glob.Compile(ignore)
			if err != nil {
				// TODO - we should add a warning to the workflow run
				log.Error().Err(err).Msg("Failed to compile glob")
				continue
			}

			if g.Match(triggerEvent.Branch) {
				return false, nil
			}
		}

		for _, branch := range *config.On.Push.Branches {
			g, err := glob.Compile(branch)
			if err != nil {
				// TODO - we should add a warning to the workflow run
				log.Error().Err(err).Msg("Failed to compile glob")
				continue
			}

			if g.Match(triggerEvent.Branch) {
				return true, nil
			}
		}

		return false, nil
	}

	if slices.Contains([]types.RunTrigger{types.RunTriggerPullRequestOpened, types.RunTriggerPullRequestClosed, types.RunTriggerPullRequestSynchronize}, triggerEvent.Trigger) {

		if config.On == nil || config.On.Pr == nil {
			return false, nil
		}

		if triggerEvent.TargetBranch == nil {
			return false, fmt.Errorf("Target branch is nil for a PR")
		}

		if config.On.Pr.Events == nil {
			config.On.Pr.Events = &[]string{}
		}

		hasEvent := false
		for _, event := range *config.On.Pr.Events {
			runTrigger := types.RawPullRequestEventToRunTrigger(event)
			if runTrigger == triggerEvent.Trigger {
				hasEvent = true
				break
			}
		}

		if !hasEvent {
			return false, nil
		}

		if config.On.Pr.TargetBranchesIgnore == nil {
			config.On.Pr.TargetBranchesIgnore = []string{}
		}

		for _, ignore := range config.On.Pr.TargetBranchesIgnore {
			g, err := glob.Compile(ignore)
			if err != nil {
				// TODO - we should add a warning to the workflow run
				log.Error().Err(err).Msg("Failed to compile glob")
				continue
			}

			if g.Match(*triggerEvent.TargetBranch) {
				return false, nil
			}
		}

		if config.On.Pr.TargetBranches == nil {
			config.On.Pr.TargetBranches = &[]string{"*"}
		}

		for _, branch := range *config.On.Pr.TargetBranches {
			g, err := glob.Compile(branch)
			if err != nil {
				// TODO - we should add a warning to the workflow run
				log.Error().Err(err).Msg("Failed to compile glob")
				continue
			}

			if g.Match(*triggerEvent.TargetBranch) {
				return true, nil
			}
		}

		return false, nil
	}

	return false, nil
}

func (h *Handler) GetWorkflowDefinitions(ctx context.Context, project typesDB.Project, triggerEvent types.TriggerEvent) ([]types.WorkflowDefintion, error) {
	gitIntegration, err := h.queries.GetGitIntegration(ctx, project.GitIntegrationID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get git integration")
		return nil, err
	}

	githubClient, err := h.gitHandler.GetClient(gitIntegration.Type)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get github client")
		return nil, err
	}

	installClient, err := githubClient.NewInstallationClient(ctx, gitIntegration.ProviderID)
	if err != nil {
		return nil, err
	}

	gitRepoData, err := installClient.GetProjectGitRepoData(ctx, project)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get project git repo data")
		return nil, err
	}

	files, err := installClient.GetWorkflowFiles(ctx, project, triggerEvent)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get workflow files")
		return nil, err
	}

	environmentVars, err := h.getDecryptedProjectEnvs(ctx, project, triggerEvent.Branch)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get decrypted project envs")
		return nil, err
	}

	var workflows []types.WorkflowDefintion

	for _, file := range files {

		// used if we come across an error
		tempWorkflow := typesDB.WorkflowRun{
			GitBranch:      triggerEvent.Branch,
			GitTitle:       &triggerEvent.GitTitle,
			ProjectID:      project.ID,
			Name:           file.Path,
			GitSha:         triggerEvent.SHA,
			CommitterEmail: triggerEvent.Committer.Email,
			PrNumber:       triggerEvent.PrNumber,
			UserID:         triggerEvent.Committer.UserID,
			Trigger:        triggerEvent.Trigger,
		}

		config, err := scannerShared.ExtractWorkflowConfig([]byte(file.Content))
		if err != nil {
			log.Error().Err(err).Msg("Failed to extract workflow config")
			if err := tempWorkflow.AppendAlert(types.WorkflowRunAlert{
				Type:    types.WorkflowRunAlertTypeError,
				Title:   "Failed to parse workflow config",
				Message: err.Error(),
			}); err != nil {
				log.Error().Err(err).Msg("Failed to append alert")
			}

			if err := h.queries.CreateFailedWorkflowRun(ctx, &tempWorkflow); err != nil {
				log.Error().Err(err).Msg("Failed to create failed workflow run")
			}

			return nil, err
		}

		if shouldRun, err := shouldRun(ctx, config.Config, triggerEvent); err != nil {
			log.Error().Err(err).Msg("Failed to check if we should run")

			if err := tempWorkflow.AppendAlert(types.WorkflowRunAlert{
				Type:    types.WorkflowRunAlertTypeError,
				Title:   "Failed to check if we should run",
				Message: err.Error(),
			}); err != nil {
				log.Error().Err(err).Msg("Failed to append alert")
			}

			if err := h.queries.CreateFailedWorkflowRun(ctx, &tempWorkflow); err != nil {
				log.Error().Err(err).Msg("Failed to create failed workflow run")
			}

			return nil, err
		} else if !shouldRun {
			log.Info().Msg("Skipping workflow")
			continue
		}

		if config.Config.Name == "" {
			config.Config.Name = file.Path
		}

		workflows = append(workflows, types.WorkflowDefintion{
			FileHash:  file.Hash,
			Committer: triggerEvent.Committer,
			GitTitle:  triggerEvent.GitTitle,
			RunWorkflowRequest: &pb.RunnerServiceStartWorkflowRequest{
				Name: config.Config.Name,
				// TimeoutAt: timestamppb.Now(), // TODO - get from config
				Env:      environmentVars,
				Image:    fmt.Sprintf("denoland/deno:alpine-%s", "2.1.10"),
				FilePath: file.Path,
				Trigger:  types.RunTriggerToProto(triggerEvent.Trigger),
				PrNumber: triggerEvent.PrNumber,
				GitInfo: &pb.GitRepo{
					Url:        gitRepoData.URL,
					Sha:        triggerEvent.SHA,
					Branch:     triggerEvent.Branch,
					FetchDepth: 1,
				},
				LanguageConfig: &pb.RunnerServiceStartWorkflowRequest_DenoConfig_{
					DenoConfig: &pb.RunnerServiceStartWorkflowRequest_DenoConfig{
						Version: "2.1.7",
					},
				},
			},
		})
	}

	return workflows, nil
}

func (h *Handler) getDecryptedProjectEnvs(ctx context.Context, project typesDB.Project, branch string) ([]*pb.EnvironmentVariable, error) {

	projectVariables, err := h.queries.GetProjectVariablesWithEnvironments(ctx, &project)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get project variables")
		return nil, err
	}

	var matchedVariables []*pb.EnvironmentVariable

	for _, variable := range projectVariables {

		match := len(variable.Environments) == 0
		for _, environment := range variable.Environments {
			if matched, err := matchWithTimeout(environment.BranchPattern, branch); err != nil {
				log.Error().Err(err).Msg("Failed to match branch pattern")
				return nil, err
			} else if matched {
				match = true
				break
			}
		}

		if match {
			index := len(matchedVariables)

			// check if variable key already exists, if so delete it as we want to overwrite it
			for i, matchedVariable := range matchedVariables {
				if matchedVariable.Key == variable.Key {
					index = i
				}
			}

			key, err := env.GetEncryptionKey(variable.EncryptionKeyID)
			if err != nil {
				log.Error().Err(err).Msg("Failed to get encryption key")
				return nil, err
			}

			decrpytedValue, err := encryption.Decrypt(variable.Value, variable.InitialisationVector, *key)
			if err != nil {
				log.Error().Err(err).Msg("Failed to decrypt variable")
				return nil, err
			}

			matchedVariables = append(matchedVariables[:index], append([]*pb.EnvironmentVariable{{
				Key:   variable.Key,
				Value: decrpytedValue,
			}}, matchedVariables[index:]...)...)

		}
	}

	log.Debug().Interface("matchedVariables", matchedVariables).Msg("Matched Variables")

	return matchedVariables, nil
}

// since we allow users to define branch patterns for environments,
// we want a timeout to avoid regex DoS attacks
func matchWithTimeout(pattern, input string) (bool, error) {
	ctx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)

	result := make(chan bool, 1)
	errChan := make(chan error, 1)

	go func() {
		g, err := glob.Compile(pattern)
		if err != nil {
			errChan <- err
			return
		}
		result <- g.Match(input)
	}()

	select {
	case <-ctx.Done():
		return false, ctx.Err() // Timeout or cancellation
	case err := <-errChan:
		log.Error().Err(err).Msg("Failed to compile regex")
		return false, err // Regex compilation error
	case res := <-result:
		return res, nil // Match result
	}

}
