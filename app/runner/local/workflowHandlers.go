package runnerLocal

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/pandaci-com/pandaci/pkg/docker"
	"github.com/pandaci-com/pandaci/pkg/utils"
	"github.com/pandaci-com/pandaci/pkg/utils/env"
	pb "github.com/pandaci-com/pandaci/proto/go/v1"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handler) StartWorkflow(ctx context.Context, workflowID string, workflowReq *pb.RunnerServiceStartWorkflowRequest) (*pb.RunnerServiceStartWorkflowResponse, error) {
	defer utils.MeasureTime(time.Now(), "StartWorkflow local runner")

	image := "4f9f1458a17253ba9584db3d8d06cdde56b038d248a49fcaa3181396d296659c" // TODO - figure out a better way to get the image. Ideally we'd build it at some point

	denoCache, err := h.getDenoCacheMount(ctx)
	if err != nil {
		return nil, err
	}

	if !utils.IsValidDenoWorkflow(workflowReq.GetFilePath()) {
		return nil, fmt.Errorf("unsupported workflow file type")
	}

	volumeMounts := []mount.Mount{
		denoCache,
	}

	containerName, err := docker.GenerateSafeDockerName(workflowID, workflowReq.GetFilePath())
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

	config := &pb.WorkflowRunnerInitConfig{
		OrchestratorAddress: *orchestratorAddress,
		PresignedOutputUrl:  workflowReq.PresignedOutputUrl,
		WorkflowMeta: &pb.WorkflowMeta{
			WorkflowJwt: workflowReq.WorkflowJwt,
			Name:        workflowReq.Name,
			Runnner:     "ubuntu-2x",
			StartedAt:   timestamppb.New(time.Now()),
			TimeoutAt:   timestamppb.New(time.Now().Add(time.Hour)),
			Repo:        workflowReq.GitInfo,
			RunnerMeta: &pb.WorkflowMeta_DockerMeta{
				DockerMeta: &pb.DockerRunMeta{
					CreatedAt: timestamppb.New(time.Now()),
				},
			},
		},
		File:         workflowReq.FilePath,
		PublicJwtKey: *publicKey,
		Host:         "localhost",
	}

	configBytes, err := proto.Marshal(config)
	if err != nil {
		return nil, err
	}

	resp, err := h.docker.ContainerCreate(ctx, &container.Config{
		Image: image,
		Cmd:   []string{base64.StdEncoding.EncodeToString(configBytes)},
		Labels: map[string]string{
			"panda-ci":   "true",
			"workflowID": workflowID,
		},
		AttachStdout: true,
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
				log.Info().Bytes("data", data).Msg("Workflow log")
				return nil
			}); err != nil {
				log.Error().Err(err).Str("containerID", resp.ID).Msg("failed to stream logs 1")
			}
		}
	}()

	return &pb.RunnerServiceStartWorkflowResponse{
		RunnerMeta: &pb.RunnerServiceStartWorkflowResponse_DockerMeta{
			DockerMeta: &pb.DockerRunMeta{
				CreatedAt: timestamppb.New(time.Now()),
			},
		},
	}, nil
}

func (h *Handler) StopWorkflow(ctx context.Context, workflowID string, workflowReq *pb.RunnerServiceStopWorkflowRequest) (*pb.RunnerServiceStopWorkflowResponse, error) {
	defer utils.MeasureTime(time.Now(), "StopWorkflow local runner")

	if err := docker.DeleteContainers(ctx, h.docker, docker.DeleteContainerFilter{
		WorkflowID: &workflowID,
	}); err != nil {
		// We don't want to return an error here as we still want to delete the volumes
		log.Error().Err(err).Msg("failed to delete job containers")
	}

	if err := docker.DeleteVolumes(ctx, h.docker, docker.DeleteVolumeFilter{
		WorkflowID: &workflowID,
	}); err != nil {
		// no point in returning an error here since it's just cleaning up
		log.Error().Err(err).Msg("failed to delete job volumes")
	}

	return &pb.RunnerServiceStopWorkflowResponse{}, nil
}
