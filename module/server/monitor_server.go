package server

import (
	"agent/common"
	"agent/global/journal"
	"agent/global/variable"
	"fmt"
	"net"
	"time"
)

type MonitorServer struct {
	IsStart bool
}

func (m *MonitorServer) Start() {
	if !variable.Config.Modules.Monitor.Enable {
		return
	}
	if m.IsStart {
		return
	}

Restart:
	if len(variable.Config.Server.Ip) <= 0 || len(variable.Config.Server.MonitorPort) <= 0 {
		common.Output("[MonitorServer Start] ip or monitor port is empty", journal.Warn)
		return
	}
	host := variable.Config.Server.Ip + variable.Config.Server.MonitorPort
	conn, err := net.Dial("tcp", host)
	if err != nil {
		common.Output(fmt.Sprintf("[MonitorServer Start] remote host connection error: %s", err), journal.Error)
		time.Sleep(5e9)
		goto Restart
	}
	defer conn.Close()

	m.IsStart = true
	for {
		if !m.IsStart {
			break
		}
		monitorDataBuffer := <-variable.MonitorDataBufferChan
		if monitorDataBuffer.DataBuffer.Len() > 0 && monitorDataBuffer.DataBuffer.Len() == int(monitorDataBuffer.Length) {
			data := make([]byte, monitorDataBuffer.DataBuffer.Len())
			n, err := monitorDataBuffer.DataBuffer.Read(data)
			if n <= 0 || err != nil {
				continue
			}
			_, err = conn.Write(data) // Send data
			if err != nil {
				common.Output(fmt.Sprintf("[MonitorServer Start] write monitor data error: %s", err), journal.Error)
				time.Sleep(5e9)
				goto Restart
			}
		} else {
			time.Sleep(1e8)
		}
	}
}

func (m *MonitorServer) Stop() {
	m.IsStart = false
}
