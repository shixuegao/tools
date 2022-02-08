package cmuV30

const (
	CmdGetSwitchConfig = "GetSwitchConfig"
	CmdPutSwitchConfig = "PutSwitchConfig"
)

type SwitchV30 struct {
	Cmd    string    `json:"CMD"`
	Switch SubSwitch `json:"Switch"`
	IPC    []IPC     `json:"IPC"`
	IMB    IMB       `json:"IMB"`
}

type SubSwitch struct {
	IP         string `json:"IP"`
	UpLinkPort int    `json:"UpLinkPort"`
	PowOutCh   int    `json:"PowOutCh"`
	Supplier   int    `json:"Supplier"`
	Model      string `json:"Model"`
}

type IPC struct {
	IP         string `json:"IP"`
	SwitchPort int    `json:"SwitchPort"`
	PowOutCh   int    `json:"PowOutCh"`
}

type IMB struct {
	SwitchPort int `json:"SwitchPort"`
}

type SwitchResponse struct {
	Cmd        string `json:"CMD"`
	No         string `json:"No"`
	ResultCode int    `json:"resultCode"`
}

func NewSwitchV30() *SwitchV30 {
	return &SwitchV30{
		Cmd: CmdGetSwitchConfig,
		Switch: SubSwitch{
			IP:         "127.0.0.1",
			UpLinkPort: 0,
			PowOutCh:   0,
			Supplier:   0,
			Model:      "xxxx",
		},
		IPC: make([]IPC, 0),
		IMB: IMB{
			SwitchPort: 0,
		},
	}
}
