package protocol

import "math/rand"

const (
	//经济版
	DevTypeEconomic = iota
	//综合版
	DevTypeSynthetic
)

var subDevTypeEconomic = []string{"UMX-3RTU1", "UMX-3RTU2"}
var subDevTypeSynthetic = []string{"UMX-3IMB3A570-C13P08WLP"}

func GetDevType(name string) int {
	if name == "UMX-3RTU2" {
		return DevTypeEconomic
	} else if name == "UMX-3RTU3" {
		return DevTypeSynthetic
	}
	return DevTypeEconomic
}

func GetDevTypeName(devType int) string {
	if devType == DevTypeEconomic {
		return "UMX-3RTU2"
	} else if devType == DevTypeSynthetic {
		return "UMX-3RTU3"
	}
	return "UMX-3RTU2"
}

//随机的子类型
func RandomSubDevType(devType int) string {
	if devType == DevTypeSynthetic {
		index := rand.Intn(len(subDevTypeSynthetic))
		return subDevTypeSynthetic[index]
	} else {
		index := rand.Intn(len(subDevTypeEconomic))
		return subDevTypeEconomic[index]
	}
}
