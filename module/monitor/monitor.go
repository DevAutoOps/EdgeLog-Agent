package monitor

import (
	"agent/common"
	"agent/global/journal"
	"agent/global/logger"
	"agent/model/logger/offset_record"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Monitor struct {
}

func (m *Monitor) AppShutDown() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
Loop:
	for {
		select {
		case _, ok := <-logger.LogEndChan:
			if !ok {
				common.Output("[Monitor AppShutDown] stop monitor.")
				m.recordFile()
				break Loop
			}
		case s := <-c:
			common.Output(fmt.Sprintf("[Monitor AppShutDown] get signal: %v", s), journal.Warn)
			m.recordFile()
			os.Exit(0)
		}
	}
}

func (m *Monitor) recordFile() {
	offsetRecords := make([]offset_record.OffsetRecord, 0)
	for _, tail := range logger.TailList {
		if tail == nil || tail.Tail == nil {
			continue
		}
		offset, _ := tail.Tail.Tell()
		offsetRecords = append(offsetRecords, offset_record.OffsetRecord{FilePath: tail.Tail.Filename, Offset: offset})
	}
	if len(offsetRecords) > 0 {
		offsetJson, err := json.Marshal(offsetRecords)
		if err == nil {
			common.Output(fmt.Sprintf("[Monitor recordFile] save offset: %s\n", offsetJson))
			_ = common.WriteFileByte(logger.OffsetRecordFile, offsetJson)
			time.Sleep(time.Millisecond * 100)
		}
	}
}
