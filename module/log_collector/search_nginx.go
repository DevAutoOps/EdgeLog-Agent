package log_collector

import (
	"agent/common"
	"agent/global/journal"
	"agent/utils/ngconf"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type NginxConfSearcher struct {
	Type string
}

func (n *NginxConfSearcher) scan(confPath, confDir string) []string {
	logFiles := make([]string, 0)
	root, err := ngconf.DecodeFile(confPath)
	if err != nil {
		common.Output(fmt.Sprintf("[NginxConfSearcher scan] decode file (%s) error:%s", confPath, err), journal.Error)
		return nil
	}
	if root == nil {
		common.Output(fmt.Sprintf("[NginxConfSearcher scan] (%s) object is null after decode", confPath), journal.Warn)
		return nil
	}
	localAccess := make([]string, 0)
	localAccess = n.collectAll(root, confDir)
	includeNodes := root.GetDeepByDirective("include")
	if includeNodes != nil && len(includeNodes) > 0 {
		includeAccessFiles := n.includeHandler(includeNodes, confDir)
		if includeAccessFiles != nil && len(includeAccessFiles) > 0 {
			localAccess = append(localAccess, includeAccessFiles...)
		}
	}
	if localAccess != nil && len(localAccess) > 0 {
		logFiles = append(logFiles, localAccess...)
	}
	return logFiles
}

func (n *NginxConfSearcher) collectAll(node *ngconf.Node, confDir string) []string {
	accessFiles := make([]string, 0)
	accessLogs := node.GetDeepByDirective("access_log")
	if accessLogs == nil || len(accessLogs) <= 0 {
		//common.Output(fmt.Sprintf("[CommonAction collectLogFiles] no found access_log"), journal.Warn)
		return accessFiles
	}
	for _, accessLog := range accessLogs {
		accessFile := accessLog.Args[0]
		if len(accessFile) <= 0 {
			common.Output(fmt.Sprintf("[NginxConfSearcher collectAll] access_log file is empty"), journal.Warn)
			continue
		}
		if accessFile == "off" {
			continue
		}
		if !strings.HasPrefix(accessFile, "/") {
			accessFile = filepath.Join(confDir, "..", accessFile)
		}
		common.Output(fmt.Sprintf("[NginxConfSearcher collectAll] find log file %s", accessFile))
		accessFiles = append(accessFiles, accessFile)
		if _, ok := localLogFileMap[n.Type][accessFile]; !ok {
			localLogFileMap[n.Type][accessFile] = struct{}{}
		}
	}
	return accessFiles
}

func (n *NginxConfSearcher) includeHandler(includeNodes []*ngconf.Node, confDir string) []string {
	accessFiles := make([]string, 0)
	if includeNodes != nil && len(includeNodes) > 0 {
		for _, includeNode := range includeNodes {
			if includeNode.Args == nil || len(includeNode.Args) <= 0 ||
				len(includeNode.Args[0]) <= 0 ||
				strings.HasSuffix(includeNode.Args[0], "mime.types") {
				continue
			}
			includeFile := includeNode.Args[0]
			if !strings.HasPrefix(includeFile, "/") {
				includeFile = filepath.Join(confDir, includeFile)
			}
			includeFileBase := filepath.Base(includeFile)
			var includeAccessFiles []string
			if strings.Contains(includeFileBase, "*") {
				includeDir := filepath.Dir(includeFile)
				if _, ok := includeFileMap[includeDir]; ok {
					continue
				} else {
					includeFileMap[includeDir] = struct{}{}
				}
				includeAccessFiles = n.scanIncludeDir(includeFile, confDir)
			} else {
				if _, ok := includeFileMap[includeFile]; ok {
					continue
				} else {
					includeFileMap[includeFile] = struct{}{}
				}
				includeAccessFiles = n.scanIncludeFile(includeFile, confDir)
			}
			if includeAccessFiles != nil && len(includeAccessFiles) > 0 {
				accessFiles = append(accessFiles, includeAccessFiles...)
			}
		}
	}
	return accessFiles
}

func (n *NginxConfSearcher) scanIncludeDir(filePath string, confDir string) []string {
	accessFiles := make([]string, 0)
	fileDir := filepath.Dir(filePath)
	rd, err := ioutil.ReadDir(fileDir)
	if err != nil {
		common.Output(fmt.Sprintf("[NginxConfSearcher scanIncludeDir] read dir %s err: %s", fileDir, err), journal.Error)
		return nil
	}
	for _, fi := range rd {
		if fi.IsDir() {
			continue
		}
		targetFile := filepath.Join(fileDir, fi.Name())
		includeAccessFiles := n.scanIncludeFile(targetFile, confDir)
		if includeAccessFiles != nil && len(includeAccessFiles) > 0 {
			accessFiles = append(accessFiles, includeAccessFiles...)
		}
	}
	return accessFiles
}

func (n *NginxConfSearcher) scanIncludeFile(filePath string, confDir string) []string {
	accessFiles := make([]string, 0)
	root, err := ngconf.DecodeFile(filePath)
	if err != nil {
		common.Output(fmt.Sprintf("[NginxConfSearcher scanIncludeFile] decode file (%s) error:%s", filePath, err), journal.Error)
		return nil
	}
	if root == nil {
		common.Output(fmt.Sprintf("[NginxConfSearcher scanIncludeFile] (%s) object is null after decode", filePath), journal.Warn)
		return nil
	}
	localAccess := make([]string, 0)
	localAccess = n.collectAll(root, confDir)
	includeNodes := root.GetDeepByDirective("include")
	if includeNodes != nil && len(includeNodes) > 0 {
		includeAccessFiles := n.includeHandler(includeNodes, confDir)
		if includeAccessFiles != nil && len(includeAccessFiles) > 0 {
			localAccess = append(localAccess, includeAccessFiles...)
		}
	}
	if localAccess != nil && len(localAccess) > 0 {
		accessFiles = append(accessFiles, localAccess...)
	}
	return accessFiles
}
