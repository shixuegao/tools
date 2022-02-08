package cmu

const (
	CmdDevControl = "DevControl"
)

type DevControlWrapper struct {
	Cmd     string      `json:"CMD"`
	Imb     Imb         `json:"IMB"`
	PwdRst  int         `json:"PwdRst"`
	PwdPort PwdPort     `json:"PwdPort"`
	Heat    ControlUnit `json:"Heat"`
}

type Imb struct {
	Fan       ControlUnit `json:"Fan"`
	Buzzer    ControlUnit `json:"Buzzer"`
	Indicator ControlUnit `json:"Indicator"`
	Lamp      ControlUnit `json:"Lamp"`
	Guard     ControlUnit `json:"Guard"`
}

type PwdPort struct {
	Ac1  ControlUnit `json:"AC1"`
	Ac2  ControlUnit `json:"AC2"`
	Ac3  ControlUnit `json:"AC3"`
	Ac4  ControlUnit `json:"AC4"`
	Ac5  ControlUnit `json:"AC5"`
	Ac6  ControlUnit `json:"AC6"`
	Ac7  ControlUnit `json:"AC7"`
	Ac8  ControlUnit `json:"AC8"`
	Ipc1 ControlUnit `json:"IPC1"`
	Ipc2 ControlUnit `json:"IPC2"`
	Ipc3 ControlUnit `json:"IPC3"`
	Ipc4 ControlUnit `json:"IPC4"`
	Ipc5 ControlUnit `json:"IPC5"`
	Ipc6 ControlUnit `json:"IPC6"`
}

type ControlUnit struct {
	Remote string `json:"Remote"`
	Open   int    `json:"Open"`
}

type DevControlResponse struct {
	Cmd        string `json:"CMD"`
	No         string `json:"No"`
	ResultCode int    `json:"resultCode"`
}
