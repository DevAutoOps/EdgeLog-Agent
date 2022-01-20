package logger

// LogData Log Data
type LogData struct {
	Index        string  // index
	CreatedAt    string  // Access time
	RemoteAddr   string  // Access IP address
	RequestTime  float64 // Response time
	Status       string  // Status code
	DocumentUri  string  // Page view
	Size         int64   // Website traffic
	Method       string  // Request mode
	UpstreamHost string
	NodeType     uint8  // Node type
	NodeId       uint32 // Node ID
	HttpHost     string // Site
}
