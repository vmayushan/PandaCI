package git

import (
	"context"

	typesDB "github.com/alfiejones/panda-ci/types/database"
)

func (h *Handler) UpdateRunStatusInRepo(ctx context.Context, run typesDB.WorkflowRun) error {
	project, err := h.queries.Unsafe_GetProjectByID(ctx, run.ProjectID)
	if err != nil {
		return err
	}

	org, err := h.queries.Unsafe_GetOrgByID(ctx, project.OrgID)
	if err != nil {
		return err
	}

	gitInstallation, err := h.queries.GetGitIntegration(ctx, project.GitIntegrationID)
	if err != nil {
		return err
	}

	gitClient, err := h.GetClient(gitInstallation.Type)
	if err != nil {
		return err
	}

	install, err := gitClient.NewInstallationClient(ctx, gitInstallation.ProviderID)
	if err != nil {
		return err
	}

	return install.UpdateRunInRepo(ctx, *org, *project, run)
}
