package orchestrator

import (
	"context"
	"fmt"
	"os"

	"connectrpc.com/connect"
	queriesWorkflow "github.com/pandaci-com/pandaci/app/queries/workflow"
	"github.com/pandaci-com/pandaci/pkg/jwt"
	"github.com/pandaci-com/pandaci/platform/storage"
	"github.com/pandaci-com/pandaci/types"
	typesDB "github.com/pandaci-com/pandaci/types/database"
	"github.com/rs/zerolog/log"
)

func (h *Handler) hasMinutesLeft(ctx context.Context, org *typesDB.OrgDB) (bool, error) {
	license, err := org.GetLicense()
	if err != nil {
		return false, err
	}

	if license.Plan == types.CloudSubscriptionPlanPaused {
		return false, nil
	}

	billingPeriod := license.GetBillingPeriod()

	if license.Plan == types.CloudSubscriptionPlanFree {
		buildMinutes, err := h.queries.GetBuildMinutes(ctx, org.ID, billingPeriod.StartsAt, billingPeriod.EndsAt)
		if err != nil {
			return false, err
		}

		return buildMinutes < license.Features.MaxBuildMinutes, nil
	}

	return true, nil
}

func (h *Handler) hasCommittersLeft(ctx context.Context, org *typesDB.OrgDB, workflowDefs []types.WorkflowDefintion) (bool, error) {
	license, err := org.GetLicense()
	if err != nil {
		return false, err
	}

	if license.Plan == types.CloudSubscriptionPlanPaused {
		return false, nil
	}

	billingPeriod := license.GetBillingPeriod()

	if license.Plan == types.CloudSubscriptionPlanFree {

		newCommitters := make([]queriesWorkflow.Committer, len(workflowDefs))

		for i, def := range workflowDefs {
			newCommitters[i] = queriesWorkflow.Committer{
				UserID:         def.Committer.UserID,
				CommitterEmail: def.Committer.Email,
			}
		}

		committersCount, err := h.queries.CountCommitters(ctx, org.ID, billingPeriod.StartsAt, billingPeriod.EndsAt, &newCommitters)
		if err != nil {
			return false, err
		}

		return committersCount <= license.Features.MaxCommitters, nil
	}

	return true, nil
}

func (h *Handler) StartWorkflows(ctx context.Context, project *typesDB.Project, workflowDefs []types.WorkflowDefintion) ([]typesDB.WorkflowRun, error) {

	org, err := h.queries.Unsafe_GetOrgByID(ctx, project.OrgID)
	if err != nil {
		return nil, err
	}

	hasMinutesLeft, err := h.hasMinutesLeft(ctx, org)
	if err != nil {
		log.Err(err).Msg("Failed to check if we have minutes left")
		return nil, err
	}

	hasCommittersLeft, err := h.hasCommittersLeft(ctx, org, workflowDefs)
	if err != nil {
		log.Err(err).Msg("Failed to check if we have committers left")
		return nil, err
	}

	bucketClient, err := storage.GetClient(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to create s3 client")
		return nil, err
	}

	var runs []typesDB.WorkflowRun
	for _, def := range workflowDefs {

		if !hasMinutesLeft {
			log.Debug().Msg("No minutes left, skipping workflow")

			// add an error to the workflow run
			workflow := typesDB.WorkflowRun{
				ProjectID:      project.ID,
				Name:           def.RunWorkflowRequest.Name,
				CommitterEmail: def.Committer.Email,
				UserID:         def.Committer.UserID,
				GitTitle:       &def.GitTitle,
				GitSha:         def.RunWorkflowRequest.GitInfo.Sha,
				GitBranch:      def.RunWorkflowRequest.GitInfo.Branch,
				Trigger:        types.RunTriggerFromProto(def.RunWorkflowRequest.GetTrigger()),
				PrNumber:       def.RunWorkflowRequest.PrNumber,
				Runner:         "ubuntu-2x", // TODO - we need an unknown runner maybe
			}

			workflow.AppendAlert(types.WorkflowRunAlert{
				Title:   "Out of build minutes",
				Message: "Please upgrade your plan for more minutes",
				Type:    types.WorkflowRunAlertTypeError,
			})

			if err := h.queries.CreateFailedWorkflowRun(ctx, &workflow); err != nil {
				log.Err(err).Msg("Failed to create failed workflow run")
			}

			continue
		} else if !hasCommittersLeft {
			// add an error to the workflow run
			workflow := typesDB.WorkflowRun{
				ProjectID:      project.ID,
				Name:           def.RunWorkflowRequest.Name,
				CommitterEmail: def.Committer.Email,
				UserID:         def.Committer.UserID,
				GitTitle:       &def.GitTitle,
				GitSha:         def.RunWorkflowRequest.GitInfo.Sha,
				GitBranch:      def.RunWorkflowRequest.GitInfo.Branch,
				Trigger:        types.RunTriggerFromProto(def.RunWorkflowRequest.GetTrigger()),
				PrNumber:       def.RunWorkflowRequest.PrNumber,
				Runner:         "ubuntu-2x", // TODO - we need an unknown runner maybe
			}

			workflow.AppendAlert(types.WorkflowRunAlert{
				Title:   "Out of committers",
				Message: "Please upgrade your plan for more committers",
				Type:    types.WorkflowRunAlertTypeError,
			})

			if err := h.queries.CreateFailedWorkflowRun(ctx, &workflow); err != nil {
				log.Err(err).Msg("Failed to create failed workflow run")
			}

			continue
		}

		if run, err := h.StartWorkflow(ctx, project, def, bucketClient); err != nil {
			log.Err(err).Msg("Failed to start workflow")
			// we still want to start the other workflows
		} else {
			runs = append(runs, *run)
		}
	}

	return runs, nil
}

func (h *Handler) StartWorkflow(ctx context.Context, project *typesDB.Project, workflowDef types.WorkflowDefintion, bucketClient *storage.BucketClient) (*typesDB.WorkflowRun, error) {
	// TOOD - handle errors here and attempt to set the workflow run to failed

	workflowRun := typesDB.WorkflowRun{
		ProjectID:      project.ID,
		Name:           workflowDef.RunWorkflowRequest.Name,
		Status:         types.RunStatusPending,
		GitSha:         workflowDef.RunWorkflowRequest.GitInfo.Sha,
		GitBranch:      workflowDef.RunWorkflowRequest.GitInfo.Branch,
		Runner:         "ubuntu-2x",
		Trigger:        types.RunTriggerFromProto(workflowDef.RunWorkflowRequest.GetTrigger()),
		PrNumber:       workflowDef.RunWorkflowRequest.PrNumber,
		CommitterEmail: workflowDef.Committer.Email,
		UserID:         workflowDef.Committer.UserID,
		GitTitle:       &workflowDef.GitTitle,
	}

	if err := h.queries.CreateWorkflowRun(ctx, &workflowRun); err != nil {
		return nil, err
	}

	if err := h.git.UpdateRunStatusInRepo(ctx, workflowRun); err != nil {
		log.Err(err).Msg("Failed to update run status in repo")
		// we still want to start the workflow
		// hopefully, the next time we try to update the status it will work
	}

	workflowJWT, err := h.jwt.CreateWorkflowToken(jwt.WorkflowClaims{
		WorkflowID: workflowRun.ID,
		ProjectID:  project.ID,
		OrgID:      project.OrgID,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create workflow jwt")
		h.queries.FailWorkflowRun(ctx, &workflowRun,
			types.WorkflowRunAlert{
				Type:    types.WorkflowRunAlertTypeError,
				Title:   "Failed to start workflow",
				Message: "Failed to create workflow jwt, please try again",
			})

		return nil, err
	}

	presignedWorkflowLogsURL, err := bucketClient.PutObject(ctx, os.Getenv("WORKFLOW_LOGS_BUCKET"), fmt.Sprintf("%s/%s/%s/workflow.csv", project.OrgID, project.ID, workflowRun.ID), "text/csv", 60*60*24*7)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create presigned workflow logs url")

		h.queries.FailWorkflowRun(ctx, &workflowRun, types.WorkflowRunAlert{
			Type:    types.WorkflowRunAlertTypeError,
			Title:   "Failed to start workflow",
			Message: "Failed to create presigned workflow logs url, please try again",
		})

		return nil, err
	}

	workflowDef.RunWorkflowRequest.WorkflowJwt = workflowJWT
	workflowDef.RunWorkflowRequest.PresignedOutputUrl = presignedWorkflowLogsURL.URL

	log.Debug().Str("url", presignedWorkflowLogsURL.URL).Msg("Presigned workflow logs url")

	log.Debug().Msgf("Starting workflow with id: %s", workflowRun.ID)

	go func() {
		if _, err = h.runnerClient.StartWorkflow(context.Background(), connect.NewRequest(workflowDef.RunWorkflowRequest)); err != nil {
			log.Error().Err(err).Msg("Failed to start workflow")

			if err := h.queries.FailWorkflowRun(context.Background(), &workflowRun, types.WorkflowRunAlert{
				Type:    types.WorkflowRunAlertTypeError,
				Title:   "Failed to start workflow",
				Message: "Failed to start workflow with runner, please try again",
			}); err != nil {
				log.Error().Err(err).Msg("Failed to fail workflow run")
			}
		}
	}()

	return &workflowRun, nil
}
