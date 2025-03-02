package stream

import (
	"time"

	utilsCSV "github.com/alfiejones/panda-ci/pkg/utils/csv"
	pb "github.com/alfiejones/panda-ci/proto/go/v1"
)

type LogStream struct {
	*Stream
}

type LogType string

const (
	LogTypeStdout LogType = "stdout"
	LogTypeStderr LogType = "stderr"
	LogTypeExit   LogType = "exit"
)

type Log struct {
	Timestamp time.Time
	Type      LogType
	Data      string
}

func ProtoLogTypeToLogType(protoType pb.LogMessage_ExecData_Type) LogType {
	switch protoType {
	case pb.LogMessage_ExecData_TYPE_STDERR:
		return LogTypeStderr
	case pb.LogMessage_ExecData_TYPE_STDOUT:
		return LogTypeStdout
	default:
		return LogTypeExit
	}
}

func (ls *LogStream) WriteLog(log Log) error {

	entry, err := utilsCSV.FormatCSVRow([]string{log.Timestamp.Format(time.RFC3339Nano), string(log.Type), log.Data})
	if err != nil {
		return err
	}

	ls.Write(entry)

	return nil
}

func NewLogStream(inital []string) *LogStream {
	logStore := &LogStream{
		Stream: NewStream(inital),
	}

	return logStore
}
