package server

import (
	"agent/common"
	"agent/global/consts"
	"agent/global/journal"
	"agent/global/variable"
	"agent/model"
	"agent/model/check"
	check2 "agent/module/check"
	"agent/module/encode"
	"bytes"
	"fmt"
	"github.com/robfig/cron/v3"
)

type MonitorDataCollector struct {
	crontab    *cron.Cron
	taskIds    map[uint16]cron.EntryID
	IsStart    bool
	checkCount uint
}

func (m *MonitorDataCollector) Start() {
	if !variable.Config.Modules.Monitor.Enable {
		return
	}
	if m.IsStart {
		return
	}

	if m.crontab == nil {
		m.crontab = cron.New(cron.WithSeconds())
	}
	freq := variable.Config.Modules.Monitor.Freq
	if freq == 0 {
		freq = consts.DefaultFreq
	}
	freq = freq / 1000
	spec := fmt.Sprintf("*/%d * * * * ?", freq)
	taskId, err := m.crontab.AddFunc(spec, m.run)
	if err != nil {
		common.Output(fmt.Sprintf("[MonitorDataCollector Start] add task error: %s", err), journal.Error)
		return
	}
	if m.taskIds == nil {
		m.taskIds = make(map[uint16]cron.EntryID)
	}
	m.taskIds[consts.MonitorData] = taskId
	m.crontab.Start()
	m.IsStart = true
	select {}
}

func (m *MonitorDataCollector) Stop() {
	if m.crontab != nil {
		m.crontab.Stop()
		m.crontab = nil
	}
	m.IsStart = false
}

func (m *MonitorDataCollector) run() {
	//_, ok := variable.Result.(*check.CheckedResult)
	//if !ok {
	//	common.Output("[MonitorDataCollector run] get result object err", journal.Error)
	//	return
	//}

	// 监测命令暂不支持Windows
	//if common.GetOs() != consts.OsWindows {
	if variable.Config.Modules.Monitor.CheckFreq > 1 {
		// Circular judgment
		if m.checkCount%variable.Config.Modules.Monitor.CheckFreq == 0 {
			variable.Config.Modules.Monitor.Check = true
		} else {
			variable.Config.Modules.Monitor.Check = false
			//setCheckList(checkedResult, consts.MonitorCpuUsage, 255, false)
		}
	}
	if variable.Config.Modules.Monitor.Check {
		(&check2.ServerCheck{}).Check()
	}
	//}

	(&check2.AppCheck{}).Check()

	dataBuffer, err := encode.EncodeMonitor()
	if err != nil {
		common.Output(fmt.Sprintf("[MonitorDataCollector run] encode monitor data error: %s", err), journal.Error)
		return
	}
	if dataBuffer == nil {
		common.Output(fmt.Sprintf("[MonitorDataCollector run] encode monitor data is nil"), journal.Warn)
		return
	}
	monitorData := &model.DataBufferClient{
		DataBuffer: bytes.NewBuffer(dataBuffer.Bytes()),
		Length:     uint32(dataBuffer.Len()),
	}
	variable.MonitorDataBufferChan <- monitorData
	if variable.Config.Logs.PrintLevel == 1 {
		common.Output(fmt.Sprintf("[MonitorDataCollector run] encode data length: %d, content: %v", dataBuffer.Len(), dataBuffer.Bytes()))
	}
	if variable.Config.Modules.Monitor.CheckFreq > 1 {
		if m.checkCount%variable.Config.Modules.Monitor.CheckFreq == 0 {
			m.checkCount = 0
		}
		m.checkCount++
	}
}

func setCheckList(checkedResult *check.CheckedResult, startNum uint8, endNum uint8, status bool) {
	checkedResult.Lock()
	for i := startNum; i <= endNum; i++ {
		checkedResult.CheckList[i] = status
	}
	checkedResult.Unlock()
}
