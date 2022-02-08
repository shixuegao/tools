package main

import (
	"sync"
	"time"
	"umx/tools/pressure/server/util"
)

const (
	//tokenTimeout token持续时间 2小时
	tokenTimeout = 2 * 60 * 60 * 1000
)

type tokenMap struct {
	pMap *sync.Map
}

func newTokenMap() *tokenMap {
	return &tokenMap{
		pMap: new(sync.Map),
	}
}

func (tm *tokenMap) has(key string) bool {
	if key == "" {
		return false
	}
	_, ok := tm.pMap.Load(key)
	return ok
}

func (tm *tokenMap) add() string {
	uuid := util.UUID()
	timeout := time.Now().Nanosecond()/1e6 + tokenTimeout
	tm.pMap.Store(uuid, timeout)
	return uuid
}

func (tm *tokenMap) clear() {
	cur := time.Now().Nanosecond() / 1e6
	tm.pMap.Range(func(k interface{}, v interface{}) bool {
		timeout := v.(int)
		if cur > timeout {
			tm.pMap.Delete(k)
		}
		return true
	})
}
