package logger

type NginxLogger struct {
	RemoteAddr           string  `json:"remote_addr"`
	HttpXForwardedFor    string  `json:"http_x_forwarded_for"`
	RemoteUser           string  `json:"remote_user"`
	TimeIso8601          string  `json:"created_at"`
	RequestMethod        string  `json:"method"`
	DocumentUri          string  `json:"document_uri"`
	QueryString          string  `json:"query_string"`
	ServerProtocol       string  `json:"server_protocol"`
	Status               string  `json:"status"`
	BodyBytesSent        int64   `json:"size"`
	RequestTime          float64 `json:"request_time"`
	HttpReferer          string  `json:"http_referer"`
	HttpUserAgent        string  `json:"user_agent"`
	HttpCdnSrcIp         string  `json:"http_Cdn-Src-Ip"`
	Host                 string  `json:"host"`
	HostName             string  `json:"hostname"`
	ServerAddr           string  `json:"server_addr"`
	RemotePort           string  `json:"remote_port"`
	ServerPort           string  `json:"server_port"`
	UpstreamAddr         string  `json:"upstream_host"`
	UpstreamStatus       string  `json:"upstream_status"`
	UpstreamResponseTime string  `json:"upstream_response_time"`
	V5                   string  `json:"V5"`
}

type NginxLoggerJson struct {
	CreatedAt            string  `json:"created_at"`
	RemoteAddr           string  `json:"remote_addr"`
	Method               string  `json:"method"`
	Request              string  `json:"request"`
	Status               string  `json:"status"`
	Size                 int64   `json:"size"`
	Referer              string  `json:"referer"`
	HttpHost             string  `json:"http_host"`
	UpstreamResponseTime string  `json:"upstream_response_time"`
	RequestTime          float64 `json:"request_time"`
	HttpXForwardedFor    string  `json:"http_x_forwarded_for"`
	UserAgent            string  `json:"user_agent"`
	UpstreamHost         string  `json:"upstream_host"`
	Host                 string  `json:"-"`
	DocumentUri          string  `json:"document_uri"`
}
