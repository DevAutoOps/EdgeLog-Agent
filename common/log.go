package common

import (
	"agent/global/journal"
	"agent/global/variable"
	"fmt"
	"time"
)

// Output Format output
func Output(s string, args ...string) {
	if len(s) > 0 {
		if !variable.Config.Logs.CloseStandardOutput {
			currentTime := time.Now()
			tm := fmt.Sprintf("[%d-%d-%d %d:%d:%d.%d] ", currentTime.Year(), currentTime.Month(), currentTime.Day(), currentTime.Hour(), currentTime.Minute(), currentTime.Second(), currentTime.Nanosecond())
			fmt.Println(tm, s)
		}
		if variable.Config.Logs.Enable {
			if len(args) > 0 && args[0] == journal.Error {
				variable.ZapLog.Error(s)
			} else if len(args) > 0 && args[0] == journal.Warn {
				variable.ZapLog.Warn(s)
			} else if len(args) > 0 && args[0] == journal.Panic {
				// Will be forced to withdraw app
				variable.ZapLog.Panic(s)
			} else {
				variable.ZapLog.Info(s)
			}
		}
	}
}
