package runnerLocal

import (
	"context"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/volume"
	nanoid "github.com/matoous/go-nanoid/v2"
)

func (h *Handler) getDenoCacheMount(ctx context.Context) (mount.Mount, error) {
	name := "panda-ci-deno-cache"

	if _, err := h.docker.VolumeCreate(ctx, volume.CreateOptions{
		Name: name,
		Labels: map[string]string{
			"panda-ci": "true",
		},
	}); err != nil {
		return mount.Mount{}, err
	}

	return mount.Mount{
		Type:     mount.TypeVolume,
		Source:   name,
		Target:   "/deno-dir",
		ReadOnly: false,
	}, nil
}

func getRootTempDir() string {
	tempDir := os.TempDir()
	return filepath.Join(tempDir, "panda-ci")
}

func getTempWorkflowDir(workflowID string) (string, error) {
	randomFolderName, err := nanoid.New(8)
	if err != nil {
		return "", err
	}

	return filepath.Join(getRootTempDir(), "workflows", workflowID, randomFolderName), nil
}

func getTempFile(ext string) (string, error) {
	randomFolderName, err := nanoid.New(8)
	if err != nil {
		return "", err
	}

	return filepath.Join(getRootTempDir(), randomFolderName+ext), nil
}
