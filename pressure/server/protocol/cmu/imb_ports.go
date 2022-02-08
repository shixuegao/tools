package cmu

const (
	CmdPutImbPorts = "Report"
)

var onu = "onu"
var heat = "heat"
var DefaultUploadArray = [][]string{
	{},
	{},
	{onu},
	{onu},
	{heat},
	{heat},
	{onu, heat},
	{heat, onu},
}

type ImbPorts struct {
	Cmd   string    `json:"CMD"`
	No    string    `json:"No"`
	Time  string    `json:"Time"`
	Ports []ImbPort `json:"Ports"`
}

type ImbPort struct {
	Name string `json:"Name"`
	Phy  string `json:"Phy"`
}

func NewImbPorts(number string, info []string) *ImbPorts {
	return &ImbPorts{
		Cmd:   CmdPutImbPorts,
		No:    number,
		Time:  "2020-01-02 12:00:00",
		Ports: infoToPorts(info),
	}
}

func infoToPorts(info []string) (ipts []ImbPort) {
	ipts = []ImbPort{}
	for _, v := range info {
		switch v {
		case onu:
			{
				ipt := ImbPort{Name: "ONU", Phy: "DC12V"}
				ipts = append(ipts, ipt)
			}
		case heat:
			{
				ipt := ImbPort{Name: "Heat", Phy: "AC2"}
				ipts = append(ipts, ipt)
			}
		}
	}
	return
}
