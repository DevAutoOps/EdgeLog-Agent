package logger

import (
	"agent/model/logger/log_file"
	"agent/proto/log"
)

const (
	OffsetRecordFile = "./offset_record"
)

var (
	LogEndChan        chan struct{}
	TailList          []*log_file.LogTail
	LogSenderDataChan chan *log.LogData
)

func init() {
	LogEndChan = make(chan struct{})
	TailList = make([]*log_file.LogTail, 0)
	LogSenderDataChan = make(chan *log.LogData, 100)
}
