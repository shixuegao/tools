package cmuV30

import "math/rand"

const (
	CmdPostKeepLive = "PostKeepLive"
)

type KeepLiveWrapper struct {
	Cmd  string `json:"CMD"`
	No   string `json:"No"`
	Sn   string `json:"SN"`
	IP   string `json:"IP"`
	Type int    `json:"Type"`
}

//NewKeepLiveWrapper 生成心跳包
func NewKeepLiveWrapper(number string) *KeepLiveWrapper {
	t := rand.Intn(2)
	return &KeepLiveWrapper{
		Cmd:  CmdPostKeepLive,
		No:   number,
		Sn:   "1234567890",
		IP:   "127.0.0.1",
		Type: t,
	}
}
