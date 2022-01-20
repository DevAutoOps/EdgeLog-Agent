package model

import (
	"bytes"
	"net"
)

type DataBufferClient struct {
	DataBuffer *bytes.Buffer
	Length     uint32
	Conn       *net.Conn
}
