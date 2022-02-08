package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strings"
	"time"
	"umx/tools/pressure/server/interact"
	"umx/tools/pressure/server/interact/http"
	"umx/tools/pressure/server/interact/udp"
	"umx/tools/pressure/server/util"

	"github.com/tealeg/xlsx"
)

const (
	defaultUsername = "admin"
	defaultPassword = "admin"
)

var server *http.Server
var token *tokenMap

func openHttpServer(addr string) (err error) {
	server, paths, err := http.NewServer(addr, defaultHandlers())
	if nil != err {
		return
	}
	//print paths
	logger.Info("Urls: " + strings.Join(paths, ", "))
	token = newTokenMap()
	fErr := func(err error) {
		logger.Errorf("HTTP服务发生异常-->%s", err.Error())
	}
	server.StartAsync(fErr)
	return
}

func closeHttpServer() error {
	if server != nil {
		return server.Close()
	}
	return nil
}

func defaultHandlers() map[string]http.Handler {
	handlers := make(map[string]http.Handler)
	handlers["/"] = &MyHandler{}
	return handlers
}

type MyHandler struct{}

// /
func (mh *MyHandler) Get(iw *http.InnerWrapper) {
	ok := accessNoWriteBack(iw)
	next := ""
	if ok {
		next = "/static/index.html"
	} else {
		next = "/static/login.html"
	}
	iw.Redirection(next)
}

// /login
func (mh *MyHandler) GetLogin(iw *http.InnerWrapper) {
	iw.FormParser()
	username := iw.ParamString("username")
	password := iw.ParamString("password")
	if username == "" || password == "" {
		iw.WriteToJson(udp.FailureResponse("账号或者密码错误"))
	}
	if username == defaultUsername && password == defaultPassword {
		uuid := token.add()
		token.clear()
		iw.WriteToJson(udp.SuccessResponse(uuid))
	} else {
		iw.WriteToJson(udp.FailureResponse("账号或者密码错误"))
	}
}

// /static
func (mh *MyHandler) GetStatic(iw *http.InnerWrapper) {
	filePath := iw.RawRequest().URL.Path[1:]
	writeFile(filePath, iw)
}

// /task/show
func (mh *MyHandler) GetTaskShow(iw *http.InnerWrapper) {
	if !access(iw) {
		return
	}
	iw.FormParser()
	taskType := iw.ParamString("type")
	taskName := iw.ParamString("name")
	data := taskManager.ShowForHttp(taskType, taskName)
	iw.WriteToJson(udp.SuccessResponse(data))
}

// /task/show/numbers
func (mh *MyHandler) GetTaskShowNumbers(iw *http.InnerWrapper) {
	if !access(iw) {
		return
	}
	iw.FormParser()
	taskType := iw.ParamString("type")
	taskName := iw.ParamString("name")
	devNumbers := taskManager.DevNumbers(taskType, taskName)
	file := xlsx.NewFile()
	sheet, _ := file.AddSheet("sheet")
	for _, v := range devNumbers {
		row := sheet.AddRow()
		row.AddCell().Value = v.Number
		row.AddCell().Value = v.DevType
		row.AddCell().Value = v.SubDevType
	}
	buffer := new(bytes.Buffer)
	file.Write(buffer)
	filename := "Numbers_" + time.Now().Format("20060102150405") + ".xlsx"
	contentType := "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	err := util.WriteToHttpResponse(filename, contentType, buffer.Bytes(), iw.RawResponse())
	if nil != err {
		logger.Warnf("写Excel文件到HTTP异常-->%s", err.Error())
	}
}

// /task/show/params
func (mh *MyHandler) GetTaskShowParams(iw *http.InnerWrapper) {
	if !access(iw) {
		return
	}
	iw.FormParser()
	taskType := iw.ParamString("type")
	taskName := iw.ParamString("name")
	data := taskManager.Params(taskType, taskName)
	iw.WriteToJson(udp.SuccessResponse(data))
}

// /task/show/statistic
func (mh *MyHandler) GetTaskShowStatistic(iw *http.InnerWrapper) {
	if !access(iw) {
		return
	}
	iw.FormParser()
	taskType := iw.ParamString("type")
	taskName := iw.ParamString("name")
	data := taskManager.Statistic(taskType, taskName)
	iw.WriteToJson(udp.SuccessResponse(data))
}

// /task/add
func (mh *MyHandler) PostTaskAdd(iw *http.InnerWrapper) {
	if !access(iw) {
		return
	}
	iw.FormParser()
	add := udp.Add{
		TypeName: udp.TypeName{
			Type: iw.ParamString("type"),
			Name: iw.ParamString("name"),
		},
		LocalIP:   iw.ParamString("localIP"),
		IP:        iw.ParamString("ip"),
		Port:      iw.ParamInt("port"),
		StartPort: iw.ParamInt("startPort"),
		PortCount: iw.ParamInt("portCount"),
	}
	resp := taskManager.Add(add)
	iw.WriteToJson(resp)
}

// /task/add/file
func (mh *MyHandler) PostTaskAddFile(iw *http.InnerWrapper) {
	if !access(iw) {
		return
	}
	iw.FormParser()
	content := udp.AddContent{
		TypeName: udp.TypeName{
			Type: iw.ParamString("type"),
			Name: iw.ParamString("name"),
		},
		LocalIP:   iw.ParamString("localIP"),
		IP:        iw.ParamString("ip"),
		Port:      iw.ParamInt("port"),
		StartPort: iw.ParamInt("startPort"),
	}
	data, err := iw.FileContent("file")
	if nil != err {
		resp := udp.FailureResponse("无法读取设备序列号")
		logger.Warnf("无法读取设备序列号-->%s", err.Error())
		iw.WriteToJson(resp)
		return
	}
	file, err := xlsx.OpenBinary(data)
	if nil != err {
		resp := udp.FailureResponse("无法转成xlsx文件")
		logger.Warnf("无法转成xlsx文件-->%s", err.Error())
		iw.WriteToJson(resp)
		return
	}
	sheet := file.Sheet["sheet"]
	devNumbers := make([]interact.DevNumber, sheet.MaxRow)
	for i := 0; i < sheet.MaxRow; i++ {
		row := sheet.Row(i)
		devNumber := interact.DevNumber{
			Number:     row.Cells[0].Value,
			DevType:    row.Cells[1].Value,
			SubDevType: row.Cells[2].Value,
		}
		devNumbers[i] = devNumber
	}
	content.DevNumbers = devNumbers
	resp := taskManager.AddContent(content)
	iw.WriteToJson(resp)
}

// /task/operation
func (mh *MyHandler) PutTaskOperation(iw *http.InnerWrapper) {
	if !access(iw) {
		return
	}
	iw.FormParser()
	req := udp.Request{
		TypeName: udp.TypeName{
			Type: iw.ParamString("type"),
			Name: iw.ParamString("name"),
		},
		Order: iw.ParamString("order"),
	}
	resp := taskManager.Order(req)
	iw.WriteToJson(resp)
}

// /task/params
func (mh *MyHandler) PutTaskParams(iw *http.InnerWrapper) {
	if !access(iw) {
		return
	}
	iw.FormParser()
	bs, err := iw.Body()
	if err != nil {
		resp := udp.FailureResponse("读取Body数据失败-->" + err.Error())
		iw.WriteToJson(resp)
		return
	}
	params := udp.Params{}
	err = json.Unmarshal(bs, &params)
	if err != nil {
		resp := udp.FailureResponse("Body转Json失败-->" + err.Error())
		iw.WriteToJson(resp)
		return
	}
	resp := taskManager.SetParams(params)
	iw.WriteToJson(resp)
}

//接入判断
func access(iw *http.InnerWrapper) bool {
	auth := iw.Header("Authorization")
	ok := token.has(auth)
	if !ok {
		data := udp.FailureResponse("无效的Token")
		iw.WriteToJson(data)
	}
	return ok
}

func accessNoWriteBack(iw *http.InnerWrapper) bool {
	auth := iw.Header("Authorization")
	return token.has(auth)
}

//写文件到HTTP
func writeFile(filePath string, iw *http.InnerWrapper) {
	c, err := ioutil.ReadFile(filePath)
	if nil != err {
		logger.Errorf("读取文件%s失败-->%s", filePath, err.Error())
		iw.RawResponse().WriteHeader(404)
	} else {
		if strings.HasSuffix(filePath, ".css") {
			iw.RawResponse().Header().Set("Content-Type", "text/css")
		} else if strings.HasSuffix(filePath, ".html") {
			iw.RawResponse().Header().Set("Content-Type", "text/html")
		} else if strings.HasSuffix(filePath, ".js") {
			iw.RawResponse().Header().Set("Content-Type", "text/js")
		}
		_, err = iw.RawResponse().Write(c)
		if nil != err {
			logger.Errorf("写文件到HTTP失败-->%s", err.Error())
			return
		}
	}
}
