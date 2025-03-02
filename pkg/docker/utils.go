package docker

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"regexp"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	nanoid "github.com/matoous/go-nanoid/v2"
	"github.com/rs/zerolog/log"
	"github.com/alfiejones/panda-ci/pkg/utils"
	pb "github.com/alfiejones/panda-ci/proto/go/v1"
)

type DeleteContainerFilter struct {
	WorkflowID *string
	JobID      *string
	TaskID     *string
}

func DeleteContainers(ctx context.Context, dockerClient *client.Client, filter DeleteContainerFilter) error {
	defer utils.MeasureTime(time.Now(), "DeleteContainers")

	containerFilters := filters.NewArgs()

	containerFilters.Add("label", "panda-ci=true")

	if filter.WorkflowID != nil {
		containerFilters.Add("label", fmt.Sprintf("workflowID=%s", *filter.WorkflowID))
	}
	if filter.JobID != nil {
		containerFilters.Add("label", fmt.Sprintf("jobID=%s", *filter.JobID))
	}
	if filter.TaskID != nil {
		containerFilters.Add("label", fmt.Sprintf("taskID=%s", *filter.TaskID))
	}

	containers, err := dockerClient.ContainerList(ctx, container.ListOptions{
		Filters: containerFilters,
	})
	if err != nil {
		return err
	}

	var deleteErr error
	for _, dockerContainer := range containers {
		if err := dockerClient.ContainerRemove(ctx, dockerContainer.ID, container.RemoveOptions{
			Force: true,
		}); err != nil {
			// We want to continue deleting the rest of the containers
			// This means we will only return the last error that occurred
			deleteErr = err
			log.Error().Err(err).Msgf("Failed to remove container: %s", dockerContainer.ID)
		}
	}

	return deleteErr
}

type DeleteVolumeFilter struct {
	WorkflowID *string
	JobID      *string
}

func DeleteVolumes(ctx context.Context, dockerClient *client.Client, filter DeleteVolumeFilter) error {
	defer utils.MeasureTime(time.Now(), "DeleteVolumes")

	volumeFilters := filters.NewArgs()

	volumeFilters.Add("label", "panda-ci=true")

	if filter.WorkflowID != nil {
		volumeFilters.Add("label", fmt.Sprintf("workflowID=%s", *filter.WorkflowID))
	}
	if filter.JobID != nil {
		volumeFilters.Add("label", fmt.Sprintf("jobID=%s", *filter.JobID))
	}

	volumes, err := dockerClient.VolumeList(ctx, volume.ListOptions{
		Filters: volumeFilters,
	})
	if err != nil {
		return err
	}

	var deleteErr error
	for _, volume := range volumes.Volumes {
		if err := dockerClient.VolumeRemove(ctx, volume.Name, true); err != nil {
			// We want to continue deleting the rest of the volumes
			// This means we will only return the last error that occurred
			deleteErr = err
			log.Error().Err(err).Msgf("Failed to remove volume: %s", volume.Name)
		}
	}

	return deleteErr
}

func GenerateSafeDockerName(workflowID string, tag string) (string, error) {
	id, err := nanoid.New(6)
	if err != nil {
		return "", err
	}

	regex := regexp.MustCompile(`[^a-zA-Z0-9_.-]+`)

	safeTag := regex.ReplaceAllString(tag, "")

	// TODO - avoid using banned characters and length restrictions

	prefix := fmt.Sprintf("panda-ci-%s-%s", workflowID, safeTag)

	if len(prefix) > 249 {
		prefix = prefix[:249]
	}

	return fmt.Sprintf("%s-%s", prefix, id), nil
}

func PullImage(ctx context.Context, dockerClient *client.Client, refStr string) error {
	log.Debug().Msgf("Pulling image: %s", refStr)

	filterArgs := filters.NewArgs()
	filterArgs.Add("reference", refStr)

	res, err := dockerClient.ImageList(ctx, image.ListOptions{
		Filters: filterArgs,
	})
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	if len(res) == 0 {
		// We want to wait for the image to be pulled since we don't have one locally
		// If we have one locally then we use that and pull the latest in the background for future runs
		// This enables us very fast start times without compromising on being out of sync with images
		wg.Add(1)
	}

	go func() {
		// TODO - we should skip this if it was last done in the last 5 minutes

		if len(res) == 0 {
			defer wg.Done()
		}

		imgReader, err := dockerClient.ImagePull(context.Background(), refStr, image.PullOptions{})
		if err != nil {
			log.Err(err).Msg("Failed to pull docker image")
			errChan <- err
			return
		}

		// Wait for the image to be pulled
		_, err = io.ReadAll(imgReader)
		if err != nil {
			imgReader.Close()
			log.Err(err).Msg("Failed to read docker pull image response")
			errChan <- err
			return
		}

		if err := imgReader.Close(); err != nil {
			log.Err(err).Msg("Failed to close docker image reader")
			errChan <- err
			return
		}

		log.Info().Msgf("Latest version of docker image %s pulled", refStr)
		errChan <- nil
	}()

	if len(res) > 0 {
		return nil
	}

	wg.Wait()
	return <-errChan
}

func WaitForContainerToFinish(ctx context.Context, dockerClient *client.Client, containerID string) error {
	statusCh, errCh := dockerClient.ContainerWait(ctx, containerID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		return err
	case <-statusCh:
		return nil
	}
}

func StreamLogs(reader io.Reader, callback func(logType pb.LogMessage_ExecData_Type, data []byte) error) error {
	hdr := make([]byte, 8)
	for {
		_, err := io.ReadFull(reader, hdr)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		logType := pb.LogMessage_ExecData_TYPE_STDOUT
		if hdr[0] != 1 {
			logType = pb.LogMessage_ExecData_TYPE_STDERR
		}

		count := binary.BigEndian.Uint32(hdr[4:])
		dat := make([]byte, count)
		if _, err := io.ReadFull(reader, dat); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		if err := callback(logType, dat); err != nil {
			return err
		}
	}
}

func CreateVolume(ctx context.Context, workflowID string, dockerClient *client.Client, tag string, labels map[string]string) (string, error) {
	name, err := GenerateSafeDockerName(workflowID, tag)
	if err != nil {
		return "", err
	}

	labels["panda-ci"] = "true"
	labels["workflowID"] = workflowID

	if _, err := dockerClient.VolumeCreate(ctx, volume.CreateOptions{
		Name:   name,
		Labels: labels,
	}); err != nil {
		return "", err
	}

	return name, nil
}
