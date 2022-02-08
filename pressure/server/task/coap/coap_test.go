package coap_test

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"umx/tools/pressure/server/interact"
	"umx/tools/pressure/server/protocol"
	"umx/tools/pressure/server/task/coap"
)

func TestV30(t *testing.T) {
	devNumbers := []interact.DevNumber{
		{
			Number:  "0000000518110853",
			DevType: protocol.GetDevTypeName(protocol.DevTypeSynthetic),
		},
	}
	coapTask := coap.NewTask("test1", "localhost", 10000, 1, "localhost", 5683, devNumbers)
	err := coapTask.Start()
	if nil != err {
		fmt.Println("启动失败-->", err.Error())
		return
	}
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	<-signalChan
	err = coapTask.Close()
	if nil != err {
		fmt.Println("关闭异常-->", err.Error())
	}
}
