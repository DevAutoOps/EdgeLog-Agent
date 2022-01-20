package server

import (
	"agent/common"
	"agent/global/journal"
	"agent/global/logger"
	"agent/global/variable"
	"agent/module/encode"
	"agent/proto/log"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"google.golang.org/protobuf/proto"
	"net"
	"time"
)

type LogSender struct {
	IsStart bool
	Conn    net.Conn
}

func (s *LogSender) Start() {
	if s.IsStart {
		return
	}
Restart:
	if len(variable.Config.Server.Ip) <= 0 || len(variable.Config.Server.LogPort) <= 0 {
		common.Output("[LogSender Start] ip or log port is empty", journal.Warn)
		return
	}
	host := variable.Config.Server.Ip + variable.Config.Server.LogPort
	var err error
	s.Conn, err = net.Dial("tcp", host)
	if err != nil {
		common.Output(fmt.Sprintf("[LogSender Start] remote host (%s) connection error: %s", host, err), journal.Error)
		time.Sleep(5e9)
		goto Restart
	}
	defer s.Conn.Close() // Close connection

	s.IsStart = true
	isStop := false
Loop:
	for {
		if !s.IsStart {
			break
		}

		select {
		case _, ok := <-logger.LogEndChan:
			if !ok {
				common.Output("[LogSender Start] notify stop sender.", journal.Warn)
				isStop = true
				if len(logger.LogSenderDataChan) == 0 {
					break Loop
				}
			}
		case logSenderData, ok := <-logger.LogSenderDataChan:
			if !ok && isStop {
				common.Output("[LogSender Start] stop write.", journal.Warn)
				break Loop
			}
			if !ok {
				common.Output("[LogSender Start] need not sender log data, wait 100ms.", journal.Warn)
				time.Sleep(1e8)
				continue
			}
			if variable.Config.Logs.PrintLevel == 1 {
				common.Output(fmt.Sprintf("log -> %#v\n", logSenderData))
			}
			data, err := s.encode(logSenderData)
			if err != nil {
				common.Output(fmt.Sprintf("[LogSender Start] data encode error: %s", err), journal.Error)
				continue
			}
			_, err = s.Conn.Write(data) // Send data
			if err != nil {
				common.Output(fmt.Sprintf("[LogSender Start] write data error: %s", err), journal.Error)
				logger.LogSenderDataChan <- logSenderData
				time.Sleep(5e9)
				goto Restart
			}
		}
	}
}

func (s *LogSender) Stop() {
	s.IsStart = false
	s.Conn.Close()
}

func (s *LogSender) encode(logSenderData *log.LogData) ([]byte, error) {
	data, err := proto.Marshal(logSenderData)
	if err != nil {
		common.Output(fmt.Sprintf("[LogSender encode] encode proto error: %s", err), journal.Error)
		return nil, err
	}
	if len(data) <= 0 {
		common.Output(fmt.Sprintf("[LogSender encode] encode proto is empty: %s", errors.New("encode proto is empty")), journal.Error)
		return nil, err
	}
	//common.Output(fmt.Sprintf("[LogSender encode] encode log data length: %d", len(data)))
	var gzipBuffer bytes.Buffer
	gzipWriter := gzip.NewWriter(&gzipBuffer)
	defer gzipWriter.Close()
	gzipWriter.Write(data)
	gzipWriter.Flush()
	gzipData := gzipBuffer.Bytes()
	//common.Output(fmt.Sprintf("[LogSender encode] encode log zip data length: %d", len(gzipData)))
	resultBuffer, err := encode.EncodeLog(gzipData, uint32(len(data)))
	if err != nil {
		common.Output(fmt.Sprintf("[LogSender encode] encode log data error: %s", err), journal.Error)
		return nil, err
	}
	return resultBuffer.Bytes(), nil
}

func (s *LogSender) encode2(logSenderData *log.LogData) ([]byte, error) {
	data, err := proto.Marshal(logSenderData)
	if err != nil {
		common.Output(fmt.Sprintf("[LogSender encode] encode proto error: %s", err), journal.Error)
		return nil, err
	}
	if len(data) <= 0 {
		common.Output(fmt.Sprintf("[LogSender encode] encode proto is empty: %s", errors.New("encode proto is empty")), journal.Error)
		return nil, err
	}
	resultBuffer, err := encode.EncodeLog2(data)
	if err != nil {
		return nil, err
	}
	return resultBuffer.Bytes(), nil
}
