package statistic

import (
	"fmt"
)

type Statistic struct {
	Total   int64 //包总数
	Delay   int64 //回包平均延迟
	Timeout int64 //收发超时的包
	Lost    int64 //丢失的包
}

func (s Statistic) String() string {
	return fmt.Sprintf("[总包数: %d | 回包平均延迟(ms): %d | 收发超时包数: %d | 丢包数: %d]", s.Total, s.Delay, s.Timeout, s.Lost)
}

func CaculateBunchOfStatistic(array []*Statistic) Statistic {
	var delay int64
	sta := Statistic{}
	for _, v := range array {
		sta.Total += v.Total
		sta.Timeout += v.Timeout
		sta.Lost += v.Lost
		delay += v.Delay
	}
	sta.Delay = delay / int64(len(array))
	return sta
}
