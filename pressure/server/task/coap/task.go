package coap

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
	"umx/tools/pressure/server/interact"
	"umx/tools/pressure/server/interact/http"
	"umx/tools/pressure/server/protocol"
	"umx/tools/pressure/server/protocol/cmu"
	"umx/tools/pressure/server/protocol/cmuV30"
	"umx/tools/pressure/server/task/statistic"
)

const (
	CoapType = "coap"
	//minInterval 最小的睡眠周期
	minInterval = 10 * time.Millisecond
)

type info struct {
	name       string
	localIP    string
	startPort  int
	portCount  int
	dstIp      string
	dstPort    int
	devNumbers []interact.DevNumber
}

type params struct {
	Heartbeat int32 `json:"Heartbeat"` //秒
	DevState  int32 `json:"DevState"`  //秒
	EventInfo int32 `json:"EventInfo"` //秒
	Timeout   int32 `json:"Timeout"`   //毫秒
	Lost      int32 `json:"Lost"`      //毫秒
}

type response struct {
	ResultCode int `json:"resultCode"`
}

func validateParams(val params) error {
	if val.Heartbeat <= 0 {
		return errors.New("非法的心跳周期")
	}
	if val.DevState <= 0 {
		return errors.New("非法的测量值与状态上报周期")
	}
	if val.EventInfo <= 0 {
		return errors.New("非法的事件上报周期")
	}
	if val.Timeout <= 0 {
		return errors.New("非法的超时阈值")
	}
	if val.Lost <= 0 {
		return errors.New("非法的丢包阈值")
	}
	if val.Timeout > val.Lost {
		return errors.New("超时阈值不能大于丢包阈值")
	}
	return nil
}

//Task 任务实现
type Task struct {
	info
	params
	conns    map[string]*Connector
	handlers map[string]FuncHandler
	wait     *sync.WaitGroup
	state    int32
	closed   chan struct{}
	lock     *sync.Mutex
	//statistic
	statistic *statistic.Statistic
}

//Init 初始化
func (t *Task) init() (err error) {
	count := 0
	for i := 0; i < t.portCount; i++ {
		number := ""
		port := i + t.startPort
		devType := protocol.DevTypeEconomic
		subDevType := ""
		if nil == t.devNumbers {
			number = devNumber(port)
			devType = randomDevType()
			subDevType = protocol.RandomSubDevType(devType)
		} else {
			number = t.devNumbers[count].Number
			devType = protocol.GetDevType(t.devNumbers[count].DevType)
			subDevType = protocol.RandomSubDevType(devType)
		}
		if conn, e := NewConnector(devType, t.localIP, port, t.dstIp, t.dstPort, subDevType); nil == e {
			t.conns[number] = conn
			count++
		} else {
			err = e
		}
	}
	if count <= 0 {
		return errors.New("初始化Task失败, 无法建立Udp连接, 可能的原因-->" + err.Error())
	}
	return nil
}

func (t *Task) devTypeBuNumber(number string) int {
	conn := t.conns[number]
	if conn == nil {
		return protocol.DevTypeEconomic
	}
	return conn.devType
}

func randomDevType() int {
	i := rand.Intn(2)
	switch i {
	case protocol.DevTypeEconomic:
		return protocol.DevTypeEconomic
	case protocol.DevTypeSynthetic:
		return protocol.DevTypeSynthetic
	}
	return protocol.DevTypeEconomic
}

//Start
func (t *Task) Start() error {
	t.lock.Lock()
	defer t.lock.Unlock()
	if atomic.LoadInt32(&t.state) == 1 {
		return nil
	}
	if err := t.init(); nil != err {
		return err
	}
	atomic.StoreInt32(&t.state, 1)
	t.receive()
	t.sendDevInfo()
	t.updateStatistic()
	t.sendIntervalDevstate()
	t.sendIntervalHeartbeat()
	t.sendIntervalEventInfo()
	t.sendIntervalBroadcast()
	t.sendIntervalImbPorts()
	return nil
}

//Close
func (t *Task) Close() error {
	t.lock.Lock()
	defer t.lock.Unlock()
	if atomic.LoadInt32(&t.state) == 0 {
		return nil
	}
	select {
	case <-t.closed:
		return errors.New("任务已经关闭")
	default:
		close(t.closed)
		for _, v := range t.conns {
			if err := v.Close(); nil != err {
				//do nothing
			}
		}
		t.wait.Wait()
		atomic.StoreInt32(&t.state, 0)
		return nil
	}
}

func (t *Task) tState() string {
	if state := atomic.LoadInt32(&t.state); state == 1 {
		return "Running"
	} else {
		return "Stopped"
	}
}

func (t *Task) Show() http.Show {
	return http.Show{
		Type:      CoapType,
		Name:      t.name,
		LocalIP:   t.localIP,
		IP:        t.dstIp,
		Port:      t.dstPort,
		State:     t.tState(),
		StartPort: t.startPort,
		PortCount: t.portCount,
	}
}

//device numbers
func (t *Task) DevNumbers() []interact.DevNumber {
	index := 0
	devNumbers := make([]interact.DevNumber, len(t.conns))
	for k, v := range t.conns {
		devNumbers[index] = interact.DevNumber{
			DevType:    protocol.GetDevTypeName(v.devType),
			SubDevType: v.subDevType,
			Number:     k,
		}
		index++
	}
	return devNumbers
}

func (t *Task) Params() interface{} {
	return t.params
}

func (t *Task) SetParams(v interface{}) error {
	if vv, ok := v.(map[string]interface{}); ok {
		return t.setParamsByMap(vv)
	}
	return nil
}

func (t *Task) Statistic() statistic.Statistic {
	sta := t.statistic
	if sta == nil {
		return statistic.Statistic{}
	}
	return statistic.Statistic{
		Delay:   sta.Delay,
		Total:   sta.Total,
		Timeout: sta.Timeout,
		Lost:    sta.Lost,
	}
}

func (t *Task) setParamsByMap(v map[string]interface{}) error {
	val := params{}
	if flt, ok := v["Heartbeat"].(float64); ok {
		val.Heartbeat = int32(flt)
	}
	if flt, ok := v["DevState"].(float64); ok {
		val.DevState = int32(flt)
	}
	if flt, ok := v["EventInfo"].(float64); ok {
		val.EventInfo = int32(flt)
	}
	if flt, ok := v["Timeout"].(float64); ok {
		val.Timeout = int32(flt)
	}
	if flt, ok := v["Lost"].(float64); ok {
		val.Lost = int32(flt)
	}
	if err := validateParams(val); err != nil {
		return err
	}
	fmt.Printf("---------------------------------------\n"+
		"Coap任务[%s]的参数发生改变:\n"+
		"心跳周期(秒): %d -> %d\n"+
		"测量值上报周期(秒): %d -> %d\n"+
		"事件上报周期(秒): %d -> %d\n"+
		"超时阈值(ms): %d -> %d\n"+
		"丢包阈值(ms): %d -> %d\n"+
		"---------------------------------------\n",
		t.name,
		t.Heartbeat, val.Heartbeat,
		t.DevState, val.DevState,
		t.EventInfo, val.EventInfo,
		t.Timeout, val.Timeout,
		t.Lost, val.Lost)
	atomic.StoreInt32(&t.Heartbeat, val.Heartbeat)
	atomic.StoreInt32(&t.DevState, val.DevState)
	atomic.StoreInt32(&t.EventInfo, val.EventInfo)
	atomic.StoreInt32(&t.Timeout, val.Timeout)
	atomic.StoreInt32(&t.Lost, val.Lost)
	return nil
}

//sendDevInfo 发送设备信息
func (t *Task) sendDevInfo() {
	t.wait.Add(1)
	go func() {
		defer t.wait.Done()
		message := &Message{
			Type: Confirmable,
			Code: PUT,
		}
		// var buffer bytes.Buffer
		for number, conn := range t.conns {
			devInfo := devInfo(number, conn.devType)
			payload, err := json.Marshal(devInfo)
			if nil != err {
				continue
			}
			message.MessageID = conn.getMessageID()
			message.Payload = payload
			message.SetPathString(devInfoUrl(conn.devType))
			//发送
			err = conn.Send(message)
			if nil != err {
				//do nothing
			}
			//睡眠是10ms
			if !t.intervalSleep(10 * time.Millisecond) {
				return
			}
		}
	}()
}

func devInfoUrl(devType int) string {
	if devType == protocol.DevTypeEconomic {
		return cmu.UrlDevInfo
	} else {
		return cmuV30.UrlDevInfo
	}
}

func devInfo(number string, devType int) interface{} {
	if devType == protocol.DevTypeEconomic {
		return cmu.NewDevInfo(number)
	} else {
		return cmuV30.NewDevInfo(number)
	}
}

//receive 接收数据(循环接收)
func (t *Task) receive() {
	for num, conn := range t.conns {
		num0 := num
		conn0 := conn
		t.wait.Add(1)
		go func() {
			defer t.wait.Done()
			// var buffer bytes.Buffer
			for {
				if m, err := conn0.Receive(); nil != err {
					return
				} else {
					if m.Type == Acknowledgement {
						//接收响应
						resp := response{}
						if err := json.Unmarshal(m.Payload, &resp); nil == err && resp.ResultCode == 255 {
							t.register(num0, conn0)
						}
					} else {
						//接收请求
						if handler, ok := t.handlers["/"+m.PathString()]; ok {
							handler(num0, conn0, &m)
						}
					}
				}
				select {
				case <-t.closed:
					return
				default:
				}
			}
		}()
	}
}

func (t *Task) updateStatistic() {
	t.wait.Add(1)
	go func() {
		defer t.wait.Done()
		array := make([]*statistic.Statistic, len(t.conns))
		for {
			timeout := atomic.LoadInt32(&t.Timeout)
			lost := atomic.LoadInt32(&t.Lost)
			if !t.intervalSleep(time.Duration(2*lost) * time.Millisecond) {
				return
			}
			count := 0
			for _, conn := range t.conns {
				sta := conn.statistic.clearStatistic(int64(timeout), int64(lost))
				array[count] = &sta
				count++
			}
			assembled := statistic.CaculateBunchOfStatistic(array)
			t.statistic = &assembled
		}
	}()
}

//注册
func (t *Task) register(num string, conn *Connector) {
	if conn.registered {
		return
	}
	registerPkt := cmuV30.DefaultRegister(num)
	payload, _ := json.Marshal(registerPkt)
	message := &Message{
		Type:      Confirmable,
		Code:      POST,
		MessageID: conn.getMessageID(),
		Payload:   payload,
	}
	message.SetPathString(cmuV30.UrlRegister)
	err := conn.Send(message)
	if nil != err {
		//do nothing
	}
	conn.registered = true
}

//新建coap任务
func NewTask(name, localIP string, startPort, portCount int, ip string, port int, devNumbers []interact.DevNumber) *Task {
	t := Task{
		info: info{
			name:       name,
			localIP:    localIP,
			startPort:  startPort,
			portCount:  portCount,
			dstIp:      ip,
			dstPort:    port,
			devNumbers: devNumbers,
		},
		//单位: 秒
		params: params{
			Heartbeat: 30,
			DevState:  120,
			EventInfo: 300,
			Timeout:   100, //默认100ms没收到回包则视为超时
			Lost:      500, //默认500ms没收到回包则视为丢包
		},
		conns:     make(map[string]*Connector),
		handlers:  make(map[string]FuncHandler),
		wait:      new(sync.WaitGroup),
		closed:    make(chan struct{}),
		statistic: new(statistic.Statistic),
		lock:      new(sync.Mutex),
	}
	t.handlers[cmu.UrlDevInfo] = receiveDevInfo(&t)
	t.handlers[cmu.UrlDevState] = receiveDevState(&t)
	t.handlers[cmu.UrlDevConfig] = receiveDevConfig(&t)
	t.handlers[cmu.UrlDevControl] = receiveDevControl(&t)
	//v30
	t.handlers[cmuV30.UrlDevInfo] = receiveDevInfo(&t)
	t.handlers[cmuV30.UrlDevState] = receiveDevState(&t)
	t.handlers[cmuV30.UrlDevConfig] = receiveDevConfig(&t)
	t.handlers[cmuV30.UrlDevControl] = receiveDevControl(&t)
	t.handlers[cmuV30.UrlSwitch] = receiveSwitch(&t)
	return &t
}
