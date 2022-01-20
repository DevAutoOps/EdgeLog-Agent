package log_file

import "github.com/hpcloud/tail"

type LogFile struct {
	FilePath string // File path
}

type LogTail struct {
	LogFile
	Tail *tail.Tail // Tail object
}
