package coap

import (
	"sync"
	"umx/tools/pressure/server/task/statistic"
	"umx/tools/pressure/server/util"
)

type record struct {
	id          uint16
	recordTime  int64
	confirmTime int64
}

type coapStatistic struct {
	total   int64
	timeout int64
	lost    int64
	//delay
	count   int64
	delay   int64
	lock    *sync.Mutex
	records map[uint16]*record
}

func newCoapStatistic() *coapStatistic {
	return &coapStatistic{
		lock:    new(sync.Mutex),
		records: make(map[uint16]*record),
	}
}

//只记录Confirmable包(即发出去的包)
func (cs *coapStatistic) record(cType COAPType, id uint16) {
	if cType != Confirmable {
		return
	}
	rec := record{id: id, recordTime: util.NowMillisecond()}
	cs.lock.Lock()
	defer cs.lock.Unlock()
	cs.total++
	cs.records[id] = &rec
}

//只确认Acknowledgement包(即收到的回包)
func (cs *coapStatistic) confirm(cType COAPType, id uint16) {
	if cType != Acknowledgement {
		return
	}
	cs.lock.Lock()
	defer cs.lock.Unlock()
	rec := cs.records[id]
	if rec != nil {
		rec.confirmTime = util.NowMillisecond()
	}
}

//清理统计信息
func (cs *coapStatistic) clearStatistic(timeout, loss int64) statistic.Statistic {
	cs.lock.Lock()
	defer cs.lock.Unlock()
	var count int64
	var delay int64
	discard := 2 * loss
	now := util.NowMillisecond()
	for k, rec := range cs.records {
		//已确认
		if rec.confirmTime != 0 {
			count++
			delay += rec.confirmTime - rec.recordTime
			cs.adjustOffset(rec.confirmTime-rec.recordTime, timeout, loss)
			delete(cs.records, k)
			continue
		}
		//未确认且滞留时间超过2倍于丢失时间, 则废弃
		if now-rec.recordTime >= discard {
			cs.adjustOffset(discard, timeout, loss)
			delete(cs.records, k)
		}
	}
	//计算平均回包延迟
	if count > 0 {
		delay += cs.count * cs.delay
		cs.count += count
		delay /= cs.count
		cs.delay = delay
	}
	return statistic.Statistic{
		Total:   cs.total,
		Delay:   cs.delay,
		Timeout: cs.timeout,
		Lost:    cs.lost,
	}
}

func (cs *coapStatistic) adjustOffset(offset, timeout, loss int64) {
	if offset >= timeout && offset < loss {
		cs.timeout++
	} else if offset >= loss {
		cs.lost++
	}
}
