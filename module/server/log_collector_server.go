package server

import (
	"agent/common"
	"agent/global/journal"
	"agent/global/logger"
	"agent/global/variable"
	"agent/module/log_collector"
)

type LogCollectorServer struct {
	IsStart bool
}

func (l *LogCollectorServer) Start() {
	if l.IsStart {
		return
	}
	if !variable.Config.Modules.LogCollector.Enable {
		common.Output("[LogCollectorServer Start] log collector enable is false", journal.Warn)
		return
	}
	l.IsStart = true

	go (&log_collector.LogReader{}).Start()
}

func (l *LogCollectorServer) Stop() {
	l.IsStart = false
	close(logger.LogEndChan)
}
