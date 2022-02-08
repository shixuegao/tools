package coap

import (
	"encoding/json"
	"sync/atomic"
	"time"
	"umx/tools/pressure/server/protocol"
	"umx/tools/pressure/server/protocol/cmu"
	"umx/tools/pressure/server/protocol/cmuV30"
	"umx/tools/pressure/server/util"
)

//sendIntervalHeartbeat 定期发送心跳
func (t *Task) sendIntervalHeartbeat() {
	t.wait.Add(1)
	go func() {
		defer t.wait.Done()
		for {
			interval := atomic.LoadInt32(&t.Heartbeat)
			if !t.intervalSleep(time.Duration(interval) * time.Second) {
				return
			}
			//send event
			for number, conn := range t.conns {
				heartbeat := keepLive(number, conn.devType)
				message := &Message{
					Type:      Confirmable,
					Code:      POST,
					MessageID: conn.getMessageID(),
				}
				payload, err := json.Marshal(heartbeat)
				if nil != err {
					continue
				}
				message.Payload = payload
				message.SetPathString(keepLiveUrl(conn.devType))
				err = conn.Send(message)
				if nil != err {
					//do nothing
				}
				//睡眠10ms
				if !t.intervalSleep(10 * time.Millisecond) {
					return
				}
			}
		}
	}()
}

func keepLiveUrl(devType int) string {
	if devType == protocol.DevTypeEconomic {
		return cmu.UrlKeepLive
	} else {
		return cmuV30.UrlKeepLive
	}
}

func keepLive(number string, devType int) interface{} {
	if devType == protocol.DevTypeEconomic {
		return cmu.NewKeepLiveWrapper(number)
	} else {
		return cmuV30.NewKeepLiveWrapper(number)
	}
}

//sendIntervalDevstate 定期发送设备状态
func (t *Task) sendIntervalDevstate() {
	t.wait.Add(1)
	go func() {
		defer t.wait.Done()
		for {
			interval := atomic.LoadInt32(&t.DevState)
			if !t.intervalSleep(time.Duration(interval) * time.Second) {
				return
			}
			for number, conn := range t.conns {
				devState := devState(number, conn.devType)
				message := &Message{
					Type:      Confirmable,
					Code:      PUT,
					MessageID: conn.getMessageID(),
				}
				payload, err := json.Marshal(devState)
				if nil != err {
					continue
				}
				message.Payload = payload
				message.SetPathString(devStateUrl(conn.devType))
				err = conn.Send(message)
				if nil != err {
					//do nothing
				}
				//睡眠10ms
				if !t.intervalSleep(10 * time.Millisecond) {
					return
				}
			}
		}
	}()
}

func devStateUrl(devType int) string {
	if devType == protocol.DevTypeEconomic {
		return cmu.UrlDevState
	} else {
		return cmuV30.UrlDevState
	}
}

func devState(number string, devType int) interface{} {
	if devType == protocol.DevTypeEconomic {
		return cmu.NewDevState(number)
	} else {
		return cmuV30.NewDevState(number)
	}
}

//sendIntervalEventInfo 定期发送事件
func (t *Task) sendIntervalEventInfo() {
	t.wait.Add(1)
	go func() {
		defer t.wait.Done()
		for {
			interval := atomic.LoadInt32(&t.EventInfo)
			if !t.intervalSleep(time.Duration(interval) * time.Second) {
				return
			}
			timeStr := util.CurrentTime()
			for number, conn := range t.conns {
				eventWrapper := eventInfo(number, timeStr, conn.devType)
				message := &Message{
					Type:      Confirmable,
					Code:      PUT,
					MessageID: conn.getMessageID(),
				}
				payload, err := json.Marshal(eventWrapper)
				if nil != err {
					continue
				}
				message.Payload = payload
				message.SetPathString(eventInfoUrl(conn.devType))
				err = conn.Send(message)
				if nil != err {
					//do nothing
				}
				//睡眠10ms
				if !t.intervalSleep(10 * time.Millisecond) {
					return
				}
			}
		}
	}()
}

//sendIntervalBroadcast 定期发送广播回送
func (t *Task) sendIntervalBroadcast() {
	t.wait.Add(1)
	go func() {
		defer t.wait.Done()
		for {
			interval := 10
			if !t.intervalSleep(time.Duration(interval) * time.Second) {
				return
			}
			for number, conn := range t.conns {
				broadcastPkt := cmuV30.DefaultBroadcast(number)
				payload, _ := json.Marshal(broadcastPkt)
				message := &Message{
					Type:      Confirmable,
					Code:      POST,
					MessageID: conn.getMessageID(),
					Payload:   payload,
				}
				message.SetPathString(cmuV30.UrlBroadcast)
				err := conn.Send(message)
				if nil != err {
					//do nothing
				}
				//睡眠10ms
				if !t.intervalSleep(10 * time.Millisecond) {
					return
				}
			}
		}
	}()
}

func eventInfoUrl(devType int) string {
	if devType == protocol.DevTypeEconomic {
		return cmu.UrlEventInfo
	} else {
		return cmuV30.UrlEventInfo
	}
}

func eventInfo(number, timeStr string, devType int) interface{} {
	if devType == protocol.DevTypeEconomic {
		event, state := cmu.RandomEvent()
		eventWrapper := &cmu.EventWrapper{
			Cmd:      cmu.CmdPutEvent,
			No:       number,
			Event:    event.Id,
			EveType:  event.Type,
			EveState: state,
			Time:     timeStr,
		}
		return eventWrapper
	} else {
		event, state := cmuV30.RandomEvent()
		eventWrapper := &cmuV30.EventWrapper{
			Cmd:      cmuV30.CmdPutEvent,
			No:       number,
			Event:    event.Id,
			EveType:  event.Type,
			EveState: state,
			PowerCh:  cmuV30.RandomPowerCh(),
			SwitchCh: cmuV30.RandomSwitchCh(),
			More:     cmuV30.RandomMore(event.Id),
			Time:     timeStr,
		}
		return eventWrapper
	}
}

//intervalSleep 睡眠一定时期
func (t *Task) intervalSleep(total time.Duration) bool {
	if total < minInterval {
		total = minInterval
	}
	count := int(total / minInterval)
	for i := 0; i < count; i++ {
		time.Sleep(minInterval)
		select {
		case <-t.closed:
			return false
		default:
		}
	}
	return true
}

//sendIntervalImbPorts 定时发送Imb端口信息(周期性发送来检测服务器的稳定性)
func (t *Task) sendIntervalImbPorts() {
	t.wait.Add(1)
	go func() {
		defer t.wait.Done()
		index := 0
		total := len(cmu.DefaultUploadArray)
		for {
			//每120秒送一次
			if !t.intervalSleep(10 * time.Second) {
				return
			}
			for number, conn := range t.conns {
				imbPortsPkt := cmu.NewImbPorts(number, cmu.DefaultUploadArray[index])
				payload, _ := json.Marshal(imbPortsPkt)
				message := &Message{
					Type:      Confirmable,
					Code:      PUT,
					MessageID: conn.getMessageID(),
					Payload:   payload,
				}
				message.SetPathString(cmu.UrlImbPorts)
				err := conn.Send(message)
				if nil != err {
					//do nothing
				}
				//睡眠10ms
				if !t.intervalSleep(10 * time.Millisecond) {
					return
				}
			}
			index++
			if index >= total {
				index = 0
			}
		}
	}()
}
