package cmu

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
	{16, 1, "IPC1异常离线"},
	{17, 1, "IPC2异常离线"},
	{18, 1, "IPC3异常离线"},
	{19, 1, "IPC4异常离线"},
	{20, 1, "IPC5异常离线"},
	{21, 1, "IPC6异常离线"},
	{22, 0, "IP冲突"},
	{23, 0, "网络设备异常"},
	{24, 0, "监控箱倾斜"},
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
	Time     string `json:"Time"`
}

//RandomEvent 获取随机的事件
func RandomEvent() (Event, bool) {
	state := rand.Intn(2) == 1
	size := len(events)
	index := rand.Intn(size)
	return events[index], state
}

//生成新的事件包
func NewEventWrapper(number string, event Event, state bool) *EventWrapper {
	return &EventWrapper{
		Cmd:      CmdPutEvent,
		No:       number,
		Event:    event.Id,
		EveState: state,
		EveType:  event.Type,
		Time:     util.CurrentTime(),
	}
}
