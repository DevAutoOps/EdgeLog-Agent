package variable

import (
	"agent/global/consts"
	"agent/model"
	"agent/model/check"
	"agent/model/config"
	"go.uber.org/zap"
	"log"
	"net"
	"os"
	"path/filepath"
)

var (
	BasePath       string // Defines the root directory of the project
	ConfigFilePath string
	// Global log pointer
	ZapLog *zap.Logger
	// Global configuration object
	Config *config.Config
	// Communication dictionary
	CommunicationDict map[*net.Conn]chan bool
	// Command data buffer
	CommandBufferChan chan *model.DataBufferClient
	// Monitoring data buffer
	MonitorDataBufferChan chan *model.DataBufferClient
	// CMD command execution error counter
	CmdErrCounter = 0
	CmdErrLimit   = 2160
	// Test result object
	Result interface{}
	// save last monitor data
	LastMonitorData *check.MonitorRecord
)

func init() {
	// 1.Initialize program root directory
	if path, err := os.Getwd(); err == nil {
		BasePath = path
		ConfigFilePath = filepath.Join(BasePath, "config", "config.yml")
	} else {
		log.Fatal(consts.ErrorsBasePath)
	}
}
