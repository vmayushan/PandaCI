package git

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/rs/zerolog/log"
	"github.com/pandaci-com/pandaci/types"
)

// Attempts to log the current git version
// If it fails (probably because git isn't installed), we return an error
func VerifyGit(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "git", "--version")

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error().Err(err).Bytes("out", output).Msg("failed to get git version. Check that git is installed")
		return err
	}

	log.Info().Str("version", string(output)).Msg("git info")

	return nil
}

func CloneRepo(ctx context.Context, dest string, opts types.GitCloneOptions) error {
	clone := exec.CommandContext(ctx, "git", "clone", "--no-checkout", opts.URL, dest)
	output, err := clone.CombinedOutput()
	if err != nil {
		log.Error().Err(err).Bytes("out", output).Msg("failed to clone git repo")
		return err
	}

	fetch := exec.CommandContext(ctx, "git", "fetch", "--depth", fmt.Sprintf("%d", opts.FetchDepth), "origin", opts.Sha)
	fetch.Dir = dest
	output, err = fetch.CombinedOutput()
	if err != nil {
		log.Error().Err(err).Bytes("out", output).Msg("failed to clone git repo")
		return err
	}

	checkout := exec.CommandContext(ctx, "git", "checkout", opts.Sha)
	checkout.Dir = dest
	output, err = checkout.CombinedOutput()
	if err != nil {
		log.Error().Err(err).Bytes("out", output).Msg("failed to clone git repo")
		return err
	}

	log.Info().Msg("git repo cloned")

	return nil
}
