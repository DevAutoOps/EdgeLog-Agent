package check

import "sync"

type CheckedResult struct {
	CheckList map[uint8]bool
	MonitorData
	AppResult
	sync.RWMutex
}
