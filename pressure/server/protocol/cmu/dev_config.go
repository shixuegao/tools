package cmu

const (
	CmdGetDevConfig = "GetDevConfig"
	CmdPutDevConfig = "PutDevConfig"
)

type DevConfigWrapper struct {
	Cmd              string           `json:"CMD"`
	No               string           `json:"No"`
	Guard            Guard            `json:"Guard"`
	Threshold        Threshold        `json:"Threshold"`
	Net              Net              `json:"Net"`
	WirelessOfConfig WirelessOfConfig `json:"Wireless"`
	Interval         int              `json:"Interval"`
	HeartInt         int              `json:"HeartInt"`
	Protocol         Protocol         `json:"Protocol"`
	ForbidPing       int              `json:"ForbidPing"`
	ONUDiagnose      int              `json:"ONUDiagnose"`
	CommType         int              `json:"CommType"`
}

type Guard struct {
	DoorPsw  string `json:"DoorPsw"`
	Interval int    `json:"Interval"`
	EnterNum int    `json:"EnterNum"`
}

type Threshold struct {
	FanTemVal int `json:"FanTemVal"`
	HeatMax   int `json:"HeatMax"`
	Yaw       int `json:"Yaw"`
	Roll      int `json:"Roll"`
	Pitch     int `json:"Pitch"`
}

type Net struct {
	Local  Local  `json:"Local"`
	Server Server `json:"Server"`
}

type Local struct {
	IP     string `json:"IP"`
	Mask   string `json:"MASK"`
	Gate   string `json:"GATE"`
	AutoIp int    `json:"AutoIp"`
}

type Server struct {
	IP   string `json:"IP"`
	Port int    `json:"Port"`
}

type WirelessOfConfig struct {
	IP   string `json:"IP"`
	Port int    `json:"Port"`
}

type Protocol struct {
	Type    int    `json:"Type"`
	Version string `json:"Version"`
}

type DevConfigRequest struct {
	Cmd string `json:"CMD"`
}

type DevConfigResponse struct {
	Cmd        string `json:"CMD"`
	No         string `json:"No"`
	ResultCode int    `json:"resultCode"`
}

func NewDevConfig(number string) *DevConfigWrapper {
	return &DevConfigWrapper{
		Cmd: CmdGetDevConfig,
		No:  number,
		Guard: Guard{
			DoorPsw:  "1234",
			Interval: 2,
			EnterNum: 4,
		},
		Threshold: Threshold{
			FanTemVal: 5,
			HeatMax:   6,
			Yaw:       11,
			Roll:      11,
			Pitch:     11,
		},
		Net: Net{
			Local: Local{
				IP:     "127.0.0.1",
				Mask:   "255.255.255.255",
				Gate:   "255.255.255.255",
				AutoIp: 0,
			},
			Server: Server{
				IP:   "127.0.0.1",
				Port: 10000,
			},
		},
		WirelessOfConfig: WirelessOfConfig{
			IP:   "127.0.0.1",
			Port: 1000,
		},
		Interval: 2,
		HeartInt: 3,
		Protocol: Protocol{
			Type:    0,
			Version: "2.0",
		},
		ForbidPing:  0,
		ONUDiagnose: 0,
		CommType:    1,
	}
}
