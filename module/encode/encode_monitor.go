package encode

import (
	"agent/common"
	"agent/global/consts"
	"agent/global/journal"
	"agent/global/variable"
	"agent/model/check"
	"bytes"
	"encoding/binary"
	"fmt"
)

func EncodeMonitor() (*bytes.Buffer, error) {
	var pkg = new(bytes.Buffer)
	var count uint8 = 0
	checkedResult, ok := variable.Result.(*check.CheckedResult)
	if !ok {
		common.Output("[EncodeMonitor2] get result object err", journal.Warn)
		return nil, nil
	}
	if variable.Config.Modules.Monitor.Check {
		checkedResult.RLock()
		for itemId, isCheck := range checkedResult.CheckList {
			if itemId >= consts.MonitorNginx {
				continue
			}
			if isCheck {
				switch itemId {
				case consts.MonitorCpuUsage:
					err := packageCPUUsage(pkg, checkedResult, itemId)
					if err != nil {
						common.Output(fmt.Sprintf("[EncodeMonitor2] write MonitorCpuUsage err: %s", err), journal.Error)
						checkedResult.RUnlock()
						return nil, err
					}
					count++
				case consts.MonitorCpuBaseUsage:
					err := packageCPUBaseUsage(pkg, checkedResult, itemId)
					if err != nil {
						common.Output(fmt.Sprintf("[EncodeMonitor2] write MonitorCpuBaseUsage err: %s", err), journal.Error)
						checkedResult.RUnlock()
						return nil, err
					}
					count++
				case consts.MonitorCpuLoad:
					err := packageCPULoad(pkg, checkedResult, itemId)
					if err != nil {
						common.Output(fmt.Sprintf("[EncodeMonitor2] write MonitorCpuLoad err: %s", err), journal.Error)
						checkedResult.RUnlock()
						return nil, err
					}
					count++
				case consts.MonitorMemUsed:
					err := packageMemUsed(pkg, checkedResult, itemId)
					if err != nil {
						common.Output(fmt.Sprintf("[EncodeMonitor2] write MonitorMemUsed err: %s", err), journal.Error)
						checkedResult.RUnlock()
						return nil, err
					}
					count++
				case consts.MonitorMemRate:
					err := packageMemRate(pkg, checkedResult, itemId)
					if err != nil {
						common.Output(fmt.Sprintf("[EncodeMonitor2] write MonitorMemRate err: %s", err), journal.Error)
						checkedResult.RUnlock()
						return nil, err
					}
					count++
				case consts.MonitorDiskPartRate:
					err := packageDiskPartPate(pkg, checkedResult, itemId)
					if err != nil {
						common.Output(fmt.Sprintf("[EncodeMonitor2] write MonitorDiskPartRate err: %s", err), journal.Error)
						checkedResult.RUnlock()
						return nil, err
					}
					count++
				case consts.MonitorDiskReadFlow:
					err := packageDiskReadFlow(pkg, checkedResult, itemId)
					if err != nil {
						common.Output(fmt.Sprintf("[EncodeMonitor2] write MonitorDiskReadFlow err: %s", err), journal.Error)
						checkedResult.RUnlock()
						return nil, err
					}
					count++
				case consts.MonitorDiskWriteFlow:
					err := packageDiskWriteFlow(pkg, checkedResult, itemId)
					if err != nil {
						common.Output(fmt.Sprintf("[EncodeMonitor2] write MonitorDiskWriteFlow err: %s", err), journal.Error)
						checkedResult.RUnlock()
						return nil, err
					}
					count++
				case consts.MonitorDiskIoWait:
					err := packageDiskIOWait(pkg, checkedResult, itemId)
					if err != nil {
						common.Output(fmt.Sprintf("[EncodeMonitor2] write MonitorDiskIoWait err: %s", err), journal.Error)
						checkedResult.RUnlock()
						return nil, err
					}
					count++
				case consts.MonitorDiskIoService:
					err := packageDiskIOService(pkg, checkedResult, itemId)
					if err != nil {
						common.Output(fmt.Sprintf("[EncodeMonitor2] write MonitorDiskIoService err: %s", err), journal.Error)
						checkedResult.RUnlock()
						return nil, err
					}
					count++
				case consts.MonitorDiskIoBusy:
					err := packageDiskIOBusy(pkg, checkedResult, itemId)
					if err != nil {
						common.Output(fmt.Sprintf("[EncodeMonitor2] write MonitorDiskIoBusy err: %s", err), journal.Error)
						checkedResult.RUnlock()
						return nil, err
					}
					count++
				case consts.MonitorDiskPartTotal:
					err := packageDiskPartTotal(pkg, checkedResult, itemId)
					if err != nil {
						common.Output(fmt.Sprintf("[EncodeMonitor2] write MonitorDiskPartTotal err: %s", err), journal.Error)
						checkedResult.RUnlock()
						return nil, err
					}
					count++
				case consts.MonitorDiskPartUsage:
					err := packageDiskPartUsage(pkg, checkedResult, itemId)
					if err != nil {
						common.Output(fmt.Sprintf("[EncodeMonitor2] write MonitorDiskPartUsage err: %s", err), journal.Error)
						checkedResult.RUnlock()
						return nil, err
					}
					count++
				case consts.MonitorNginxConcurrency:
					err := packageNginxConcurrency(pkg, checkedResult, itemId)
					if err != nil {
						common.Output(fmt.Sprintf("[EncodeMonitor2] write MonitorNginxConcurrency err: %s", err), journal.Error)
						checkedResult.RUnlock()
						return nil, err
					}
					count++
				case consts.MonitorNetworkReception:
					err := packageNetworkReception(pkg, checkedResult, itemId)
					if err != nil {
						common.Output(fmt.Sprintf("[EncodeMonitor2] write MonitorNetworkReception err: %s", err), journal.Error)
						checkedResult.RUnlock()
						return nil, err
					}
					count++
				case consts.MonitorNetworkSending:
					err := packageNetworkSending(pkg, checkedResult, itemId)
					if err != nil {
						common.Output(fmt.Sprintf("[EncodeMonitor2] write MonitorNetworkSending err: %s", err), journal.Error)
						checkedResult.RUnlock()
						return nil, err
					}
					count++
				}
			}
		}
		checkedResult.RUnlock()
	}
	checkedResult.RLock()
	for itemId, isCheck := range checkedResult.CheckList {
		if itemId < consts.MonitorNginx {
			continue
		}
		if isCheck {
			switch itemId {
			case consts.MonitorNginx:
				err := packageNginxStatus(pkg, checkedResult, itemId)
				if err != nil {
					common.Output(fmt.Sprintf("[EncodeMonitor2] write MonitorNginx err: %s", err), journal.Error)
				}
				count++
			}
		}
	}
	checkedResult.RUnlock()

	fullData, err := packageHead(pkg.Bytes(), count)
	if err != nil {
		common.Output(fmt.Sprintf("[EncodeMonitor2] package err: %s", err), journal.Error)
		return nil, err
	}
	return fullData, nil
}

func packageHead(data []byte, count uint8) (*bytes.Buffer, error) {
	var pkg = new(bytes.Buffer)
	// 写入前缀
	err := binary.Write(pkg, binary.LittleEndian, []byte(consts.MonitorPrefix))

	if len(data) <= 0 || count <= 0 {
		count = 0
	}
	err = binary.Write(pkg, binary.LittleEndian, count)
	if err != nil {
		common.Output(fmt.Sprintf("[EncodeMonitor2 packageHead] write err: %s", err), journal.Error)
		return nil, err
	}
	if count == 0 {
		return pkg, nil
	}

	var dataLength = uint32(len(data))
	err = binary.Write(pkg, binary.LittleEndian, dataLength)
	if err != nil {
		common.Output(fmt.Sprintf("[EncodeMonitor2 packageHead] write err: %s", err), journal.Error)
		return nil, err
	}
	err = binary.Write(pkg, binary.LittleEndian, data)
	if err != nil {
		common.Output(fmt.Sprintf("[EncodeMonitor2 packageHead] write err: %s", err), journal.Error)
		return nil, err
	}
	return pkg, nil
}

func packageMonitor(pkg *bytes.Buffer, intValue int32, sType uint8) error {
	err := binary.Write(pkg, binary.LittleEndian, sType)
	if err != nil {
		return err
	}
	err = binary.Write(pkg, binary.LittleEndian, intValue)
	if err != nil {
		return err
	}
	return nil
}

func packageCPUUsage(pkg *bytes.Buffer, checkedResult *check.CheckedResult, sType uint8) error {
	return packageMonitor(pkg, checkedResult.CPUUsage, sType)
}

func packageCPUBaseUsage(pkg *bytes.Buffer, checkedResult *check.CheckedResult, sType uint8) error {
	return packageMonitor(pkg, checkedResult.CPUBaseUsage, sType)
}

func packageCPULoad(pkg *bytes.Buffer, checkedResult *check.CheckedResult, sType uint8) error {
	return packageMonitor(pkg, checkedResult.CPULoad, sType)
}

func packageMemUsed(pkg *bytes.Buffer, checkedResult *check.CheckedResult, sType uint8) error {
	return packageMonitor(pkg, checkedResult.MemUsed, sType)
}

func packageMemRate(pkg *bytes.Buffer, checkedResult *check.CheckedResult, sType uint8) error {
	return packageMonitor(pkg, checkedResult.MemRate, sType)
}

func packageDiskPartPate(pkg *bytes.Buffer, checkedResult *check.CheckedResult, sType uint8) error {
	return packageMonitor(pkg, checkedResult.DiskPartRate, sType)
}

func packageDiskReadFlow(pkg *bytes.Buffer, checkedResult *check.CheckedResult, sType uint8) error {
	return packageMonitor(pkg, checkedResult.DiskReadFlow, sType)
}

func packageDiskWriteFlow(pkg *bytes.Buffer, checkedResult *check.CheckedResult, sType uint8) error {
	return packageMonitor(pkg, checkedResult.DiskWriteFlow, sType)
}

func packageDiskIOWait(pkg *bytes.Buffer, checkedResult *check.CheckedResult, sType uint8) error {
	return packageMonitor(pkg, checkedResult.DiskIOWait, sType)
}

func packageDiskIOService(pkg *bytes.Buffer, checkedResult *check.CheckedResult, sType uint8) error {
	return packageMonitor(pkg, checkedResult.DiskIOService, sType)
}

func packageDiskIOBusy(pkg *bytes.Buffer, checkedResult *check.CheckedResult, sType uint8) error {
	return packageMonitor(pkg, checkedResult.DiskIOBusy, sType)
}

func packageDiskPartTotal(pkg *bytes.Buffer, checkedResult *check.CheckedResult, sType uint8) error {
	return packageMonitor(pkg, checkedResult.DiskPartTotal, sType)
}

func packageDiskPartUsage(pkg *bytes.Buffer, checkedResult *check.CheckedResult, sType uint8) error {
	return packageMonitor(pkg, checkedResult.DiskPartUsage, sType)
}

func packageNginxConcurrency(pkg *bytes.Buffer, checkedResult *check.CheckedResult, sType uint8) error {
	return packageMonitor(pkg, checkedResult.NginxConcurrency, sType)
}

func packageNginxStatus(pkg *bytes.Buffer, checkedResult *check.CheckedResult, sType uint8) error {
	return packageMonitor(pkg, checkedResult.NginxStatus, sType)
}

func packageNetworkReception(pkg *bytes.Buffer, checkedResult *check.CheckedResult, sType uint8) error {
	return packageMonitor(pkg, checkedResult.NetworkReception, sType)
}

func packageNetworkSending(pkg *bytes.Buffer, checkedResult *check.CheckedResult, sType uint8) error {
	return packageMonitor(pkg, checkedResult.NetworkSending, sType)
}
