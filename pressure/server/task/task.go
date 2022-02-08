package task

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"umx/tools/pressure/server/interact"
	"umx/tools/pressure/server/interact/http"
	"umx/tools/pressure/server/interact/udp"
	"umx/tools/pressure/server/log"
	"umx/tools/pressure/server/task/coap"
	"umx/tools/pressure/server/task/statistic"
)

var logger = log.Logger

type Task interface {
	Start() error
	Close() error
	Show() http.Show
	DevNumbers() []interact.DevNumber
	Params() interface{}
	SetParams(v interface{}) error
	Statistic() statistic.Statistic
}

type TaskManager struct {
	tasks map[string]Task
	ports *Ports
	lock  *sync.Mutex
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks: make(map[string]Task),
		ports: newPorts(),
		lock:  new(sync.Mutex),
	}
}

func (tm *TaskManager) Close() {
	tm.lock.Lock()
	defer tm.lock.Unlock()
	for k, v := range tm.tasks {
		if err := v.Close(); nil != err {
			logger.Warnf("关闭任务[%s]异常-->%s", k, err.Error())
		}
	}
}

func (tm *TaskManager) Show(req udp.Request) udp.Response {
	if err := udp.ValidateOrder(req); nil != err {
		return udp.FailureResponse(err.Error())
	}
	tm.lock.Lock()
	defer tm.lock.Unlock()
	var temp map[string]Task
	if req.Type == "" {
		temp = tm.tasks
	} else if req.Name == "" {
		prefix := req.Type + "#"
		temp = make(map[string]Task)
		for k, v := range tm.tasks {
			if strings.Contains(k, prefix) {
				temp[k] = v
			}
		}
	} else {
		temp = make(map[string]Task, 1)
		sign := signature(req.Type, req.Name)
		for k, v := range tm.tasks {
			if k == sign {
				temp[k] = v
				break
			}
		}
	}
	if len(temp) <= 0 {
		return udp.FailureResponse("没有可用的数据")
	}
	data := ""
	if udp.OrderStatus == req.Order {
		count := 0
		total := make([]string, len(temp))
		for _, v := range temp {
			total[count] = showToStr(v.Show())
			count++
		}
		data = strings.Join(total, "\r\n")
	} else if udp.OrderNumbers == req.Order {
		for _, v := range temp {
			devNumbers := v.DevNumbers()
			length := len(devNumbers)
			builder := bytes.Buffer{}
			for i := 0; i < length; i++ {
				v1 := devNumbers[i]
				builder.WriteString(v1.Number)
				builder.WriteString(" ")
				builder.WriteString(v1.DevType)
				if i != length-1 {
					builder.WriteString("\r\n")
				}
			}
			data = builder.String()
			break
		}
	} else if udp.OrderParams == req.Order {
		for _, v := range temp {
			Params := udp.Params{Type: v.Show().Type, Name: v.Show().Name, Data: v.Params()}
			bs, err := json.Marshal(Params)
			if err != nil {
				return udp.FailureResponse("操作异常-->" + err.Error())
			}
			data = string(bs)
			break
		}
	} else if udp.OrderStatistic == req.Order {
		for _, v := range temp {
			data = v.Statistic().String()
			break
		}
	} else {
		return udp.FailureResponse("无法识别的命令: " + req.Order)
	}
	return udp.SuccessResponse(data)
}

func (tm *TaskManager) ShowForHttp(taskType, taskName string) []http.Show {
	tm.lock.Lock()
	defer tm.lock.Unlock()
	var temp map[string]Task
	if taskType == "" {
		temp = tm.tasks
	} else if taskName == "" {
		prefix := taskType + "#"
		temp = make(map[string]Task)
		for k, v := range tm.tasks {
			if strings.Contains(k, prefix) {
				temp[k] = v
			}
		}
	} else {
		temp = make(map[string]Task, 1)
		sign := signature(taskType, taskName)
		for k, v := range tm.tasks {
			if k == sign {
				temp[k] = v
				break
			}
		}
	}
	count := 0
	data := make([]http.Show, len(temp))
	for _, v := range temp {
		data[count] = v.Show()
		count++
	}
	return data
}

func (tm *TaskManager) DevNumbers(taskType, taskName string) []interact.DevNumber {
	tm.lock.Lock()
	defer tm.lock.Unlock()
	sign := signature(taskType, taskName)
	for k, v := range tm.tasks {
		if k == sign {
			return v.DevNumbers()
		}
	}
	return nil
}

func (tm *TaskManager) Params(taskType, taskName string) interface{} {
	tm.lock.Lock()
	defer tm.lock.Unlock()
	sign := signature(taskType, taskName)
	for k, v := range tm.tasks {
		if k == sign {
			return v.Params()
		}
	}
	return struct{}{}
}

func (tm *TaskManager) Statistic(taskType, taskName string) statistic.Statistic {
	tm.lock.Lock()
	defer tm.lock.Unlock()
	sign := signature(taskType, taskName)
	for k, v := range tm.tasks {
		if k == sign {
			return v.Statistic()
		}
	}
	return statistic.Statistic{}
}

func (tm *TaskManager) Add(req udp.Add) udp.Response {
	if err := udp.ValidateAdd(req); nil != err {
		return udp.FailureResponse(err.Error())
	}
	tm.lock.Lock()
	defer tm.lock.Unlock()
	if tm.isSignatureRepeated(req.TypeName) {
		return udp.FailureResponse("任务名称重复")
	}
	sign := signature(req.Type, req.Name)
	if !tm.ports.addPorts(sign, req.StartPort, req.PortCount) {
		return udp.FailureResponse("端口已被占用")
	}
	signature := signature(req.Type, req.Name)
	if coap.CoapType == req.Type {
		task := coap.NewTask(req.Name, req.LocalIP, req.StartPort, req.PortCount, req.IP, req.Port, nil)
		tm.tasks[signature] = task
	}
	return udp.SuccessResponse(nil)
}

func (tm *TaskManager) isSignatureRepeated(tn udp.TypeName) bool {
	signature := signature(tn.Type, tn.Name)
	for k := range tm.tasks {
		if k == signature {
			return true
		}
	}
	return false
}

func (tm *TaskManager) AddContent(req udp.AddContent) udp.Response {
	if err := udp.ValidateAddContent(req); nil != err {
		return udp.FailureResponse(err.Error())
	}
	tm.lock.Lock()
	defer tm.lock.Unlock()
	if tm.isSignatureRepeated(req.TypeName) {
		return udp.FailureResponse("任务名称重复")
	}
	signature := signature(req.Type, req.Name)
	if coap.CoapType == req.Type {
		portCount := len(req.DevNumbers)
		if portCount <= 0 {
			return udp.FailureResponse("缺少设备编号")
		}
		if portCount > 1000 {
			return udp.FailureResponse("一个任务最多占用1000个端口")
		}
		if !tm.ports.addPorts(signature, req.StartPort, portCount) {
			return udp.FailureResponse("端口已被占用")
		}
		task := coap.NewTask(req.Name, req.LocalIP, req.StartPort, portCount, req.IP, req.Port, req.DevNumbers)
		tm.tasks[signature] = task
	}
	return udp.SuccessResponse(nil)
}

func (tm *TaskManager) Order(req udp.Request) udp.Response {
	if err := udp.ValidateRequest(req); nil != err {
		return udp.FailureResponse(err.Error())
	}
	tm.lock.Lock()
	defer tm.lock.Unlock()
	signature := signature(req.Type, req.Name)
	defer func() {
		if udp.OrderClose == req.Order {
			tm.ports.removePorts(signature)
			delete(tm.tasks, signature)
		}
	}()
	if task, ok := tm.tasks[signature]; ok {
		if udp.OrderStart == req.Order {
			if err := task.Start(); nil != err {
				return udp.FailureResponse("启动任务失败: " + err.Error())
			}
		} else if udp.OrderClose == req.Order {
			if err := task.Close(); nil != err {
				return udp.FailureResponse("关闭任务异常: " + err.Error())
			}
		} else {
			return udp.FailureResponse("无法识别的命令: " + req.Order)
		}
	} else {
		return udp.FailureResponse(fmt.Sprintf("找不到名称为%s的%s任务", req.Name, req.Type))
	}
	return udp.SuccessResponse(nil)
}

func (tm *TaskManager) SetParams(params udp.Params) udp.Response {
	if v, err := udp.ValidateParams(params); err != nil {
		return udp.FailureResponse(err.Error())
	} else {
		tm.lock.Lock()
		defer tm.lock.Unlock()
		signature := signature(params.Type, params.Name)
		if task, ok := tm.tasks[signature]; ok {
			err = task.SetParams(v)
			if err != nil {
				return udp.FailureResponse(err.Error())
			}
		} else {
			return udp.FailureResponse(fmt.Sprintf("找不到名称为%s的%s任务", params.Name, params.Type))
		}
		return udp.SuccessResponse(nil)
	}
}

func signature(taskType, taskName string) string {
	return taskType + "#" + taskName
}

func showToStr(show http.Show) string {
	var builder strings.Builder
	//type
	builder.WriteString("[类型: ")
	builder.WriteString(show.Type)
	builder.WriteString(" | ")
	//name
	builder.WriteString("名称: ")
	builder.WriteString(show.Name)
	builder.WriteString(" | ")
	//status
	builder.WriteString("状态: ")
	builder.WriteString(show.State)
	builder.WriteString(" | ")
	//start port
	builder.WriteString("起始端口: ")
	builder.WriteString(strconv.Itoa(show.StartPort))
	builder.WriteString(" | ")
	//port count
	builder.WriteString("端口数量: ")
	builder.WriteString(strconv.Itoa(show.PortCount))
	builder.WriteString("]")
	return builder.String()
}
