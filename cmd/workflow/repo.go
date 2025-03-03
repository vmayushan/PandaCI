package main

import (
	"context"
	"fmt"
	"time"

	"github.com/pandaci-com/pandaci/pkg/git"
	"github.com/pandaci-com/pandaci/pkg/utils"
	pb "github.com/pandaci-com/pandaci/proto/go/v1"
	"github.com/pandaci-com/pandaci/types"
)

func cloneRepo(ctx context.Context, repo *pb.GitRepo) error {
	defer utils.MeasureTime(time.Now(), "cloning repo")
	if err := git.VerifyGit(ctx); err != nil {
		return err
	}

	cloneConfig := types.GitCloneOptions{
		FetchDepth: int(repo.FetchDepth),
		Sha:        repo.Sha,
		URL:        repo.Url,
	}

	if err := git.CloneRepo(ctx, fmt.Sprintf("/home/pandaci/repo"), cloneConfig); err != nil {
		return err
	}

	return nil
}
