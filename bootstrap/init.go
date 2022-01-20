package bootstrap

import (
	"agent/common"
	"agent/global/variable"
	"agent/model"
	"agent/model/check"
	"agent/model/config"
	"agent/utils/sys_log_hook"
	"agent/utils/zap_factory"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
)

func init() {
	initConfig()
	checkDir()
	variable.ZapLog = zap_factory.CreateZapFactory(sys_log_hook.ZapLogHandler)
	variable.CommunicationDict = make(map[*net.Conn]chan bool)
	variable.CommandBufferChan = make(chan *model.DataBufferClient, 1)
	variable.MonitorDataBufferChan = make(chan *model.DataBufferClient, 1)
	variable.Result = &check.CheckedResult{CheckList: make(map[uint8]bool)}
}

func initConfig() {
	var config2 config.Config
	configContent, err := ioutil.ReadFile(variable.ConfigFilePath)
	if err != nil {
		fmt.Println("not found config.yml, ", err)
		os.Exit(1)
	}
	err = yaml.Unmarshal(configContent, &config2)
	if err != nil {
		fmt.Println("yml Unmarshal error, ", err)
		os.Exit(1)
	}
	variable.Config = &config2
}

func checkDir() {
	logsPath := filepath.Join(".", variable.Config.Logs.LogName)
	logsDir := filepath.Dir(logsPath)
	if !common.IsExists(logsDir, false) {
		err := common.CreateFolder(logsDir)
		if err != nil {
			fmt.Println("create dir fail,", err)
		}
	}
}
