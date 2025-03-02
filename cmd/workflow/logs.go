package main

import (
	"bytes"
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/alfiejones/panda-ci/pkg/stream"
	"github.com/alfiejones/panda-ci/pkg/uploads"
	"github.com/alfiejones/panda-ci/pkg/utils"
	pb "github.com/alfiejones/panda-ci/proto/go/v1"
)

type ArrayWriter struct {
	logStream *stream.LogStream
}

func (aw *ArrayWriter) Write(p []byte) (int, error) {

	// TODO - infer the log type from the data
	// we'd have to attempt to parse it and get the level from the json

	log := stream.Log{Timestamp: time.Now(), Type: "stdout", Data: string(p)}

	if err := aw.logStream.WriteLog(log); err != nil {
		return 0, err
	}

	return len(p), nil
}

func initLogs() *ArrayWriter {
	logStore := &ArrayWriter{
		logStream: stream.NewLogStream([]string{"timestamp,type,data\n"}),
	}

	multi := zerolog.MultiLevelWriter(os.Stdout, logStore)

	log.Logger = zerolog.New(multi).With().Logger()

	return logStore
}

func (aw *ArrayWriter) UploadLogs(ctx context.Context, config *pb.WorkflowRunnerInitConfig) error {
	defer utils.MeasureTime(time.Now(), "uploading logs")

	var buffer bytes.Buffer

	for _, entry := range aw.logStream.Entries() {
		buffer.WriteString(entry)
	}

	if err := uploads.UploadFile(ctx, config.PresignedOutputUrl, &buffer, "text/csv"); err != nil {
		log.Error().Err(err).Msg("uploading logs")
		return err
	}

	return nil
}
