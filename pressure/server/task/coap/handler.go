package coap

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"umx/tools/pressure/server/protocol"
	"umx/tools/pressure/server/protocol/cmu"
	"umx/tools/pressure/server/protocol/cmuV30"
)

type FuncHandler func(number string, conn *Connector, m *Message)

//receiveDevInfo 设备信息接收
func receiveDevInfo(t *Task) FuncHandler {
	f := func(number string, conn *Connector, m *Message) {
		var payload []byte
		devType := t.devTypeBuNumber(number)
		if devType == protocol.DevTypeEconomic {
			req := cmu.DevInfoRequest{}
			if err := json.Unmarshal(m.Payload, &req); err != nil {
				return
			}
			devInfo := cmu.NewDevInfo(number)
			if p, err := json.Marshal(devInfo); err != nil {
				payload = []byte("{}")
			} else {
				payload = p
			}
		} else {
			req := cmuV30.DevInfoRequest{}
			if err := json.Unmarshal(m.Payload, &req); err != nil {
				return
			}
			devInfo := cmuV30.NewDevInfo(number)
			if p, err := json.Marshal(devInfo); err != nil {
				payload = []byte("{}")
			} else {
				payload = p
			}
		}
		res := &Message{
			Type:      Acknowledgement,
			Code:      Content,
			Token:     m.Token,
			MessageID: m.MessageID,
			Payload:   payload,
		}
		err := conn.Send(res)
		if nil != err {
			//do nothing
		}
	}
	return FuncHandler(f)
}

//receiveDevState
func receiveDevState(t *Task) FuncHandler {
	f := func(number string, conn *Connector, m *Message) {
		var payload []byte
		devType := t.devTypeBuNumber(number)
		if devType == protocol.DevTypeEconomic {
			devState := cmu.NewDevState(number)
			if p, err := json.Marshal(devState); err != nil {
				payload = []byte("{}")
			} else {
				payload = p
			}
		} else {
			devState := cmuV30.NewDevState(number)
			if p, err := json.Marshal(devState); err != nil {
				payload = []byte("{}")
			} else {
				payload = p
			}
		}
		res := &Message{
			Type:      Acknowledgement,
			Code:      Content,
			Token:     m.Token,
			MessageID: m.MessageID,
			Payload:   payload,
		}
		err := conn.Send(res)
		if nil != err {
			//do nothing
		}
	}
	return FuncHandler(f)
}

//receiveDevConfig
func receiveDevConfig(t *Task) FuncHandler {
	f := func(number string, conn *Connector, m *Message) {
		if m.Code == GET {
			payload, err := devConfig(number, conn.devType)
			if nil != err {
				payload = []byte("{}")
			}
			res := &Message{
				Type:      Acknowledgement,
				Token:     m.Token,
				Code:      Content,
				MessageID: m.MessageID,
				Payload:   payload,
			}
			err = conn.Send(res)
			if nil != err {
				//do nothing
			}
		} else if m.Code == PUT {
			devConfigRes := cmu.DevConfigResponse{
				Cmd:        cmu.CmdPutDevConfig,
				No:         number,
				ResultCode: 0,
			}
			payload, err := json.Marshal(&devConfigRes)
			if nil != err {
				payload = []byte("{}")
			}
			res := &Message{
				Type:      Acknowledgement,
				Token:     m.Token,
				Code:      Content,
				MessageID: m.MessageID,
				Payload:   payload,
			}
			err = conn.Send(res)
			if nil != err {
				//do nothing
			}
		}
	}
	return FuncHandler(f)
}

func devConfig(number string, devType int) (payload []byte, err error) {
	if devType == protocol.DevTypeEconomic {
		payload, err = json.Marshal(cmu.NewDevConfig(number))
	} else if devType == protocol.DevTypeSynthetic {
		payload, err = json.Marshal(cmuV30.NewDevConfig(number))
	} else {
		err = errors.New("无法识别的设备类型")
	}
	return
}

//receiveDevControl
func receiveDevControl(t *Task) FuncHandler {
	f := func(number string, conn *Connector, m *Message) {
		var payload []byte
		devType := t.devTypeBuNumber(number)
		if devType == protocol.DevTypeEconomic {
			resp := cmu.DevControlResponse{Cmd: cmu.CmdDevControl, No: number, ResultCode: 0}
			if p, err := json.Marshal(&resp); err != nil {
				payload = []byte("{}")
			} else {
				payload = p
			}
		} else {
			resp := cmuV30.DevControlResponse{Cmd: cmu.CmdDevControl, No: number, ResultCode: 0}
			if p, err := json.Marshal(&resp); err != nil {
				payload = []byte("{}")
			} else {
				payload = p
			}
		}
		res := &Message{
			Type:      Acknowledgement,
			Token:     m.Token,
			Code:      Content,
			MessageID: m.MessageID,
			Payload:   payload,
		}
		err := conn.Send(res)
		if nil != err {
			//do nothing
		}
	}
	return FuncHandler(f)
}

//receiveSwitch
func receiveSwitch(t *Task) FuncHandler {
	f := func(number string, conn *Connector, m *Message) {
		if m.Code == GET {
			payload, err := json.Marshal(cmuV30.NewSwitchV30())
			if nil != err {
				payload = []byte("{}")
			}
			res := &Message{
				Type:      Acknowledgement,
				Token:     m.Token,
				Code:      Content,
				MessageID: m.MessageID,
				Payload:   payload,
			}
			err = conn.Send(res)
			if nil != err {
				//do nothing
			}
		} else if m.Code == PUT {
			resp := cmuV30.SwitchResponse{
				Cmd:        cmu.CmdPutDevConfig,
				No:         number,
				ResultCode: 0,
			}
			payload, err := json.Marshal(&resp)
			if nil != err {
				payload = []byte("{}")
			}
			res := &Message{
				Type:      Acknowledgement,
				Token:     m.Token,
				Code:      Content,
				MessageID: m.MessageID,
				Payload:   payload,
			}
			err = conn.Send(res)
			if nil != err {
				//do nothing
			}
		}
	}
	return FuncHandler(f)
}

//devNumber 生成设备编号
func devNumber(port int) string {
	time := time.Now().Format("01021504")
	strPort := fmt.Sprintf("%08d", port)
	return time + strPort
}
