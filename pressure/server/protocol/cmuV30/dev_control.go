package cmuV30

const (
	CmdDevControl = "DevControl"
)

type DevControlWrapper struct {
	Cmd     string                 `json:"CMD"`
	Imb     Imb                    `json:"IMB"`
	PwdRst  int                    `json:"PwdRst"`
	PwdPort map[string]ControlUnit `json:"PwdPort"`
	Heat    ControlUnit            `json:"Heat"`
	Switch  map[string]ControlUnit `json:"Switch"`
}

type Imb struct {
	Fan   ControlUnit `json:"Fan"`
	Lamp  ControlUnit `json:"Lamp"`
	Guard ControlUnit `json:"Guard"`
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
