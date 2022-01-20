package log_collector

import (
	"agent/common"
	"agent/global/consts"
	"agent/global/journal"
	"agent/global/variable"
	"agent/model/logger/log_file"
	"errors"
	"fmt"
	"path/filepath"
)

func search() ([]log_file.LogFile, error) {
	includeFileMap = nil
	includeFileMap = make(map[string]struct{})
	localLogFileMap = nil
	localLogFileMap = make(map[string]map[string]struct{})
	if variable.Config.Modules.LogCollector.Targets == nil ||
		len(variable.Config.Modules.LogCollector.Targets) <= 0 ||
		variable.Config.Modules.LogCollector.TargetConfs == nil ||
		len(variable.Config.Modules.LogCollector.TargetConfs) <= 0 {
		common.Output(fmt.Sprintf("[search log file] The target to be searched is missing in the configuration file. Please fill in 'Targets' and 'Targetconfs'"), journal.Warn)
		return nil, errors.New("the target to be searched is missing in the configuration file. Please fill in 'Targets' and 'Targetconfs'")
	}
	if len(variable.Config.Modules.LogCollector.Targets) !=
		len(variable.Config.Modules.LogCollector.TargetConfs) {
		common.Output(fmt.Sprintf("[search log file] The number of 'Targets' and 'Targetconfs' filled in the configuration file is inconsistent"), journal.Warn)
		return nil, errors.New("the number of 'Targets' and 'Targetconfs' filled in the configuration file is inconsistent")
	}

	logFileList := make([]log_file.LogFile, 0)
	for i := 0; i < len(variable.Config.Modules.LogCollector.Targets); i++ {
		currentTarget := variable.Config.Modules.LogCollector.Targets[i]
		currentTargetConf := variable.Config.Modules.LogCollector.TargetConfs[i]
		if len(currentTarget) <= 0 || len(currentTargetConf) <= 0 {
			common.Output(fmt.Sprintf("[search log file] The content of 'Targets' or 'Targetconfs' in the configuration file is empty"), journal.Warn)
			continue
		}
		logFiles, err := searchLogs("Nginx", currentTargetConf)
		if err != nil {
			common.Output(fmt.Sprintf("[search log file] The log file related configuration could not be found in the configuration file '%s' to search'%s'", currentTarget, currentTargetConf), journal.Warn)
		} else {
			logFileList = append(logFileList, logFiles...)
		}
	}
	if len(logFileList) <= 0 {
		return nil, errors.New(consts.ErrorsNotFoundFile)
	} else {
		return logFileList, nil
	}
}

func searchLogs(logType, confPath string) ([]log_file.LogFile, error) {
	confDir := filepath.Dir(confPath)
	logFiles := make([]log_file.LogFile, 0)
	localLogFileMap[logType] = nil
	localLogFileMap[logType] = make(map[string]struct{})
	_ = (&NginxConfSearcher{Type: logType}).scan(confPath, confDir)
	if len(localLogFileMap[logType]) <= 0 {
		return nil, errors.New(consts.ErrorsNotFoundFile)
	}
	for logFile := range localLogFileMap[logType] {
		if len(logFile) > 0 {
			logFiles = append(logFiles, log_file.LogFile{FilePath: logFile})
		}
	}

	if len(logFiles) <= 0 {
		return nil, errors.New(consts.ErrorsNotFoundFile)
	} else {
		return logFiles, nil
	}
}
