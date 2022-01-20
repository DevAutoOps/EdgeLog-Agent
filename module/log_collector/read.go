package log_collector

import (
	"agent/common"
	"agent/global/journal"
	"agent/global/logger"
	"agent/model/logger/log_file"
	"agent/model/logger/offset_record"
	"agent/proto/log"
	"encoding/json"
	"fmt"
	"github.com/hpcloud/tail"
	"io"
	"os"
	"time"
)

type LogReader struct {
	IsStart bool
}

var localLogFileMap map[string]map[string]struct{}
var includeFileMap map[string]struct{}

func (r *LogReader) Init() bool {
	files, err := search()
	if err != nil {
		common.Output(fmt.Sprintf("[LogReader Init] search log file fail, err: %s", err), journal.Warn)
		return false
	}
	if files == nil || len(files) == 0 {
		common.Output("[LogReader Init] not found log file", journal.Warn)
		return false
	}

	for _, file := range files {
		location := &tail.SeekInfo{Offset: 0, Whence: io.SeekEnd}
		if common.IsExists(logger.OffsetRecordFile, true) {
			fileByte, err := common.ReadFileByte(logger.OffsetRecordFile)
			if err == nil && len(fileByte) > 0 {
				var offsetRecords []offset_record.OffsetRecord
				err = json.Unmarshal(fileByte, &offsetRecords)
				if err == nil {
					for _, offsetRecord := range offsetRecords {
						if offsetRecord.FilePath == file.FilePath && offsetRecord.Offset > 0 {
							location = &tail.SeekInfo{Offset: offsetRecord.Offset, Whence: io.SeekCurrent}
							break
						}
					}
				}
			}
		}
		config := tail.Config{
			ReOpen:    true, // Reopen
			Follow:    true, // Follow
			Location:  location,
			MustExist: false, // No error is reported if the file does not exist
			Poll:      true,
		}
		tailFile, err := tail.TailFile(file.FilePath, config)
		if err != nil {
			common.Output(fmt.Sprintf("[LogReader Init] tail file (%s) failed, err: %s", file.FilePath, err), journal.Error)
			continue
		}
		logTail := &log_file.LogTail{}
		logTail.FilePath = file.FilePath
		logTail.Tail = tailFile
		logger.TailList = append(logger.TailList, logTail)
	}
	if common.IsExists(logger.OffsetRecordFile, true) {
		// Delete record file
		_ = os.Remove(logger.OffsetRecordFile)
	}
	return true
}

func (r *LogReader) Start() {
	if r.IsStart {
		return
	}
	r.IsStart = true

	if !r.Init() {
		common.Output("[LogReader Start] init fail, exit.", journal.Warn)
		return
	}
	if len(logger.TailList) <= 0 {
		common.Output("[LogReader Start] log file not found, exit.", journal.Warn)
		return
	}

Loop:
	for {
		for _, tail := range logger.TailList {
			if tail == nil || tail.Tail == nil {
				continue
			}
			select {
			case _, ok := <-logger.LogEndChan:
				if !ok {
					common.Output("[LogReader Start] stop read.", journal.Warn)
					break Loop
				}
			case line, ok := <-tail.Tail.Lines:
				if !ok {
					common.Output(fmt.Sprintf("[LogReader Start] tail file close reopen, filename:%s\n", tail.Tail.Filename), journal.Warn)
					continue
				}
				go r.parse(line)
			default:
				time.Sleep(time.Millisecond * 100)
			}
		}
	}
}

func (r *LogReader) Stop() {
	r.IsStart = false
}

func (r *LogReader) parse(line *tail.Line) {
	if len(line.Text) <= 0 {
		return
	}

	logger.LogSenderDataChan <- r.setData(line.Text)
}

func (r *LogReader) setData(logText string) *log.LogData {
	var data log.LogData
	data.Log = logText
	return &data
}
