package encode

import (
	"agent/common"
	"agent/global/consts"
	"agent/global/variable"
	"bytes"
	"encoding/binary"
	"fmt"
)

func EncodeLog(data []byte, oriLength uint32) (*bytes.Buffer, error) {
	var pkg = new(bytes.Buffer)
	// Write prefix
	err := binary.Write(pkg, binary.LittleEndian, []byte(consts.LogPrefix))
	// Read the length of the message and convert it to uint32 type (4 bytes)
	var length = uint32(len(data))
	// Write header
	err = binary.Write(pkg, binary.LittleEndian, length)
	if err != nil {
		return nil, err
	}
	// Write raw data length
	err = binary.Write(pkg, binary.LittleEndian, oriLength)
	if err != nil {
		return nil, err
	}
	if variable.Config.Logs.PrintLevel == 1 {
		common.Output(fmt.Sprintf("[EncodeLog] encode log data length: %d, oriLength: %d", length, oriLength))
	}
	// Write message body
	err = binary.Write(pkg, binary.LittleEndian, data)
	if err != nil {
		return nil, err
	}
	return pkg, nil
}

func EncodeLog2(data []byte) (*bytes.Buffer, error) {
	var pkg = new(bytes.Buffer)
	// Write prefix
	err := binary.Write(pkg, binary.LittleEndian, []byte(consts.LogPrefix))
	// Read the length of the message and convert it to uint32 type (4 bytes)
	var length = uint32(len(data))
	// Write header
	err = binary.Write(pkg, binary.LittleEndian, length)
	if err != nil {
		return nil, err
	}
	if variable.Config.Logs.PrintLevel == 1 {
		common.Output(fmt.Sprintf("[EncodeLog2] encode log data length: %d", length))
	}
	// Write message body
	err = binary.Write(pkg, binary.LittleEndian, data)
	if err != nil {
		return nil, err
	}
	return pkg, nil
}
