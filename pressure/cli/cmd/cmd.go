package cmd

import (
	"fmt"
	"umx/tools/pressure/cli/interact"
	"umx/tools/pressure/cli/interact/coap"

	"github.com/spf13/cobra"
)

var (
	RootCmd       *cobra.Command
	taskCmd       *cobra.Command
	showCmd       *cobra.Command
	addCmd        *cobra.Command
	addFileCmd    *cobra.Command
	operationCmd  *cobra.Command
	paramsCmd     *cobra.Command
	paramsCoapCmd *cobra.Command
)

var (
	ServerIP   string = "localhost"
	ServerPort int    = 65000

	Type    string = ""
	Name    string = ""
	Order   string = ""
	LocalIP string = ""

	IP        string = ""
	Port      int    = 0
	StartPort int    = 0
	PortCount int    = 0

	Filename string = ""
)

var (
	CoapHeartbeat int = 30
	CoapDevState  int = 120
	CoapEventInfo int = 300
	CoapTimeout   int = 100
	CoapLost      int = 500
)

func init() {
	initCmd()
}

func initCmd() {
	RootCmd = &cobra.Command{
		Use: "cli",
	}
	taskCmd = &cobra.Command{
		Use:   "task",
		Short: "执行任务相关操作",
	}
	showCmd = &cobra.Command{
		Use:   "show",
		Short: "执行查看任务操作",
		Run: func(cmd *cobra.Command, args []string) {
			request := interact.Request{
				TypeName: interact.TypeName{
					Type: Type,
					Name: Name,
				},
				Order: Order,
			}
			handleTaskShow(&request)
		},
	}
	addCmd = &cobra.Command{
		Use:   "add",
		Short: "执行添加任务操作",
		Run: func(cmd *cobra.Command, args []string) {
			add := interact.Add{
				TypeName: interact.TypeName{
					Type: Type,
					Name: Name,
				},
				LocalIP:   LocalIP,
				IP:        IP,
				Port:      Port,
				StartPort: StartPort,
				PortCount: PortCount,
			}
			handleTaskAdd(&add)
		},
	}
	addFileCmd = &cobra.Command{
		Use:   "addFile",
		Short: "通过文件添加任务",
		Run: func(cmd *cobra.Command, args []string) {
			devNumbers, err := readFromExcel(Filename)
			if nil != err {
				fmt.Println(err.Error())
				return
			}
			addContent := interact.AddContent{
				TypeName: interact.TypeName{
					Type: Type,
					Name: Name,
				},
				LocalIP:    LocalIP,
				IP:         IP,
				Port:       Port,
				StartPort:  StartPort,
				DevNumbers: devNumbers,
			}
			handleTaskAddByFile(&addContent)
		},
	}
	operationCmd = &cobra.Command{
		Use:   "order",
		Short: "执行启动、关闭任务操作",
		Run: func(cmd *cobra.Command, args []string) {
			req := interact.Request{
				TypeName: interact.TypeName{
					Type: Type,
					Name: Name,
				},
				Order: Order,
			}
			handleTaskOperation(&req)
		},
	}
	paramsCmd = &cobra.Command{
		Use:   "params",
		Short: "任务参数变更",
	}
	paramsCoapCmd = &cobra.Command{
		Use:   "coap",
		Short: "coap参数变更",
		Run: func(cmd *cobra.Command, args []string) {
			val := coap.Params{
				Heartbeat: CoapHeartbeat,
				DevState:  CoapDevState,
				EventInfo: CoapEventInfo,
				Timeout:   CoapTimeout,
				Lost:      CoapLost,
			}
			if val.Heartbeat <= 0 {
				fmt.Println("非法的心跳上报周期")
				return
			}
			if val.DevState <= 0 {
				fmt.Println("非法的测量值上报周期")
				return
			}
			if val.EventInfo <= 0 {
				fmt.Println("非法的事件上报周期")
				return
			}
			handleTaskCoapParamsSet("coap", Name, val)
		},
	}

	RootCmd.AddCommand(taskCmd)
	RootCmd.PersistentFlags().StringVarP(&ServerIP, "serverIP", "v", "localhost", "服务IP")
	RootCmd.PersistentFlags().IntVarP(&ServerPort, "serverPort", "s", 65000, "服务端口")
	taskCmd.AddCommand(showCmd, addCmd, addFileCmd, operationCmd, paramsCmd)
	//show
	showCmd.Flags().StringVarP(&Type, "type", "t", "", "任务类型(coap)")
	showCmd.Flags().StringVarP(&Name, "name", "n", "", "任务名称")
	showCmd.Flags().StringVarP(&Order, "order", "o", "", "操作类型(status, numbers, params, statistic)")
	showCmd.MarkFlagRequired("order")
	//add
	addCmd.Flags().StringVarP(&Type, "type", "t", "", "任务类型(coap)")
	addCmd.Flags().StringVarP(&Name, "name", "n", "", "任务名称")
	addCmd.Flags().StringVarP(&LocalIP, "localIP", "l", "", "本地IP")
	addCmd.Flags().StringVarP(&IP, "ip", "i", "", "目的IP")
	addCmd.Flags().IntVarP(&Port, "port", "p", 0, "目的端口")
	addCmd.Flags().IntVarP(&StartPort, "startPort", "o", 0, "任务起始端口")
	addCmd.Flags().IntVarP(&PortCount, "portCount", "c", 0, "任务端口数量")
	addCmd.MarkFlagRequired("type")
	addCmd.MarkFlagRequired("name")
	addCmd.MarkFlagRequired("localIP")
	addCmd.MarkFlagRequired("ip")
	addCmd.MarkFlagRequired("port")
	addCmd.MarkFlagRequired("startPort")
	addCmd.MarkFlagRequired("portCount")
	//addFile
	addFileCmd.Flags().StringVarP(&Type, "type", "t", "", "任务类型(coap)")
	addFileCmd.Flags().StringVarP(&Name, "name", "n", "", "任务名称")
	addFileCmd.Flags().StringVarP(&LocalIP, "localIP", "l", "", "本地IP")
	addFileCmd.Flags().StringVarP(&IP, "ip", "i", "", "目的IP")
	addFileCmd.Flags().IntVarP(&Port, "port", "p", 0, "目的端口")
	addFileCmd.Flags().IntVarP(&StartPort, "startPort", "o", 0, "任务起始端口")
	addFileCmd.Flags().StringVarP(&Filename, "file", "f", "", "文件名称")
	addFileCmd.MarkFlagRequired("type")
	addFileCmd.MarkFlagRequired("name")
	addFileCmd.MarkFlagRequired("localIP")
	addFileCmd.MarkFlagRequired("ip")
	addFileCmd.MarkFlagRequired("port")
	addFileCmd.MarkFlagRequired("startPort")
	addFileCmd.MarkFlagRequired("file")
	//operation
	operationCmd.Flags().StringVarP(&Type, "type", "t", "", "任务类型(coap)")
	operationCmd.Flags().StringVarP(&Name, "name", "n", "", "任务名称")
	operationCmd.Flags().StringVarP(&Order, "order", "o", "", "命令(start, close)")
	operationCmd.MarkFlagRequired("type")
	operationCmd.MarkFlagRequired("name")
	operationCmd.MarkFlagRequired("order")
	//params
	paramsCmd.AddCommand(paramsCoapCmd)
	paramsCoapCmd.Flags().StringVarP(&Name, "name", "n", "", "任务名称")
	paramsCoapCmd.Flags().IntVarP(&CoapHeartbeat, "heartbeat", "e", 30, "心跳周期(秒)")
	paramsCoapCmd.Flags().IntVarP(&CoapDevState, "devState", "d", 120, "测量值与状态上报周期(秒)")
	paramsCoapCmd.Flags().IntVarP(&CoapEventInfo, "eventInfo", "o", 300, "事件上报周期(秒)")
	paramsCoapCmd.Flags().IntVarP(&CoapTimeout, "timeout", "m", 100, "超时阈值(毫秒)")
	paramsCoapCmd.Flags().IntVarP(&CoapLost, "lost", "l", 500, "丢包阈值(毫秒)")
}
