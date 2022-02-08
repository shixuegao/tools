package cmu

const (
	CmdPutDevMsg = "PutDevMsg"
	CmdGetDevMsg = "GetDevMsg"
)

type DevStateWrapper struct {
	Cmd      string          `json:"CMD"`
	No       string          `json:"No"`
	Conect   bool            `json:"Conect"`
	Environ  Environ         `json:"Environ"`
	Door     Door            `json:"Door"`
	Spd      Spd             `json:"SPD"`
	Esw      Esw             `json:"eSW"`
	Index    Index           `json:"Index"`
	Ext      Ext             `json:"EXT"`
	Cur      Cur             `json:"CUR"`
	Vol      Vol             `json:"VOL"`
	Battery  int             `json:"Battery"`
	AC220    int             `json:"AC220"`
	DC12     int             `json:"DC12"`
	Wireless WirelessOfState `json:"Wireless"`
	Time     string          `json:"Time"`
	Bias     Bias            `json:"Bias"`
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
	Heat   string `json:"Heat"`
	Fan    string `json:"Fan"`
	FanSts int    `json:"FanSts"`
	Buz    string `json:"Buz"`
}

type Ext struct {
	Ac1  string `json:"AC1"`
	Ac2  string `json:"AC2"`
	Ac3  string `json:"AC3"`
	Ac4  string `json:"AC4"`
	Ac5  string `json:"AC5"`
	Ac6  string `json:"AC6"`
	Ac7  string `json:"AC7"`
	Ac8  string `json:"AC8"`
	Ipc1 string `json:"IPC1"`
	Ipc2 string `json:"IPC2"`
	Ipc3 string `json:"IPC3"`
	Ipc4 string `json:"IPC4"`
	Ipc5 string `json:"IPC5"`
	Ipc6 string `json:"IPC6"`
}

type Cur struct {
	Ac1  int `json:"AC1"`
	Ac2  int `json:"AC2"`
	Ac3  int `json:"AC3"`
	Ac4  int `json:"AC4"`
	Ac5  int `json:"AC5"`
	Ac6  int `json:"AC6"`
	Ac7  int `json:"AC7"`
	Ac8  int `json:"AC8"`
	Ipc1 int `json:"IPC1"`
	Ipc2 int `json:"IPC2"`
	Ipc3 int `json:"IPC3"`
	Ipc4 int `json:"IPC4"`
	Ipc5 int `json:"IPC5"`
	Ipc6 int `json:"IPC6"`
}

type Vol struct {
	Ac   int `json:"AC"`
	Ipc1 int `json:"IPC1"`
	Ipc2 int `json:"IPC2"`
	Ipc3 int `json:"IPC3"`
	Ipc4 int `json:"IPC4"`
	Ipc5 int `json:"IPC5"`
	Ipc6 int `json:"IPC6"`
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
			Heat:   "ON",
			Fan:    "ON",
			FanSts: 0,
			Buz:    "ON",
		},
		Ext: Ext{
			Ac1:  "ON",
			Ac2:  "ON",
			Ac3:  "ON",
			Ac4:  "ON",
			Ipc1: "ON",
			Ipc2: "ON",
			Ipc3: "ON",
		},
		Cur: Cur{
			Ac1:  55,
			Ac2:  55,
			Ac3:  55,
			Ac4:  55,
			Ipc1: 55,
			Ipc2: 55,
			Ipc3: 55,
		},
		Vol: Vol{
			Ac:   220,
			Ipc1: 220,
			Ipc2: 220,
			Ipc3: 220,
		},
		Battery: 50,
		AC220:   0,
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
	}
}
