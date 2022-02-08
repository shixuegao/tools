package cmu

import "math/rand"

const (
	CmdPostKeepLive = "PostKeepLive"
)

type KeepLiveWrapper struct {
	Cmd  string `json:"CMD"`
	No   string `json:"No"`
	Sn   string `json:"SN"`
	Type int    `json:"Type"`
}

//NewKeepLiveWrapper 生成心跳包
func NewKeepLiveWrapper(number string) *KeepLiveWrapper {
	t := rand.Intn(2)
	return &KeepLiveWrapper{
		Cmd:  CmdPostKeepLive,
		No:   number,
		Sn:   "1234567890",
		Type: t,
	}
}
