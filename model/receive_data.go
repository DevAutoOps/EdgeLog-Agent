package model

import "bytes"

type ReceiveDataProperties struct {
	IsBegin            bool
	IsFinish           bool
	TotalLength        uint32
	UnReadFinishLength uint32
	DataBuffer         *bytes.Buffer
}
