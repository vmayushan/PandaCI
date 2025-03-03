package streamer

import (
	"context"

	pb "github.com/pandaci-com/pandaci/proto/go/v1"
)

type (
	Handler interface {
		Log(ctx context.Context, logReq *pb.LogMessage) error
	}
)
