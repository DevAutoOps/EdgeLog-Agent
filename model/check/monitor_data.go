package check

type MonitorData struct {
	CPU
	Mem
	Disk
	NginxStatistics
	Network
}

type CPU struct {
	CPUUsage     int32 // CPU utilization
	CPUBaseUsage int32 // Basic CPU utilization
	CPULoad      int32 // Average CPU load
}

type Mem struct {
	MemUsed int32 // Memory usage
	MemRate int32 // Memory utilization
}

type Disk struct {
	DiskPartRate  int32 // Hard disk partition utilization
	DiskReadFlow  int32 // Hard disk read traffic(KB/s)
	DiskWriteFlow int32 // Hard disk write traffic(KB/s)
	DiskReadIOPS  int32 // Hard disk read IOPs (PCs.)
	DiskWriteIOPS int32 // Hard disk write IOPs (PCs.)
	DiskIOWait    int32 // Hard disk IO wait time(ms)
	DiskIOService int32 // Hard disk IO service time(ms)
	DiskIOBusy    int32 // Hard disk IO busy ratio(%)
	DiskPartTotal int32 // Total capacity of hard disk partition(GB)
	DiskPartUsage int32 // Used amount of hard disk partition(GB)
}

type NginxStatistics struct {
	NginxConcurrency int32 // Single nginx concurrency
}

type Network struct {
	NetworkReception int32
	NetworkSending   int32
}
