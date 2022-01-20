package log_tail

import "github.com/hpcloud/tail"

type LogTail struct {
	FilePath string     // File path
	LogType  string     // Log file type
	Tail     *tail.Tail // Tail object
}
