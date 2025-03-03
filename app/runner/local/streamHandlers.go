package runnerLocal

import (
	"context"

	pb "github.com/pandaci-com/pandaci/proto/go/v1"
)

func (h *Handler) GetLogStream(ctx context.Context, workflowID string, req *pb.RunnerServiceGetLogStreamRequest) (*pb.RunnerServiceGetLogStreamResponse, error) {

	return &pb.RunnerServiceGetLogStreamResponse{
		Url: "TODO",
	}, nil
}
