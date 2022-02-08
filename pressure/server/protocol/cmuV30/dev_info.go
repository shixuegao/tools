package cmuV30

const (
	CmdGetDevInfo = "GetDevInfo"
	CmdPutDevInfo = "PutDevInfo"
)

type DevInfoWrapper struct {
	Cmd       string    `json:"CMD"`
	No        string    `json:"No"`
	Sn        string    `json:"SN"`
	IP        string    `json:"IP"`
	Mac       string    `json:"MAC"`
	Icid      string    `json:"ICID"`
	Imsi      string    `json:"IMSI"`
	Imei      string    `json:"IMEI"`
	Hw        string    `json:"HW"`
	Sw        string    `json:"SW"`
	Date      string    `json:"DATE"`
	Pos       Pos       `json:"Pos"`
	DevName   string    `json:"DevName"`
	PowMod    PowMod    `json:"PowMod"`
	SwitchMod SwitchMod `json:"SwitchMod"`
}

type Pos struct {
	Longitude float32 `json:"Longitude"`
	Latitude  float32 `json:"Latitude"`
}

type PowMod struct {
	Hw   string `json:"HW"`
	Sw   string `json:"SW"`
	Date string `json:"DATE"`
}

type SwitchMod struct {
	Model string `json:"Model"`
}

type DevInfoRequest struct {
	Cmd string `json:"CMD"`
}

func NewDevInfo(number string) *DevInfoWrapper {
	return &DevInfoWrapper{
		Cmd:  CmdPutDevInfo,
		Sn:   "0000",
		IP:   "127.0.0.1",
		Mac:  "00-00-00-00-00-00",
		Icid: "0000",
		Imsi: "0000",
		Imei: "0000",
		Hw:   "V00.00",
		Sw:   "V00.00",
		Date: "2020-01-02 12:00:00",
		Pos: Pos{
			Longitude: 4.50,
			Latitude:  6.77,
		},
		DevName: "Virtual Device",
		PowMod: PowMod{
			Hw:   "v00.00",
			Sw:   "v00.00",
			Date: "2000-01-01 00:00:00",
		},
		SwitchMod: SwitchMod{
			Model: "V9.9",
		},
	}
}
