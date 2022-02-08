package cmuV30

const (
	CmdPutDevMsg = "PutDevMsg"
	CmdGetDevMsg = "GetDevMsg"
)

type DevStateWrapper struct {
	Cmd         string          `json:"CMD"`
	No          string          `json:"No"`
	Conect      bool            `json:"Conect"`
	Environ     Environ         `json:"Environ"`
	Door        Door            `json:"Door"`
	Spd         Spd             `json:"SPD"`
	Esw         Esw             `json:"eSW"`
	Index       Index           `json:"Index"`
	MaPowInSt   MaPowInSt       `json:"MaPowInSt"`
	PowModInfos []PowModInfo    `json:"PowModInfo"`
	Switch      Switch          `json:"Switch"`
	Wireless    WirelessOfState `json:"Wireless"`
	Time        string          `json:"Time"`
	Bias        Bias            `json:"Bias"`
	Pos         PosOfState      `json:"Pos"`
}

type MaPowInSt struct {
	Ac220 int `json:"AC220"`
	Dc12  int `json:"DC12"`
}

type PowModInfo struct {
	Ch      int `json:"Ch"`
	InOut   int `json:"InOut"`
	Status  int `json:"Status"`
	ACDC    int `json:"ACDC"`
	Current int `json:"Current"`
	Voltage int `json:"Voltage"`
}

type Environ struct {
	Humber int  `json:"Humber"`
	Temper int  `json:"Temper"`
	Influ  bool `json:"Influ"`
}

type Door struct {
	Guard bool `json:"Guard"`
	Open  bool `json:"Open"`
}

type Spd struct {
	Status int `json:"Status"`
	Times  int `json:"Times"`
}

type Esw struct {
	Com    int     `json:"com"`
	Status string  `json:"Status"`
	Ac220  int     `json:"AC220"`
	Irms   float32 `json:"Irms"`
	Leak   float32 `json:"Leak"`
	Oc     int     `json:"OC"`
	Ov     int     `json:"OV"`
}

type Index struct {
	BoLa   string `json:"BoLa"`
	Alarm  string `json:"Alarm"`
	Fan    string `json:"Fan"`
	FanSts int    `json:"FanSts"`
	Buz    string `json:"Buz"`
}

type WirelessOfState struct {
	Status int `json:"Status"`
	Csq    int `json:"CSQ"`
}

type Bias struct {
	Yaw   int `json:"Yaw"`
	Roll  int `json:"Roll"`
	Pitch int `json:"Pitch"`
}

type PosOfState struct {
	Longitude float32 `json:"Longitude"`
	Latitude  float32 `json:"Latitude"`
}

type Switch struct {
	ChanInfo []ChanInfo `json:"ChanInfo"`
}

type ChanInfo struct {
	Ch       int `json:"Ch"`
	Status   int `json:"Status"`
	InSpeed  int `json:"InSpeed"`
	OutSpeed int `json:"OutSpeed"`
}

func NewDevState(number string) *DevStateWrapper {
	return &DevStateWrapper{
		Cmd:    CmdPutDevMsg,
		No:     number,
		Conect: true,
		Environ: Environ{
			Humber: 11,
			Temper: 11,
			Influ:  false,
		},
		Door: Door{
			Guard: true,
			Open:  false,
		},
		Spd: Spd{
			Status: 0,
			Times:  11,
		},
		Esw: Esw{
			Com:    0,
			Status: "ON",
			Ac220:  220,
			Irms:   1.2,
			Leak:   1.3,
			Oc:     0,
			Ov:     0,
		},
		Index: Index{
			BoLa:   "ON",
			Alarm:  "ON",
			Fan:    "ON",
			FanSts: 0,
			Buz:    "ON",
		},
		MaPowInSt:   MaPowInSt{Ac220: 0, Dc12: 0},
		PowModInfos: []PowModInfo{},
		Switch: Switch{
			ChanInfo: make([]ChanInfo, 0),
		},
		Wireless: WirelessOfState{
			Status: 0,
			Csq:    14,
		},
		Time: "2021-01-02 05:05:05",
		Bias: Bias{
			Yaw:   10,
			Roll:  10,
			Pitch: 10,
		},
		Pos: PosOfState{
			Longitude: 4.50,
			Latitude:  6.77,
		},
	}
}
