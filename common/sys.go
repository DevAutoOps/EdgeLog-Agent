package common

import (
	"agent/global/consts"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
)

func GetOs() string {
	return strings.ToLower(runtime.GOOS)
}

func CmdLinux(cmdStr string) (result string, success bool, err error) {
	result = ""
	success = false
	cmd := exec.Command("bash", "-c", cmdStr)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		//Setpgid: true,	// Comment out this line. This structure does not have this attribute under Windows compilation
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
	result = string(out)
	success = true
	return
}

func CmdWindows(cmdStr string) (result string, success bool, err error) {
	result = ""
	success = false
	cmd := exec.Command("cmd", "/C", cmdStr)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
	result = string(out)
	success = true
	return
}

func Cmd(cmdStr string) (result string, success bool, err error) {
	switch GetOs() {
	case consts.OsWindows:
		return CmdWindows(cmdStr)
	case consts.OsLinux, consts.OsMac:
		return CmdLinux(cmdStr)
	default:
		return "", false, errors.New("I won't support it")
	}
}

func Cmd2(cmdStr, os string) (result string, success bool, err error) {
	switch os {
	case consts.OsWindows:
		return CmdWindows(cmdStr)
	case consts.OsLinux, consts.OsMac:
		return CmdLinux(cmdStr)
	default:
		return "", false, errors.New("I won't support it")
	}
}

// IsProcessExist Determine whether the process exists
func IsProcessExist(appName string) bool {
	os := GetOs()
	cmdStr := ""
	if os == consts.OsWindows {
		cmdStr = fmt.Sprintf("tasklist | findstr %s", appName+".exe")
	} else if os == consts.OsLinux || os == consts.OsMac {
		cmdStr = fmt.Sprintf("ps -ef | grep %s | grep -v grep", appName)
	}
	exists := false
	resultStr, success, err := Cmd2(cmdStr, os)
	if err == nil && success && strings.Contains(resultStr, appName) {
		exists = true
	}
	return exists
}
