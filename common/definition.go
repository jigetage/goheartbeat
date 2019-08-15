package common

type Msg struct {
	CliName string `json:"cli_name"`
	Info    string `json:"info"`
}

var (
	Timeout int64 = 5
	RecvBuf int64 = 128
)
