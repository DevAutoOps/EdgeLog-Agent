package check

import (
	"agent/common"
	"agent/global/consts"
	"agent/global/journal"
	"agent/global/variable"
	"agent/model/check"
	"fmt"
	"strings"
)

type AppCheck struct {
}

func (a *AppCheck) Check() {
	if variable.Config.Modules.Monitor.Apps == nil ||
		len(variable.Config.Modules.Monitor.Apps) <= 0 ||
		variable.Config.Modules.Monitor.AppProcess == nil ||
		len(variable.Config.Modules.Monitor.AppProcess) <= 0 {
		common.Output(fmt.Sprintf("[AppCheck Check] The apps to be monitored is missing in the configuration file. Please fill in 'Apps' and 'AppProcess'"), journal.Warn)
		return
	}
	if len(variable.Config.Modules.Monitor.Apps) != len(variable.Config.Modules.Monitor.AppProcess) {
		common.Output(fmt.Sprintf("[AppCheck Check] The number of 'Apps' and 'AppProcess' filled in the configuration file is inconsistent"), journal.Warn)
		return
	}

	for i := 0; i < len(variable.Config.Modules.Monitor.Apps); i++ {
		currentApp := variable.Config.Modules.Monitor.Apps[i]
		currentAppProcess := variable.Config.Modules.Monitor.AppProcess[i]
		if len(currentApp) <= 0 || len(currentAppProcess) <= 0 {
			common.Output(fmt.Sprintf("[AppCheck Check] The content of 'Apps' or 'AppProcess' in the configuration file is empty"), journal.Warn)
			continue
		}
		a.appCheck(currentAppProcess)
	}
}

func (a *AppCheck) appCheck(process string) {
	if checkedResult, ok := variable.Result.(*check.CheckedResult); ok {
		isExists := a.checkHandle(process)
		checkedResult.Lock()
		checkedResult.CheckList[consts.MonitorNginx] = true
		checkedResult.Unlock()
		if isExists {
			checkedResult.NginxStatus = 1
		} else {
			checkedResult.NginxStatus = 0
		}
	}
}

func (a *AppCheck) checkHandle(process string) bool {
	processList := strings.Split(process, ",")
	isExists := false
	for _, processName := range processList {
		if len(processName) <= 0 {
			continue
		}
		isExists = common.IsProcessExist(processName)
	}
	return isExists
}
