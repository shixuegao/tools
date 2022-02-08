package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"
	"umx/tools/pressure/server/interact/udp"
	"umx/tools/pressure/server/log"
	"umx/tools/pressure/server/task"
	"umx/tools/pressure/server/util"
)

var logger = log.Logger
var taskManager *task.TaskManager

func init() {
	taskManager = task.NewTaskManager()
}

func main() {
	defer log.Clear()
	var ip string = "localhost"
	var udpPort int = 65000
	var httpPort int = 65001
	var printMemState bool = false
	flag.StringVar(&ip, "i", "localhost", "IP地址")
	flag.IntVar(&udpPort, "p", 65000, "UDP通信端口, 默认65000")
	flag.IntVar(&httpPort, "h", 65001, "HTTP服务端口, 默认65001")
	flag.BoolVar(&printMemState, "m", false, "是否打印内存信息, 默认关闭")
	flag.Parse()
	if !util.IsLegalPort(udpPort) {
		logger.Error("非法的UDP通信端口")
		os.Exit(-1)
	}
	if !util.IsLegalPort(httpPort) {
		logger.Error("非法的HTTP服务端口")
		os.Exit(-1)
	}
	//create udp server
	addr := ip + ":" + strconv.Itoa(udpPort)
	interact, err := udp.NewInteractor(addr)
	if nil != err {
		logger.Errorf("新建UDP服务失败-->%s", err.Error())
		os.Exit(-1)
	}
	//create and run http server
	addr = ip + ":" + strconv.Itoa(httpPort)
	err = openHttpServer(addr)
	if nil != err {
		logger.Errorf("HTTP服务启动失败-->%s", err.Error())
		os.Exit(-1)
	}
	logger.Infof("HTTP服务启动成功, IP: %s, 端口: %d", ip, httpPort)
	//run udp server
	openUdpServer(interact)
	logger.Infof("UDP服务启动成功, IP: %s, 端口: %d", ip, udpPort)
	//print mem state
	wait, exit := intervalPrintMemoryState(printMemState)
	//signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	<-signalChan
	//close http server
	err = closeHttpServer()
	if nil != err {
		logger.Warnf("关闭HTTP服务异常-->%s", err.Error())
	}
	//close udp server
	err = interact.Close()
	if nil != err {
		logger.Warnf("关闭UDP服务异常-->%s", err.Error())
	}
	//close print
	if wait != nil && exit != nil {
		close(exit)
		wait.Wait()
	}
}

func intervalPrintMemoryState(enable bool) (*sync.WaitGroup, chan byte) {
	if !enable {
		return nil, nil
	}
	fmt.Println("开启内存信息打印...")
	exit := make(chan byte)
	wait := new(sync.WaitGroup)
	wait.Add(1)
	go func() {
		defer func() {
			wait.Done()
			fmt.Println("内存信息打印关闭...")
		}()
		var ms runtime.MemStats
		for {
			//sleep 10s
			for i := 0; i < 100; i++ {
				select {
				case <-exit:
					return
				default:
				}
				time.Sleep(100 * time.Millisecond)
			}
			//print
			runtime.ReadMemStats(&ms)
			bsCount := ms.Alloc
			kbCount := bsCount / 1024
			fmt.Printf("占用内存: %dKb\n", kbCount)
		}
	}()
	return wait, exit
}
