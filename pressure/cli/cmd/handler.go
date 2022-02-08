package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"umx/tools/pressure/cli/interact"
	coap "umx/tools/pressure/cli/interact/coap"

	"github.com/tealeg/xlsx"
)

func newConn() (*interact.Interactor, error) {
	addr := ServerIP + ":" + strconv.Itoa(ServerPort)
	c, err := interact.NewInteractor(addr)
	if nil != err {
		return nil, err
	}
	return c, nil
}

func handleTaskShow(req *interact.Request) {
	conn, err := newConn()
	if nil != err {
		fmt.Printf("建立UDP异常-->%s\r\n", err.Error())
		return
	}
	defer conn.Close()
	resp, err := conn.Exchange(interact.PathTaskShow, req)
	if nil != err {
		fmt.Printf("UDP通信异常-->%s", err.Error())
		return
	}
	if resp.Code == interact.Failure {
		fmt.Println(resp.Message)
		return
	}
	if Order == interact.OrderNumbers {
		if v, ok := resp.Data.(string); ok {
			devNumbers := toDevNumbers(v)
			writeToExcel(devNumbers)
		} else {
			fmt.Println("devNumbers不是字符串, 无法写文件")
		}
	} else if Order == interact.OrderParams {
		if v, ok := resp.Data.(string); ok {
			params := interact.Params{}
			err := json.Unmarshal([]byte(v), &params)
			if err != nil {
				fmt.Println("解析Json失败-->" + err.Error())
				return
			}
			handleTaskParams(params)
		}
		if v, ok := resp.Data.(interact.Params); ok {
			handleTaskParams(v)
		}
	} else {
		fmt.Println(resp.Data)
	}
}

func handleTaskParams(params interact.Params) {
	if params.Type == "" {
		return
	}
	if strings.ToLower(params.Type) == "coap" {
		if v, ok := params.Data.(map[string]interface{}); ok {
			cParams := coap.Params{}
			cParams.Heartbeat = int(v["Heartbeat"].(float64))
			cParams.DevState = int(v["DevState"].(float64))
			cParams.EventInfo = int(v["EventInfo"].(float64))
			cParams.Timeout = int(v["Timeout"].(float64))
			cParams.Lost = int(v["Lost"].(float64))
			fmt.Printf("[类型: %s | 名称: %s | 心跳周期: %d(s) | 设备测量值上报周期: %d(s) | 事件上报周期: %d(s) | 超时阈值: %d(ms) | 丢包阈值: %d(ms)]\n",
				params.Type, params.Name, cParams.Heartbeat, cParams.DevState, cParams.EventInfo, cParams.Timeout, cParams.Lost)
		}
	}
}

func handleTaskCoapParamsSet(taskType, taskName string, val coap.Params) {
	conn, err := newConn()
	if nil != err {
		fmt.Printf("建立UDP异常-->%s\r\n", err.Error())
		return
	}
	defer conn.Close()
	params := interact.Params{
		Type: taskType,
		Name: taskName,
		Data: val,
	}
	resp, err := conn.Exchange(interact.PathTaskParams, &params)
	if nil != err {
		fmt.Printf("UDP通信异常-->%s", err.Error())
		return
	}
	if resp.Code == interact.Success {
		fmt.Println("操作成功")
	} else {
		fmt.Println("操作失败-->" + resp.Message)
	}
}

func toDevNumbers(content string) []interact.DevNumber {
	if content == "" {
		return []interact.DevNumber{}
	}
	split := strings.Split(content, "\r\n")
	if len(split) <= 0 {
		return []interact.DevNumber{}
	}
	index := 0
	devNumbers := make([]interact.DevNumber, len(split))
	for _, str := range split {
		subSplit := strings.Split(str, " ")
		devNumber := interact.DevNumber{}
		devNumber.Number = subSplit[0]
		devNumber.DevType = subSplit[1]
		devNumbers[index] = devNumber
		index++
	}
	return devNumbers
}

//读Excel
func readFromExcel(filename string) ([]interact.DevNumber, error) {
	if !strings.HasSuffix(filename, ".xlsx") {
		return nil, errors.New("非法的文件类型")
	}
	file, err := xlsx.OpenFile(filename)
	if nil != err {
		return nil, err
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
	return devNumbers, nil
}

//写Excel
func writeToExcel(devNumbers []interact.DevNumber) {
	file := xlsx.NewFile()
	sheet, _ := file.AddSheet("sheet")
	for _, v := range devNumbers {
		row := sheet.AddRow()
		row.AddCell().Value = v.Number
		row.AddCell().Value = v.DevType
		row.AddCell().Value = v.SubDevType
	}
	filename := "Numbers_" + time.Now().Format("20060102150405") + ".xlsx"
	realFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 066)
	if nil != err {
		fmt.Printf("打开文件%s失败-->%s", filename, err.Error())
		return
	}
	err = file.Write(realFile)
	if nil != err {
		fmt.Printf("写文件%s失败-->%s", filename, err.Error())
	}
}

func handleTaskAdd(add *interact.Add) {
	conn, err := newConn()
	if nil != err {
		fmt.Printf("建立UDP异常-->%s\r\n", err.Error())
		return
	}
	defer conn.Close()
	resp, err := conn.Exchange(interact.PathTaskAdd, add)
	if nil != err {
		fmt.Printf("UDP通信异常-->%s", err.Error())
		return
	}
	if resp.Code == interact.Success {
		fmt.Println("操作成功")
	} else {
		fmt.Println("操作失败-->" + resp.Message)
	}
}

func handleTaskAddByFile(addContent *interact.AddContent) {
	conn, err := newConn()
	if nil != err {
		fmt.Printf("建立UDP异常-->%s\r\n", err.Error())
		return
	}
	defer conn.Close()
	resp, err := conn.Exchange(interact.PathTaskAddContent, addContent)
	if nil != err {
		fmt.Printf("UDP通信异常-->%s", err.Error())
		return
	}
	if resp.Code == interact.Success {
		fmt.Println("操作成功")
	} else {
		fmt.Println("操作失败-->" + resp.Message)
	}
}

func handleTaskOperation(req *interact.Request) {
	conn, err := newConn()
	if nil != err {
		fmt.Printf("建立UDP异常-->%s\r\n", err.Error())
		return
	}
	defer conn.Close()
	resp, err := conn.Exchange(interact.PathTaskOrder, req)
	if nil != err {
		fmt.Printf("UDP通信异常-->%s", err.Error())
		return
	}
	if resp.Code == interact.Success {
		fmt.Println("操作成功")
	} else {
		fmt.Println("操作失败-->" + resp.Message)
	}
}
