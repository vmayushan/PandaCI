package runnerLocal

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"github.com/pandaci-com/pandaci/pkg/docker"
	"github.com/pandaci-com/pandaci/pkg/utils"
	"github.com/pandaci-com/pandaci/pkg/utils/env"
	pb "github.com/pandaci-com/pandaci/proto/go/v1"
)

func (h *Handler) CleanUpJobs(ctx context.Context, workflowID string, workflowReq *pb.RunnerServiceCleanUpJobsRequest) (*pb.RunnerServiceCleanUpJobsResponse, error) {
	return &pb.RunnerServiceCleanUpJobsResponse{}, nil
}

func (h *Handler) StartJob(ctx context.Context, workflowID string, req *pb.RunnerServiceStartJobRequest) (*pb.RunnerServiceStartJobResponse, error) {
	defer utils.MeasureTime(time.Now(), "StartJob local runner")

	image := "ad2a8432e80c13b7d2fe28da36e6d726af0be25fa4694e087d052c85d10359b4" // TODO - figure out a better way to get the image. Ideally we'd build it at some point

	// TODO - this won't work on windows
	volumeMounts := []mount.Mount{
		{
			Type:   mount.TypeBind,
			Source: "/var/run/docker.sock",
			Target: "/var/run/docker.sock",
		},
	}

	containerName, err := docker.GenerateSafeDockerName(workflowID, req.Name)
	if err != nil {
		return nil, err
	}

	log.Debug().Msg("Container about to be created")

	publicKey, err := env.GetRunnerPublicKey()
	if err != nil {
		return nil, err
	}

	orchestratorAddress, err := env.GetOrchestratorGRPCURL()
	if err != nil {
		return nil, err
	}

	config := &pb.JobRunnerInitConfig{
		WorkflowMeta: req.WorkflowMeta,
		JobMeta: &pb.JobMeta{
			Id:   req.Id,
			Name: req.Name,
			RunnerMeta: &pb.JobMeta_DockerMeta{
				DockerMeta: &pb.DockerRunMeta{
					CreatedAt: timestamppb.New(time.Now()),
				},
			},
			Address: "TODO", // We fill this in after confirming the container is running
			Runner:  req.Runner,
		},
		OrchestratorAddress: *orchestratorAddress,
		RunnerAddress:       *orchestratorAddress, // TODO - this should be the runner address.
		PublicJwtKey:        *publicKey,
		Host:                "localhost",
	}

	configBytes, err := proto.Marshal(config)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal config")
		return nil, err
	}

	resp, err := h.docker.ContainerCreate(ctx, &container.Config{
		Image: image,
		Cmd:   []string{base64.StdEncoding.EncodeToString(configBytes)},
		Labels: map[string]string{
			"panda-ci":   "true",
			"workflowID": workflowID,
			"jobID":      req.Id,
		},
		Env: []string{
			"PANDACI_LOCAL=true",
		},
		AttachStdout: true, // TODO - figure out if we actually need this, same for workflow
		AttachStderr: true,
	}, &container.HostConfig{
		NetworkMode: "host",
		Mounts:      volumeMounts,
	}, nil, nil, containerName)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create container")
		return nil, err
	}

	if err := h.docker.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		log.Error().Err(err).Msg("Failed to start container")
		return nil, err
	}

	log.Info().Str("containerID", resp.ID).Msg("container started")

	go func() {
		reader, err := h.docker.ContainerLogs(context.Background(), resp.ID, container.LogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Since:      "2000-11-11",
			Follow:     true,
		})
		if err != nil {
			log.Error().Err(err).Str("containerID", resp.ID).Msg("failed to get container logs")
		} else {

			defer reader.Close()

			if err := docker.StreamLogs(reader, func(logType pb.LogMessage_ExecData_Type, data []byte) error {
				log.Info().Bytes("data", data).Msg("job log")
				return nil
			}); err != nil {
				log.Error().Err(err).Str("containerID", resp.ID).Msg("failed to stream logs 1")
			}
		}
	}()

	return &pb.RunnerServiceStartJobResponse{
		JobMeta: config.JobMeta,
	}, nil
}

func (h *Handler) StopJob(ctx context.Context, workflowID string, req *pb.RunnerServiceStopJobRequest) (*pb.RunnerServiceStopJobResponse, error) {
	defer utils.MeasureTime(time.Now(), "StopJob local runner")

	if err := docker.DeleteContainers(ctx, h.docker, docker.DeleteContainerFilter{
		WorkflowID: &workflowID,
		JobID:      &req.JobMeta.Id,
	}); err != nil {
		// We don't want to return an error here as we still want to delete the volumes
		log.Error().Err(err).Msg("failed to delete job containers")
	}

	if err := docker.DeleteVolumes(ctx, h.docker, docker.DeleteVolumeFilter{
		WorkflowID: &workflowID,
		JobID:      &req.JobMeta.Id,
	}); err != nil {
		// no point in returning an error here since it's just cleaning up
		log.Error().Err(err).Msg("failed to delete job volumes")
	}

	return &pb.RunnerServiceStopJobResponse{}, nil
}
