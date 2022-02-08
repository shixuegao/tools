package udp

import (
	"errors"
	"umx/tools/pressure/server/util"
)

func ValidateOrder(r Request) (err error) {
	if r.Order == "" {
		err = errors.New("命令错误")
	}
	return
}

func ValidateRequest(r Request) (err error) {
	if r.Type == "" {
		err = errors.New("任务类型错误")
		return
	}
	if r.Name == "" {
		err = errors.New("任务名称错误")
		return
	}
	if r.Order == "" {
		err = errors.New("命令错误")
		return
	}
	return
}

func ValidateAdd(a Add) (err error) {
	if a.Type == "" {
		err = errors.New("任务类型错误")
		return
	}
	if a.Name == "" {
		err = errors.New("任务名称错误")
		return
	}
	if !util.IsLegalIpv4(a.IP) {
		err = errors.New("非法的目的IP")
		return
	}
	if a.Port < 1000 || a.Port > 65535 {
		err = errors.New("非法的目的端口")
		return
	}
	if !util.IsLegalPort(a.StartPort) {
		err = errors.New("非法的起始端口")
		return
	}
	if a.PortCount <= 0 {
		err = errors.New("非法的端口数量")
		return
	}
	if a.PortCount > 1000 {
		err = errors.New("一个任务最多占用1000个端口")
		return
	}
	maxPort := a.StartPort + a.PortCount - 1
	if maxPort > 65535 {
		err = errors.New("非法的端口数量")
		return
	}
	return
}

func ValidateAddContent(ac AddContent) (err error) {
	if ac.Type == "" {
		err = errors.New("任务类型错误")
		return
	}
	if ac.Name == "" {
		err = errors.New("任务名称错误")
		return
	}
	if !util.IsLegalPort(ac.StartPort) {
		err = errors.New("起始端口错误")
		return
	}
	if ac.DevNumbers == nil || len(ac.DevNumbers) <= 0 {
		err = errors.New("设备编码不能为空")
		return
	}
	return
}

func ValidateParams(params Params) (v map[string]interface{}, err error) {
	if params.Type == "" {
		err = errors.New("任务类型错误")
		return
	}
	if params.Name == "" {
		err = errors.New("任务名称错误")
		return
	}
	if vv, ok := params.Data.(map[string]interface{}); !ok {
		err = errors.New("没有可用的参数")
	} else {
		v = vv
	}
	return
}
