package check

import (
	"agent/common"
	"agent/global/consts"
	"agent/global/journal"
	"agent/global/variable"
	"agent/model/check"
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"math"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type ServerCheck struct {
}

func (s *ServerCheck) Check() {
	if checkedResult, ok := variable.Result.(*check.CheckedResult); ok {
		if cpuRateValue, cpuBaseRateValue := s.cpuRate(); cpuRateValue > -1 || cpuBaseRateValue > -1 {
			if cpuRateValue > -1 {
				checkedResult.Lock()
				checkedResult.CheckList[consts.MonitorCpuUsage] = true
				checkedResult.Unlock()
				checkedResult.CPUUsage = cpuRateValue
			}
			if cpuBaseRateValue > -1 {
				checkedResult.Lock()
				checkedResult.CheckList[consts.MonitorCpuBaseUsage] = true
				checkedResult.Unlock()
				checkedResult.CPUBaseUsage = cpuBaseRateValue
			}
		} else {
			checkedResult.Lock()
			checkedResult.CheckList[consts.MonitorCpuUsage] = false
			checkedResult.CheckList[consts.MonitorCpuBaseUsage] = false
			checkedResult.Unlock()
		}
		if cpuLoadValue := s.cpuLoad(); cpuLoadValue > -1 {
			checkedResult.Lock()
			checkedResult.CheckList[consts.MonitorCpuLoad] = true
			checkedResult.Unlock()
			checkedResult.CPULoad = cpuLoadValue
		} else {
			checkedResult.Lock()
			checkedResult.CheckList[consts.MonitorCpuLoad] = false
			checkedResult.Unlock()
		}
		if memUsedValue, memRateValue := s.memoryUsage(); memUsedValue > -1 && memRateValue > -1 {
			checkedResult.Lock()
			checkedResult.CheckList[consts.MonitorMemUsed] = true
			checkedResult.CheckList[consts.MonitorMemRate] = true
			checkedResult.Unlock()
			checkedResult.MemUsed = memUsedValue
			checkedResult.MemRate = memRateValue
		} else {
			checkedResult.Lock()
			checkedResult.CheckList[consts.MonitorMemUsed] = false
			checkedResult.CheckList[consts.MonitorMemRate] = false
			checkedResult.Unlock()
		}
		if diskReadFlowValue, diskWriteFlowValue := s.diskReadWrite(); diskReadFlowValue > -1 && diskWriteFlowValue > -1 {
			checkedResult.Lock()
			checkedResult.CheckList[consts.MonitorDiskReadFlow] = true
			checkedResult.CheckList[consts.MonitorDiskWriteFlow] = true
			checkedResult.Unlock()
			checkedResult.DiskReadFlow = diskReadFlowValue
			checkedResult.DiskWriteFlow = diskWriteFlowValue
		} else {
			checkedResult.Lock()
			checkedResult.CheckList[consts.MonitorDiskReadFlow] = false
			checkedResult.CheckList[consts.MonitorDiskWriteFlow] = false
			checkedResult.Unlock()
		}
		if diskPartTotalValue, diskPartUsageValue, diskPartRateValue := s.partitionUtilizationRate(); diskPartTotalValue > -1 && diskPartUsageValue > -1 && diskPartRateValue > -1 {
			checkedResult.Lock()
			checkedResult.CheckList[consts.MonitorDiskPartTotal] = true
			checkedResult.CheckList[consts.MonitorDiskPartUsage] = true
			checkedResult.CheckList[consts.MonitorDiskPartRate] = true
			checkedResult.Unlock()
			checkedResult.DiskPartTotal = diskPartTotalValue
			checkedResult.DiskPartUsage = diskPartUsageValue
			checkedResult.DiskPartRate = diskPartRateValue
		} else {
			checkedResult.Lock()
			checkedResult.CheckList[consts.MonitorDiskPartTotal] = false
			checkedResult.CheckList[consts.MonitorDiskPartUsage] = false
			checkedResult.CheckList[consts.MonitorDiskPartRate] = false
			checkedResult.Unlock()
		}
		if diskIOWaitValue, diskIOServiceValue, diskIOBusyValue := s.ioBusy(); diskIOWaitValue > -1 && diskIOServiceValue > -1 && diskIOBusyValue > -1 {
			checkedResult.Lock()
			checkedResult.CheckList[consts.MonitorDiskIoWait] = true
			checkedResult.CheckList[consts.MonitorDiskIoService] = true
			checkedResult.CheckList[consts.MonitorDiskIoBusy] = true
			checkedResult.Unlock()
			checkedResult.DiskIOWait = diskIOWaitValue
			checkedResult.DiskIOService = diskIOServiceValue
			checkedResult.DiskIOBusy = diskIOBusyValue
		} else {
			checkedResult.Lock()
			checkedResult.CheckList[consts.MonitorDiskIoWait] = false
			checkedResult.CheckList[consts.MonitorDiskIoService] = false
			checkedResult.CheckList[consts.MonitorDiskIoBusy] = false
			checkedResult.Unlock()
		}
		if nginxConcurrencyValue := s.appConcurrency(consts.AppNginx); nginxConcurrencyValue > -1 {
			checkedResult.Lock()
			checkedResult.CheckList[consts.MonitorNginxConcurrency] = true
			checkedResult.Unlock()
			checkedResult.NginxConcurrency = nginxConcurrencyValue
		} else {
			checkedResult.Lock()
			checkedResult.CheckList[consts.MonitorNginxConcurrency] = false
			checkedResult.Unlock()
		}
		if networkReceptionValue, networkSendingValue := s.network(); networkReceptionValue > -1 && networkSendingValue > -1 {
			checkedResult.Lock()
			checkedResult.CheckList[consts.MonitorNetworkReception] = true
			checkedResult.CheckList[consts.MonitorNetworkSending] = true
			checkedResult.Unlock()
			checkedResult.NetworkReception = networkReceptionValue
			checkedResult.NetworkSending = networkSendingValue
		} else {
			checkedResult.Lock()
			checkedResult.CheckList[consts.MonitorNetworkReception] = false
			checkedResult.CheckList[consts.MonitorNetworkSending] = false
			checkedResult.Unlock()
		}
		if variable.CmdErrCounter%variable.CmdErrLimit == 0 {
			variable.CmdErrCounter = 0
		}
		variable.CmdErrCounter++
	} else {
		common.Output("[ServerCheck Check] get result object err", journal.Warn)
	}
}

// Get CPU utilization and basic utilization
func (s *ServerCheck) cpuRate() (int32, int32) {
	os := common.GetOs()
	if os == consts.OsLinux {
		return s.cpuRateForLinux()
	} else if os == consts.OsWindows {
		return s.cpuRateForWindows()
	}
	return -1, -1
}

func (s *ServerCheck) cpuRateForLinux() (int32, int32) {
	cmd := "top -bn1 | grep -i cpu"
	output, success, _ := common.Cmd(cmd)
	if len(output) <= 0 || !success {
		return -1, -1
	}
	result := strings.Fields(output)
	if len(result) <= 0 {
		return -1, -1
	}
	value1, err := strconv.ParseFloat(result[1], 64)
	if err != nil {
		return -1, -1
	}
	result1 := int32(math.Round(value1 * 100))
	value2, err := strconv.ParseFloat(result[3], 64)
	if err != nil {
		return -1, -1
	}
	result2 := int32(math.Round(value2 * 100))
	return result1, result2
}

func (s *ServerCheck) cpuRateForWindows() (int32, int32) {
	perCpu := false
	numCpu := runtime.NumCPU()
	testCount := 3
	cpuValue := 0.0
	isAssign := false
	isValue := false
	for i := 0; i < testCount; i++ {
		duration := time.Duration(1) * time.Microsecond
		v, err := cpu.Percent(duration, perCpu)
		if err != nil {
			//common.Output(fmt.Sprintf("cpuRateForWindows err: %s", err), journal.Warn)
			continue
		}
		for _, percent := range v {
			if percent >= 0.0 && percent <= 100.0001*float64(numCpu) {
				//fmt.Println("cpuTest:", percent)
				isValue = true
				if percent > 0 && !isAssign {
					cpuValue = percent
					isAssign = true
				}
			}
		}
	}
	if isAssign || isValue {
		return int32(math.Round(cpuValue * 100)), -1
	} else {
		return -1, -1
	}
}

// Get CPU load
func (s *ServerCheck) cpuLoad() int32 {
	os := common.GetOs()
	if os == consts.OsLinux {
		return s.cpuLoadForLinux()
	} else if os == consts.OsWindows {
		return s.cpuLoadForWindows()
	}
	return -1
}

func (s *ServerCheck) cpuLoadForLinux() int32 {
	cmd := "uptime"
	output, success, _ := common.Cmd(cmd)
	if len(output) <= 0 || !success {
		return -1
	}
	outputArr := strings.Fields(output)
	if len(outputArr) <= 0 {
		return -1
	}
	load := strings.TrimSuffix(outputArr[len(outputArr)-3], ",")
	value, _ := strconv.ParseFloat(load, 64)
	return int32(math.Round(value * 100))
}

func (s *ServerCheck) cpuLoadForWindows() int32 {
	cmd := "wmic cpu get loadpercentage"
	output, success, _ := common.CmdWindows(cmd)
	if len(output) <= 0 || !success {
		return -1
	}
	outputArr := strings.Fields(output)
	if len(outputArr) <= 1 {
		return -1
	}
	value, _ := strconv.ParseFloat(outputArr[1], 64)
	return int32(math.Round(value * 100))
}

// Memory usage and memory utilization
func (s *ServerCheck) memoryUsage() (int32, int32) {
	os := common.GetOs()
	if os == consts.OsLinux {
		return s.memoryUsageForLinux()
	} else if os == consts.OsWindows {
		return s.memoryUsageForWindows()
	}
	return -1, -1
}

func (s *ServerCheck) memoryUsageForLinux() (int32, int32) {
	cmd := "free -m"
	output, success, _ := common.Cmd(cmd)
	if len(output) <= 0 || !success {
		return -1, -1
	}
	result := strings.Fields(output)
	amount, err := strconv.ParseFloat(result[8], 64) // Memory usage
	if err != nil {
		return -1, -1
	}
	total, err := strconv.ParseFloat(result[7], 64)
	if err != nil {
		return -1, -1
	}
	value, err := strconv.ParseFloat(fmt.Sprintf("%.2f", amount/total), 64) // memory utilization
	if err != nil {
		return -1, -1
	}
	return int32(math.Round(amount)), int32(math.Round(value * 100))
}

func (s *ServerCheck) memoryUsageForWindows() (int32, int32) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return -1, -1
	}
	empty := &mem.VirtualMemoryStat{}
	if v == empty {
		return -1, -1
	}
	used := math.Round(float64(v.Used) / 1024 / 1024)
	return int32(used), int32(math.Round(v.UsedPercent * 100))
}

// Hard disk read and write traffic
func (s *ServerCheck) diskReadWrite() (int32, int32) {
	os := common.GetOs()
	if os == consts.OsWindows {
		return -1, -1
	}
	ioCmd := "./sysstat/iostat -t"
	ioOutput, success, _ := common.Cmd(ioCmd)
	if len(ioOutput) <= 0 || !success {
		return -1, -1
	}
	dfCmd := "df | grep \"/dev/\""
	dfOutput, success, _ := common.Cmd(dfCmd)
	if len(dfOutput) <= 0 || !success {
		return -1, -1
	}

	readFieldIndex := 0
	writeFieldIndex := 0
	ioLines := strings.Split(ioOutput, "\n")
	ioFilterLines := make([]string, 0)
	startFilter := false
	for _, ioLine := range ioLines {
		currentLine := strings.TrimSpace(ioLine)
		if startFilter && len(currentLine) > 0 {
			ioFilterLines = append(ioFilterLines, currentLine)
		}
		if strings.HasPrefix(strings.ToLower(currentLine), "device") {
			fields := strings.Fields(currentLine)
			for fieldNum, fieldTitle := range fields {
				switch fieldTitle {
				case "kB_read/s":
					readFieldIndex = fieldNum
				case "kB_wrtn/s":
					writeFieldIndex = fieldNum
				}
			}
			startFilter = true
		}
	}

	if readFieldIndex <= 0 && writeFieldIndex <= 0 {
		return 0, 0
	}
	readValues := make([]float64, 0)
	writeValues := make([]float64, 0)
	dfLines := strings.Split(dfOutput, "\n")
	for _, ioLine := range ioFilterLines {
		ioFields := strings.Fields(strings.TrimSpace(ioLine))
		if len(ioFields) < readFieldIndex || len(ioFields) < writeFieldIndex {
			continue
		}
		for _, dfLine := range dfLines {
			currentDfLine := strings.TrimSpace(dfLine)
			if !strings.HasPrefix(strings.ToLower(currentDfLine), "/dev/") {
				continue
			}
			if len(currentDfLine) <= 0 {
				continue
			}
			dfFields := strings.Fields(currentDfLine)
			if strings.Contains(dfFields[0], ioFields[0]) {
				if readFieldIndex > 0 {
					readValue, err := strconv.ParseFloat(ioFields[readFieldIndex], 64)
					if err == nil {
						readValues = append(readValues, readValue)
					}
				}
				if writeFieldIndex > 0 {
					writeValue, err := strconv.ParseFloat(ioFields[writeFieldIndex], 64)
					if err == nil {
						writeValues = append(writeValues, writeValue)
					}
				}
			}
		}
	}

	if len(readValues) <= 0 && len(writeValues) <= 0 {
		return 0, 0
	}
	var readAverage float64 = 0
	var writeAverage float64 = 0
	if len(readValues) > 0 {
		var readTotalValue float64 = 0
		for _, value := range readValues {
			readTotalValue += value
		}
		readAverage = readTotalValue / float64(len(readValues))
	}
	if len(writeValues) > 0 {
		var writeTotalValue float64 = 0
		for _, value := range writeValues {
			writeTotalValue += value
		}
		writeAverage = writeTotalValue / float64(len(writeValues))
	}
	var readResult int32 = 0
	var writeResult int32 = 0
	if readAverage > 0 {
		readResult = int32(math.Round(readAverage * 100))
	}
	if writeAverage > 0 {
		writeResult = int32(math.Round(writeAverage * 100))
	}
	return readResult, writeResult
}

// Total capacity, used capacity and utilization of hard disk partition
func (s *ServerCheck) partitionUtilizationRate() (int32, int32, int32) {
	os := common.GetOs()
	if os == consts.OsLinux {
		return s.partitionUtilizationRateForLinux()
	} else if os == consts.OsWindows {
		return s.partitionUtilizationRateForWindows()
	}
	return -1, -1, -1
}

func (s *ServerCheck) partitionUtilizationRateForLinux() (int32, int32, int32) {
	cmd := "df -h / | tail -n 1"
	output, success, _ := common.Cmd(cmd)
	if len(output) <= 0 || !success {
		return -1, -1, -1
	}
	result := strings.Fields(output)
	Used, err := strconv.ParseFloat(result[2][:len(result[2])-1], 64) // Used
	if err != nil {
		return -1, -1, -1
	}
	total, err := strconv.ParseFloat(result[1][:len(result[1])-1], 64) // Total
	if err != nil {
		return -1, -1, -1
	}
	percentage, err := strconv.ParseFloat(result[4][:len(result[4])-1], 64) // Utilization rate
	if err != nil {
		return -1, -1, -1
	}
	return int32(math.Round(total * 100)), int32(math.Round(Used * 100)), int32(math.Round(percentage * 100))
}

func (s *ServerCheck) partitionUtilizationRateForWindows() (int32, int32, int32) {
	ret, err := disk.Partitions(false)
	if err != nil {
		return -1, -1, -1
	}
	var total uint64 = 0
	var used uint64 = 0
	for _, part := range ret {
		v, err := disk.Usage(part.Device)
		if err != nil {
			//fmt.Println("diskTest err:", err)
			continue
		}
		//fmt.Println(v)
		total += v.Total
		used += v.Used
	}
	usedPercent := float64(used) / float64(total)
	return int32(math.Round(float64(total) / 1024 / 1024 / 1024 * 100)), int32(math.Round(float64(used) / 1024 / 1024 / 1024 * 100)), int32(math.Round(usedPercent * 100))
}

// IO waiting time, service time, busy ratio
func (s *ServerCheck) ioBusy() (int32, int32, int32) {
	os := common.GetOs()
	if os == consts.OsWindows {
		return -1, -1, -1
	}
	ioCmd := "./sysstat/iostat -x"
	ioOutput, success, _ := common.Cmd(ioCmd)
	if len(ioOutput) <= 0 || !success {
		return -1, -1, -1
	}
	dfCmd := "df | grep \"/dev/\""
	dfOutput, success, _ := common.Cmd(dfCmd)
	if len(dfOutput) <= 0 || !success {
		return -1, -1, -1
	}

	awaitFieldIndex := 0
	svctmFieldIndex := 0
	utilFieldIndex := 0
	ioLines := strings.Split(ioOutput, "\n")
	ioFilterLines := make([]string, 0)
	startFilter := false
	for _, ioLine := range ioLines {
		currentLine := strings.TrimSpace(ioLine)
		if startFilter && len(currentLine) > 0 {
			ioFilterLines = append(ioFilterLines, currentLine)
		}
		if strings.HasPrefix(strings.ToLower(currentLine), "device") {
			fields := strings.Fields(currentLine)
			for fieldNum, fieldTitle := range fields {
				switch fieldTitle {
				case "await":
					awaitFieldIndex = fieldNum
				case "svctm":
					svctmFieldIndex = fieldNum
				case "%util":
					utilFieldIndex = fieldNum
				}
			}
			startFilter = true
		}
	}

	if awaitFieldIndex <= 0 && svctmFieldIndex <= 0 && utilFieldIndex <= 0 {
		return 0, 0, 0
	}
	awaitValues := make([]float64, 0)
	svctmValues := make([]float64, 0)
	utilValues := make([]float64, 0)
	dfLines := strings.Split(dfOutput, "\n")
	for _, ioLine := range ioFilterLines {
		ioFields := strings.Fields(strings.TrimSpace(ioLine))
		if len(ioFields) < awaitFieldIndex || len(ioFields) < svctmFieldIndex || len(ioFields) < utilFieldIndex {
			continue
		}
		for _, dfLine := range dfLines {
			currentDfLine := strings.TrimSpace(dfLine)
			if !strings.HasPrefix(strings.ToLower(currentDfLine), "/dev/") {
				continue
			}
			if len(currentDfLine) <= 0 {
				continue
			}
			dfFields := strings.Fields(currentDfLine)
			if strings.Contains(dfFields[0], ioFields[0]) {
				if awaitFieldIndex > 0 {
					awaitValue, err := strconv.ParseFloat(ioFields[awaitFieldIndex], 64)
					if err == nil {
						awaitValues = append(awaitValues, awaitValue)
					}
				}
				if svctmFieldIndex > 0 {
					svctmValue, err := strconv.ParseFloat(ioFields[svctmFieldIndex], 64)
					if err == nil {
						svctmValues = append(svctmValues, svctmValue)
					}
				}
				if utilFieldIndex > 0 {
					utilValue, err := strconv.ParseFloat(ioFields[utilFieldIndex], 64)
					if err == nil {
						utilValues = append(utilValues, utilValue)
					}
				}
			}
		}
	}

	if len(awaitValues) <= 0 && len(svctmValues) <= 0 && len(utilValues) <= 0 {
		return 0, 0, 0
	}
	var awaitAverage float64 = 0
	var svctmAverage float64 = 0
	var utilAverage float64 = 0
	if len(awaitValues) > 0 {
		var awaitTotalValue float64 = 0
		for _, value := range awaitValues {
			awaitTotalValue += value
		}
		awaitAverage = awaitTotalValue / float64(len(awaitValues))
	}
	if len(svctmValues) > 0 {
		var svctmTotalValue float64 = 0
		for _, value := range svctmValues {
			svctmTotalValue += value
		}
		svctmAverage = svctmTotalValue / float64(len(svctmValues))
	}
	if len(utilValues) > 0 {
		var utilTotalValue float64 = 0
		for _, value := range utilValues {
			utilTotalValue += value
		}
		utilAverage = utilTotalValue / float64(len(utilValues))
	}
	var awaitResult int32 = 0
	var svctmResult int32 = 0
	var utilResult int32 = 0
	if awaitAverage > 0 {
		awaitResult = int32(math.Round(awaitAverage * 100))
	}
	if svctmAverage > 0 {
		svctmResult = int32(math.Round(svctmAverage * 100))
	}
	if utilAverage > 0 {
		utilResult = int32(math.Round(utilAverage * 10000))
	}
	return awaitResult, svctmResult, utilResult
}

// appConcurrency Concurrency per app
func (s *ServerCheck) appConcurrency(app string) int32 {
	os := common.GetOs()
	if os == consts.OsLinux {
		return s.appConcurrencyForLinux(app)
	} else if os == consts.OsWindows {
		return s.appConcurrencyForWindows(app)
	}
	return -1
}

func (s *ServerCheck) appConcurrencyForLinux(app string) int32 {
	var result int32 = 0
	cmd := fmt.Sprintf("netstat -apn | grep %s | grep tcp", app)
	output, success, _ := common.Cmd(cmd)
	if len(output) <= 0 || !success {
		if variable.CmdErrCounter%variable.CmdErrLimit == 0 {
			common.Output(fmt.Sprintf("[ServerCheck NginxConcurrency] cmd (%s) err", cmd), journal.Warn)
		}
		return -1
	}
	lines := strings.Split(output, "\n")
	ports := make([]string, 0)
	for _, line := range lines {
		lineArr := strings.Fields(line)
		if len(lineArr) >= 4 {
			hostArr := strings.Split(lineArr[3], ":")
			if len(hostArr) == 2 {
				ports = append(ports, hostArr[1])
			}
		}
	}
	if len(ports) > 0 {
		for _, port := range ports {
			cmd2 := fmt.Sprintf("netstat -n | grep :%s | grep -v grep | awk '/^tcp/ {++S[$NF]} END {for(a in S) print a, S[a]}'", port)
			output2, success2, _ := common.Cmd(cmd2)
			if len(output2) <= 0 || !success2 {
				if variable.CmdErrCounter%variable.CmdErrLimit == 0 {
					common.Output(fmt.Sprintf("[ServerCheck NginxConcurrency] cmd (%s) err", cmd2), journal.Warn)
				}
				continue
			}
			lines2 := strings.Split(output2, "\n")
			for _, line2 := range lines2 {
				lineArr2 := strings.Fields(line2)
				if len(lineArr2) == 2 && strings.ToUpper(lineArr2[0]) == "ESTABLISHED" {
					num, err := strconv.Atoi(lineArr2[1])
					if err != nil {
						common.Output(fmt.Sprintf("[ServerCheck NginxConcurrency] err for convert to digital: %s", err), journal.Warn)
						continue
					}
					result += int32(num)
				}
			}
		}
	}
	return result
}

func (s *ServerCheck) appConcurrencyForWindows(app string) int32 {
	app = app + ".exe"
	var result int32 = 0
	cmdTask := fmt.Sprintf("tasklist | findstr /I %s", app)
	outputTask, success, _ := common.CmdWindows(cmdTask)
	if len(outputTask) <= 0 || !success {
		return -1
	}
	taskLines := strings.Split(outputTask, "\r\n")
	pidList := make([]string, 0)
	for _, taskLine := range taskLines {
		currentLine := strings.TrimSpace(taskLine)
		if len(currentLine) <= 0 {
			continue
		}
		lineArr := strings.Fields(currentLine)
		if len(lineArr) >= 2 {
			pidList = append(pidList, lineArr[1])
		}
	}
	if len(pidList) <= 0 {
		return -1
	}
	for _, pid := range pidList {
		if len(pid) <= 0 {
			continue
		}
		cmdNet := fmt.Sprintf("netstat -ano | findstr /I TCP | findstr %s | findstr /I ESTABLISHED", pid)
		outputNet, success, _ := common.CmdWindows(cmdNet)
		if len(outputNet) <= 0 || !success {
			if variable.CmdErrCounter%variable.CmdErrLimit == 0 {
				common.Output(fmt.Sprintf("[ServerCheck NginxConcurrency] cmd (%s) err", cmdNet), journal.Warn)
			}
			continue
		}
		netLines := strings.Split(outputNet, "\r\n")
		for _, netLine := range netLines {
			currentLine := strings.TrimSpace(netLine)
			if len(currentLine) <= 0 {
				continue
			}
			lineArr := strings.Fields(currentLine)
			if len(lineArr) >= 5 && lineArr[4] == pid {
				result++
			}
		}
	}
	return result
}

// net
func (s *ServerCheck) network() (int32, int32) {
	if variable.LastMonitorData == nil {
		variable.LastMonitorData = &check.MonitorRecord{
			NetworkReception: 0,
			NetworkSending:   0,
		}
	}
	netIo, err := net.IOCounters(false)
	if err != nil {
		common.Output(fmt.Sprintf("[ServerCheck network] net IOCounters err: %s", err), journal.Warn)
		return -1, -1
	}
	if len(netIo) <= 0 {
		common.Output("[ServerCheck network] net IOCounters is empty", journal.Warn)
		return -1, -1
	}
	if variable.LastMonitorData.NetworkReception == 0 && variable.LastMonitorData.NetworkSending == 0 {
		variable.LastMonitorData.NetworkReception = netIo[0].BytesRecv
		variable.LastMonitorData.NetworkSending = netIo[0].BytesSent
		return -1, -1
	} else {
		networkReception := math.Abs(float64(netIo[0].BytesRecv - variable.LastMonitorData.NetworkReception))
		networkSending := math.Abs(float64(netIo[0].BytesSent - variable.LastMonitorData.NetworkSending))
		variable.LastMonitorData.NetworkReception = netIo[0].BytesRecv
		variable.LastMonitorData.NetworkSending = netIo[0].BytesSent
		//common.Output(fmt.Sprintf("recv: %d, send: %d", int32(networkReception), int32(networkSending)))
		return int32(networkReception), int32(networkSending)
	}
}
