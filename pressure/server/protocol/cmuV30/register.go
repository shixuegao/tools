package cmuV30

type Register struct {
	DeviceNumber string `json:"deviceNumber"`
	DeviceIP     string `json:"deviceIP"`
	DevicePort   int    `json:"devicePort"`
	Expire       int    `json:"expire"`
}

func DefaultRegister(deviceNumber string) Register {
	return Register{
		DeviceNumber: deviceNumber,
		DeviceIP:     "127.0.0.1",
		DevicePort:   5683,
		Expire:       60,
	}
}
