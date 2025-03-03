package service

import (
	"context"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/rs/zerolog/log"
	"github.com/pandaci-com/pandaci/pkg/docker"
	"github.com/pandaci-com/pandaci/pkg/utils"
	pb "github.com/pandaci-com/pandaci/proto/go/v1"
)

func (h *Handler) StarDockerTask(ctx context.Context, workflowID string, req *pb.JobServiceStartTaskRequest) (*pb.JobServiceStartTaskResponse, error) {
	defer utils.MeasureTime(time.Now(), "job service start task")

	if req.Data.GetDockerData() == nil {
		// TODO - support native tasks
		return nil, fmt.Errorf("Only supprts docker tasks")
	}

	image := req.Data.GetDockerData().Image

	if err := docker.PullImage(ctx, h.docker, image); err != nil {
		log.Error().Err(err).Str("image", image).Msg("Failed to pull image")
		return nil, err
	}

	containerName, err := docker.GenerateSafeDockerName(workflowID, req.Data.Name)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate container name")
		return nil, err
	}

	repoMount := mount.Mount{
		Type:   mount.TypeBind,
		Source: "/home/pandaci/repo",
		Target: "/home/pandaci/repo",
	}

	volumeMounts := []mount.Mount{
		repoMount,
	}

	for _, taskVolume := range req.Data.GetDockerData().Volumes {
		volMount := mount.Mount{
			Type:   mount.TypeVolume,
			Source: taskVolume.Source,
			Target: path.Join("/home/pandaci/repo", taskVolume.Target),
		}
		if taskVolume.Type == pb.StartTaskData_Docker_DockerVolume_TYPE_BIND {
			volMount.Type = mount.TypeBind
			// create the folders if it doesn't already exist
			if err := os.MkdirAll(volMount.Source, os.ModePerm); err != nil {
				log.Error().Err(err).Str("source", volMount.Source).Msg("Failed to create directory")
				return nil, fmt.Errorf("failed to create directory %s: %w", volMount.Source, err)
			}
		}
		volumeMounts = append(volumeMounts, volMount)
	}

	resp, err := h.docker.ContainerCreate(ctx, &container.Config{
		Image:      req.Data.GetDockerData().Image,
		Entrypoint: []string{"sleep", "infinity"}, // TODO -  We should probably allow users to specify which cmd they want. We can probably default to bash
		Tty:        true,
		OpenStdin:  true,
		Env: []string{
			"CI=true",
		},
		Labels: map[string]string{
			"panda-ci":   "true",
			"workflowID": workflowID,
			"taskID":     req.Id,
			"jobID":      req.JobMeta.Id,
		},
	}, &container.HostConfig{
		NetworkMode: "host",
		Mounts:      volumeMounts,
		// Resources: container.Resources{
		// 	Memory:     4 * 1024 * 1024 * 1024,  // 4 GB in bytes
		// 	MemorySwap: 12 * 1024 * 1024 * 1024, // 12 GB in bytes
		// },
		//
	}, nil, nil, containerName)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create task container")
		return nil, err
	}

	if err := h.docker.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		log.Error().Err(err).Msg("Failed to start task container")
		return nil, err
	}

	return &pb.JobServiceStartTaskResponse{
		TaskMeta: &pb.TaskMeta{
			Id:   req.Id,
			Name: req.Data.Name,
			SpecificMeta: &pb.TaskMeta_DockerMeta{
				DockerMeta: &pb.TaskMeta_Docker{
					ContainerId: resp.ID,
				},
			},
		},
	}, nil
}

func (h *Handler) StopDockerTask(ctx context.Context, workflowID string, req *pb.JobServiceStopTaskRequest) (*pb.JobServiceStopTaskResponse, error) {

	if err := docker.DeleteContainers(ctx, h.docker, docker.DeleteContainerFilter{
		WorkflowID: &workflowID,
		JobID:      &req.JobMeta.Id,
		TaskID:     &req.TaskMeta.Id,
	}); err != nil {
		log.Error().Err(err).Msg("failed to delete workflow containers")
		return nil, err
	}

	return &pb.JobServiceStopTaskResponse{}, nil
}
