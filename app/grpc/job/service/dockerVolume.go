package service

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/pandaci-com/pandaci/pkg/docker"
	pb "github.com/pandaci-com/pandaci/proto/go/v1"
)

func (h *Handler) CreateJobVolume(ctx context.Context, workflowID string, volumeReq *pb.JobServiceCreateJobVolumeRequest) (*pb.JobServiceCreateJobVolumeResponse, error) {

	log.Debug().Msg("Creating job volume")

	if volumeReq.Host == nil {

		name, err := docker.CreateVolume(ctx, workflowID, h.docker, "job-volume", map[string]string{
			"jobID": volumeReq.JobMeta.GetId(),
		})
		if err != nil {
			return nil, err
		}

		return &pb.JobServiceCreateJobVolumeResponse{
			Source: name,
		}, nil
	}

	// TODO - look into some validation for the host

	return &pb.JobServiceCreateJobVolumeResponse{
		Source: *volumeReq.Host,
	}, nil
}
