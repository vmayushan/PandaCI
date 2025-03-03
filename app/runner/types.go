package runner

import (
	"context"

	pb "github.com/pandaci-com/pandaci/proto/go/v1"
	"github.com/pandaci-com/pandaci/types"
)

type (
	Handler interface {
		StartWorkflow(ctx context.Context, workflowID string, workflowReq *pb.RunnerServiceStartWorkflowRequest) (*pb.RunnerServiceStartWorkflowResponse, error)
		CleanUpJobs(ctx context.Context, workflowID string, workflowReq *pb.RunnerServiceCleanUpJobsRequest) (*pb.RunnerServiceCleanUpJobsResponse, error)
		StartJob(ctx context.Context, workflowID string, jobReq *pb.RunnerServiceStartJobRequest) (*pb.RunnerServiceStartJobResponse, error)
		StopJob(ctx context.Context, workflowID string, jobReq *pb.RunnerServiceStopJobRequest) (*pb.RunnerServiceStopJobResponse, error)

		GetRunnerType() types.RunnerType

		GetLogStream(ctx context.Context, workflowID string, req *pb.RunnerServiceGetLogStreamRequest) (*pb.RunnerServiceGetLogStreamResponse, error)
	}
)
