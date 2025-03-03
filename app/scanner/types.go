package scanner

import (
	"context"

	"github.com/pandaci-com/pandaci/types"
	typesDB "github.com/pandaci-com/pandaci/types/database"
)

type WorkflowFile struct {
	RelativePath string
}

type GithubRepoInfo struct{}

type LocalRepoInfo struct {
	AbsolutePath string
	FileGlob     string
}

type (
	Handler interface {
		GetWorkflowDefinitions(ctx context.Context, project typesDB.Project, triggerEvent types.TriggerEvent) ([]types.WorkflowDefintion, error)
	}
)
