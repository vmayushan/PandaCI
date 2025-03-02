package runnerLocal

import (
	"context"

	pb "github.com/alfiejones/panda-ci/proto/go/v1"
)

func (h *Handler) GetLogStream(ctx context.Context, workflowID string, req *pb.RunnerServiceGetLogStreamRequest) (*pb.RunnerServiceGetLogStreamResponse, error) {

	return &pb.RunnerServiceGetLogStreamResponse{
		Url: "TODO",
	}, nil
}
