package main

import (
	_ "agent/bootstrap"
	"agent/common"
	"agent/global/journal"
	_ "agent/global/logger"
	_ "agent/global/variable"
	"agent/module/monitor"
	"agent/module/server"
	"fmt"
	"runtime/debug"
)

var Version string

func main() {
	defer func() {
		err := recover()
		if err != nil {
			common.Output(fmt.Sprintf("[main] app recover, %s", debug.Stack()), journal.Error)
		}
	}()

	common.Output(fmt.Sprintf("[main] Current version number:%s", Version))

	go (&monitor.Monitor{}).AppShutDown()
	go (&server.MonitorDataCollector{IsStart: false}).Start()
	go (&server.MonitorServer{IsStart: false}).Start()
	(&server.LogCollectorServer{}).Start()
	(&server.LogSender{}).Start()
}
