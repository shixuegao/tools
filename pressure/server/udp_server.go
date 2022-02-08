package main

import (
	"encoding/json"
	"umx/tools/pressure/server/interact/udp"
	"umx/tools/pressure/server/util"
)

func openUdpServer(interact *udp.Interactor) {
	fWrapper := func(wrapper *udp.Wrapper) (data []byte) {
		defer func() {
			if ok, err := util.AssertErr(recover()); ok {
				logger.Errorf("处理客户端请求异常-->%s", err.Error())
				resp := udp.FailureResponse("系统错误")
				data, _ = json.Marshal(&resp)
			}
		}()
		path := string(wrapper.Path)
		if path == udp.PathTaskShow {
			req := udp.Request{}
			json.Unmarshal(wrapper.Content, &req)
			resp := taskManager.Show(req)
			data, _ = json.Marshal(resp)
		} else if path == udp.PathTaskAdd {
			add := udp.Add{}
			json.Unmarshal(wrapper.Content, &add)
			resp := taskManager.Add(add)
			data, _ = json.Marshal(resp)
		} else if path == udp.PathTaskAddContent {
			addContent := udp.AddContent{}
			json.Unmarshal(wrapper.Content, &addContent)
			resp := taskManager.AddContent(addContent)
			data, _ = json.Marshal(resp)
		} else if path == udp.PathTaskOrder {
			req := udp.Request{}
			json.Unmarshal(wrapper.Content, &req)
			resp := taskManager.Order(req)
			data, _ = json.Marshal(resp)
		} else if path == udp.PathTaskParams {
			params := udp.Params{}
			json.Unmarshal(wrapper.Content, &params)
			resp := taskManager.SetParams(params)
			data, _ = json.Marshal(resp)
		} else {
			resp := udp.FailureResponse("无效的路径")
			data, _ = json.Marshal(&resp)
		}
		return
	}
	fErr := func(err error, stack []byte) {
		logger.Errorf("UDP服务发生异常-->%s, 堆栈: %s", err.Error(), string(stack))
	}
	interact.StartAsync(fWrapper, fErr)
}
