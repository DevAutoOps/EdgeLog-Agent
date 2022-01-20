package sys_log_hook

import (
	"go.uber.org/zap/zapcore"
)

// GoSkeleton System operation log hook function
// 1.A single log is a structure format. This function intercepts each log and you can carry out subsequent processing, such as pushing it to Alibaba cloud log management panel, elasticsearch log library, etc

func ZapLogHandler(entry zapcore.Entry) error {

	// Parameter entry introduction
	// entry  The parameter is a single log structure, which mainly includes the following fields:
	//Level      Log level
	//Time       Current time
	//LoggerName  Log name
	//Message    Log content
	//Caller     Call path of each file
	//Stack      Code call stack

	// Here, a co process is started, and hook will not affect the program performance at all，
	go func(paramEntry zapcore.Entry) {
		//fmt.Println(" GoSkeleton  hook ....，You can continue to process system logs here....")
		//fmt.Printf("%#+v\n", paramEntry)
	}(entry)
	return nil
}
