package grpcWorkflow

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"github.com/alfiejones/panda-ci/pkg/flyClient"
	pb "github.com/alfiejones/panda-ci/proto/go/v1"
	pbConnect "github.com/alfiejones/panda-ci/proto/go/v1/v1connect"
	"github.com/rs/zerolog/log"
)

// Pings the job service every 5 seconds to check its still alive
// We need 3 failed pings to consider the workflow dead
// If dead, we'll fail the job
func (h *Handler) monitorJobStatus(jobMeta *pb.JobMeta, workflowMeta *pb.WorkflowMeta) {

	ctx, cancel := context.WithCancel(context.Background())

	h.jobMonitorContextMap[jobMeta.Id] = &cancel

	defer func() {
		cancel()
		delete(h.jobMonitorContextMap, jobMeta.Id)
	}()

	jobClient := pbConnect.NewJobServiceClient(&http.Client{
		Transport: &flyClient.FlyRoundTripper{
			Base:    http.DefaultTransport,
			AppName: &jobMeta.GetFlyMeta().AppName,
			Headers: map[string]string{
				"Authorization": "Bearer " + workflowMeta.GetWorkflowJwt(),
			},
		},
	}, jobMeta.Address+"/grpc")

	ticker := time.NewTicker(5 * time.Second)

	failedPings := 0

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_, err := jobClient.Ping(ctx, connect.NewRequest(&pb.JobServicePingRequest{}))
			if err != nil {
				failedPings++
				if failedPings >= 3 {
					log.Error().Err(err).Msg("job is dead")
					handleUnresponsiveJob(ctx, h.orchestratorClient, h.workflowMeta, jobMeta)
					return
				}
			} else {
				failedPings = 0
			}
		}
	}
}

func handleUnresponsiveJob(ctx context.Context, orchestratorClient pbConnect.OrchestratorServiceClient, workflowMeta *pb.WorkflowMeta, jobMeta *pb.JobMeta) {

	_, err := orchestratorClient.CreateWorkflowAlert(ctx, connect.NewRequest(&pb.OrchestratorServiceCreateWorkflowAlertRequest{
		WorkflowMeta: workflowMeta,
		Alert: &pb.WorkflowAlert{
			Type:    pb.WorkflowAlert_TYPE_ERROR,
			Title:   "Unresponsive Job",
			Message: fmt.Sprintf("%s was unreachable, maybe the runner ran out of memory?", jobMeta.Name),
		}}))
	if err != nil {
		log.Error().Err(err).Msg("creating workflow alert")
	}

	_, err = orchestratorClient.FinishJob(ctx, connect.NewRequest(&pb.OrchestratorServiceFinishJobRequest{
		WorkflowMeta: workflowMeta,
		JobMeta:      jobMeta,
		Conclusion:   pb.Conclusion_CONCLUSION_FAILURE,
	}))
	if err != nil {
		log.Error().Err(err).Msg("finishing job")
	}
}
