package interact

const (
	Success = iota
	Failure
)

const (
	PathTaskShow       = "/task/show"
	PathTaskAdd        = "/task/add"
	PathTaskAddContent = "/task/addContent"
	PathTaskOrder      = "/task/order"
	PathTaskParams     = "/task/params"
)

const (
	OrderStart     = "start"
	OrderClose     = "close"
	OrderNumbers   = "numbers"
	OrderStatus    = "status"
	OrderParams    = "params"
	OrderStatistic = "statistic"
)

type Wrapper struct {
	Version uint8  //版本号
	PLength uint16 //路径长度
	Path    []byte //路径
	CLength uint32 //内容长度
	Content []byte //内容
}

type TypeName struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type Request struct {
	TypeName
	Order string `json:"order"` //start, stop, numbers, status
}

type Add struct {
	TypeName
	LocalIP   string `json:"localIP"`
	IP        string `json:"ip"`
	Port      int    `json:"port"`
	StartPort int    `json:"startPort"`
	PortCount int    `json:"portCount"`
}

type AddContent struct {
	TypeName
	LocalIP    string      `json:"localIP"`
	IP         string      `json:"ip"`
	Port       int         `json:"port"`
	StartPort  int         `json:"startPort"`
	DevNumbers []DevNumber `json:"devNumbers"`
}

type Params struct {
	Type string      `json:"type"`
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

type DevNumber struct {
	DevType    string
	SubDevType string
	Number     string
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
