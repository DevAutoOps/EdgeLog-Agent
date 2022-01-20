package parser

import (
	"agent/common"
	"agent/global/variable"
	"bytes"
	"fmt"
)

type CommandParser struct {
	IsStart bool
}

func (c *CommandParser) Start() {
	if c.IsStart {
		return
	}
	c.IsStart = true

	for {
		if !c.IsStart {
			break
		}

		commandBufferClient := <-variable.CommandBufferChan
		common.Output(fmt.Sprintf("[CommandParser Start] Command data received: %d, data length: %d", commandBufferClient.Length, commandBufferClient.DataBuffer.Len()))
		if commandBufferClient.DataBuffer.Len() > 0 && commandBufferClient.DataBuffer.Len() == int(commandBufferClient.Length) {
			data := make([]byte, commandBufferClient.DataBuffer.Len())
			n, err := commandBufferClient.DataBuffer.Read(data)
			if n <= 0 || err != nil {
				continue
			}
			dataBuffer := new(bytes.Buffer)
			_, err = dataBuffer.Write(data)
			if err != nil {
				continue
			}

		}
	}
}
