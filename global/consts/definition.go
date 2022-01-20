package consts

const (
	OsWindows string = "windows"
	OsLinux   string = "linux"
	OsMac     string = "darwin"
	OsOther   string = "other"
)

const (
	DefaultFreq   = 5000
	LogPrefix     = "taq"
	MonitorPrefix = "pbn"
)

const (
	MonitorData uint16 = iota
)

const (
	AppNginx = "nginx"
)

const (
	MonitorCpuUsage         uint8 = iota // CPU utilization
	MonitorCpuBaseUsage                  // Basic CPU utilization
	MonitorCpuLoad                       // Average CPU load
	MonitorMemUsed                       // Memory usage
	MonitorMemRate                       // Memory utilization
	MonitorDiskPartRate                  // Hard disk partition utilization
	MonitorDiskReadFlow                  // Hard disk read traffic(KB/s)
	MonitorDiskWriteFlow                 // Hard disk write traffic(KB/s)
	MonitorDiskIoWait                    // Hard disk IO wait time(ms)
	MonitorDiskIoService                 // Hard disk IO service time(ms)
	MonitorDiskIoBusy                    // Hard disk IO busy ratio(%)
	MonitorDiskPartTotal                 // Total capacity of hard disk partition(GB)
	MonitorDiskPartUsage                 // Used amount of hard disk partition(GB)
	MonitorNginxConcurrency              // Single nginx concurrency
	MonitorNetworkReception              // Network Reception
	MonitorNetworkSending                // Network sending
)

const (
	MonitorNginx uint8 = 128 + iota
)
