package cmuV30

type Broadcast struct {
	BroadcastID    int    `json:"broadcastID"`
	DeviceType     string `json:"deviceType"`
	DeviceNumber   string `json:"deviceNumber"`
	DeviceName     string `json:"deviceName"`
	DeviceHardware string `json:"deviceHardware"`
	DeviceSoftware string `json:"deviceSoftware"`
	DeviceIP       string `json:"deviceIP"`
	DevicePort     int    `json:"devicePort"`
}

func DefaultBroadcast(deviceNumber string) Broadcast {
	return Broadcast{
		BroadcastID:    0,
		DeviceType:     "XXXXX",
		DeviceNumber:   deviceNumber,
		DeviceName:     deviceNumber,
		DeviceHardware: "v0.0.0",
		DeviceSoftware: "v0.0.0",
		DeviceIP:       "127.0.0.1",
		DevicePort:     5683,
	}
}
