package cmuV30

import (
	"math/rand"
	"umx/tools/pressure/server/util"
)

const (
	CmdGetEvent = "GetEvent"
	CmdPutEvent = "PutEvent"
)

var events = []Event{
	{1, 0, "市电掉电"},
	{2, 0, "有线网络断开"},
	{3, 0, "无线NBIOT网络断开"},
	{4, 0, "重合闸跳闸"},
	{5, 0, "DC12V掉电"},
	{6, 0, "自动重合闸故障"},
	{7, 0, "自动重合闸短路"},
	{8, 0, "电源防雷器失效"},
	{9, 0, "箱门打开"},
	{10, 0, "撤防"},
	{11, 0, "水浸"},
	{12, 0, "AC220V过流"},
	{13, 0, "AC220V过压"},
	{14, 0, "AC220V低压"},
	{15, 0, "锂电池失效"},
	{22, 0, "IP冲突"},
	{23, 0, "网络设备异常"},
	{24, 0, "监控箱倾斜"},
	{25, 0, "补光灯异常"},
	{26, 0, "闪光等异常"},
	{27, 1, "IPC网络异常"},
	{28, 1, "IPC电源异常"},
	{29, 0, "电源模块通讯异常"},
}

type Event struct {
	Id     int
	Type   int
	Remark string
}

type EventWrapper struct {
	Cmd      string `json:"CMD"`
	No       string `json:"No"`
	Event    int    `json:"Event"`
	EveState bool   `json:"EveState"`
	EveType  int    `json:"EveType"`
	PowerCh  int    `json:"PowerCh"`
	SwitchCh int    `json:"SwitchCh"`
	More     int    `json:"More"`
	Time     string `json:"Time"`
}

//RandomEvent 获取随机的事件
func RandomEvent() (Event, bool) {
	state := rand.Intn(2) == 1
	size := len(events)
	index := rand.Intn(size)
	return events[index], state
}

func RandomMore(event int) int {
	if event == 25 || event == 26 {
		return rand.Intn(2)
	}
	return 0
}

func RandomPowerCh() int {
	return rand.Intn(20) + 1
}

func RandomSwitchCh() int {
	return rand.Intn(8) + 1
}

//生成新的事件包
func NewEventWrapper(number string, event Event, state bool) *EventWrapper {
	return &EventWrapper{
		Cmd:      CmdPutEvent,
		No:       number,
		Event:    event.Id,
		EveState: state,
		EveType:  event.Type,
		PowerCh:  RandomPowerCh(),
		SwitchCh: RandomSwitchCh(),
		More:     RandomMore(event.Id),
		Time:     util.CurrentTime(),
	}
}
