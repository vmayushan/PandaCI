package service

import (
	"context"
	"fmt"
	"path"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/rand"
	"google.golang.org/protobuf/types/known/timestamppb"
	"github.com/pandaci-com/pandaci/pkg/docker"
	"github.com/pandaci-com/pandaci/pkg/utils"
	pb "github.com/pandaci-com/pandaci/proto/go/v1"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func (h *Handler) StartDockerStep(ctx context.Context, workflowID string, req *pb.JobServiceStartStepRequest, send func(*pb.JobServiceStartStepResponse) error) error {
	defer utils.MeasureTime(time.Now(), "job service start step")

	if req.GetExecData() == nil {
		return fmt.Errorf("Only supprts exec requests")
	} else if req.GetTaskMeta().GetDockerMeta() == nil {
		return fmt.Errorf("Only supports exec requests in a Docker task context")
	}

	env := []string{}
	for _, e := range req.GetExecData().GetEnv() {
		env = append(env, fmt.Sprintf("%s=%s", e.GetKey(), e.GetValue()))
	}

	resp, err := h.docker.ContainerExecCreate(ctx, req.GetTaskMeta().GetDockerMeta().GetContainerId(), container.ExecOptions{
		Cmd:          []string{"sh", "-c", req.GetExecData().GetCmd()},
		WorkingDir:   path.Join("/home/pandaci/repo", req.GetExecData().GetCwd()),
		AttachStdout: true,
		AttachStderr: true,
		Env:          env,
		Tty:          false,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create exec")
		return err
	}

	aresp, err := h.docker.ContainerExecAttach(ctx, resp.ID, container.ExecStartOptions{})
	if err != nil {
		log.Error().Err(err).Msg("Failed to attach to exec")
		return err
	}
	defer aresp.Close()

	if err := docker.StreamLogs(aresp.Reader, func(logType pb.LogMessage_ExecData_Type, data []byte) error {

		if logType == pb.LogMessage_ExecData_TYPE_STDOUT {
			return send(&pb.JobServiceStartStepResponse{
				Payload: &pb.JobServiceStartStepResponse_Exec{
					Exec: &pb.StepExecPayload{
						Timestamp: timestamppb.Now(),
						LogData: &pb.StepExecPayload_Stdout{
							Stdout: data,
						},
					},
				},
			})
		}

		return send(&pb.JobServiceStartStepResponse{
			Payload: &pb.JobServiceStartStepResponse_Exec{
				Exec: &pb.StepExecPayload{
					Timestamp: timestamppb.Now(),
					LogData: &pb.StepExecPayload_Stderr{
						Stderr: data,
					},
				},
			},
		})
	}); err != nil {
		log.Error().Err(err).Msg("Failed to stream logs")
		return err
	}

	log.Info().Msg("Exec finished")

	execRes, err := h.docker.ContainerExecInspect(ctx, resp.ID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to inspect exec")
		return err
	}

	if err := send(&pb.JobServiceStartStepResponse{
		Payload: &pb.JobServiceStartStepResponse_Exec{
			Exec: &pb.StepExecPayload{
				Timestamp: timestamppb.Now(),
				LogData: &pb.StepExecPayload_ExitCode{
					ExitCode: int32(execRes.ExitCode),
				},
			},
		},
	}); err != nil {
		log.Error().Err(err).Msg("Failed to send exec response")
		return err
	}

	return nil
}
