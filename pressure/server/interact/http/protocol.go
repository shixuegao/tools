package http

type Show struct {
	Type      string `json:"type"`
	Name      string `json:"name"`
	LocalIP   string `json:"localIP"`
	IP        string `json:"ip"`
	Port      int    `json:"port"`
	State     string `json:"state"`
	StartPort int    `json:"startPort"`
	PortCount int    `json:"portCount"`
}
