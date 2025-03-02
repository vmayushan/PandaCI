package streamer

import (
	"context"

	pb "github.com/alfiejones/panda-ci/proto/go/v1"
)

type (
	Handler interface {
		Log(ctx context.Context, logReq *pb.LogMessage) error
	}
)
