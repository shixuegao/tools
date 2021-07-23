package subscribe

import (
	"fmt"
	"sync"
)

type MemoryBitMap struct {
	limit  uint64
	offset uint64
	p      []byte
	lock   *sync.RWMutex
}

func (mbm *MemoryBitMap) Offset() uint64 {
	return mbm.offset
}

func (mbm *MemoryBitMap) Existed(value uint64) bool {
	if v, ok := normalized(value, mbm.offset); !ok {
		return false
	} else {
		mbm.lock.RLock()
		defer mbm.lock.RUnlock()
		if v > mbm.limit {
			return false
		}
		index, mask := indexAndMask(v)
		b := mbm.p[index]
		return (b & mask) > 0
	}
}

func (mbm *MemoryBitMap) Put(value uint64) error {
	return mbm.PutAndResize(value, false)
}

func (mbm *MemoryBitMap) Remove(value uint64) {
	if v, ok := normalized(value, mbm.offset); !ok {
		return
	} else {
		mbm.lock.Lock()
		defer mbm.lock.Unlock()
		if v > mbm.limit {
			return
		}
		index, mask := indexAndMask(v)
		mbm.p[index] = mbm.p[index] & (^mask)
	}
}

func (mbm *MemoryBitMap) Size() int {
	mbm.lock.RLock()
	defer mbm.lock.RUnlock()
	return len(mbm.p)
}

//from bitcount caculation of redis
func (mbm *MemoryBitMap) Count() int {
	mbm.lock.RLock()
	defer mbm.lock.RUnlock()
	return caculateBitCount(mbm.p)
}

func (mbm *MemoryBitMap) PutAndResize(value uint64, resize bool) error {
	if v, ok := normalized(value, mbm.offset); !ok {
		return fmt.Errorf("非法的值[%d], 值必须在[%d, %d]之间", value, mbm.offset+1, mbm.offset+LimitMaxValue)
	} else {
		mbm.lock.Lock()
		defer mbm.lock.Unlock()
		if !resize && v > mbm.limit {
			return fmt.Errorf("值[%d]大于当前最大值[%d], 请先扩容再插入值", value, mbm.limit+mbm.offset)
		}
		if resize && v > mbm.limit {
			mbm.changeCapacity(v)
		}
		index, mask := indexAndMask(v)
		mbm.p[index] = mbm.p[index] | mask
		return nil
	}
}

func (mbm *MemoryBitMap) Resize(value uint64) error {
	if v, ok := normalized(value, mbm.offset); !ok {
		return fmt.Errorf("非法的值[%d], 值必须在[%d, %d]之间", value, mbm.offset+1, mbm.offset+LimitMaxValue)
	} else {
		mbm.lock.Lock()
		defer mbm.lock.Unlock()
		mbm.changeCapacity(v)
		return nil
	}
}

func (mbm *MemoryBitMap) Range() (min, max uint64) {
	mbm.lock.RLock()
	defer mbm.lock.RUnlock()
	return 1 + mbm.offset, mbm.limit + mbm.offset
}

func (mbm *MemoryBitMap) changeCapacity(max uint64) {
	if max == mbm.limit {
		return
	}
	if mbm.limit == 0 {
		index, _ := indexAndMask(max)
		mbm.p = make([]byte, index+1)
	} else {
		oldIndex, _ := indexAndMask(mbm.limit)
		newIndex, _ := indexAndMask(max)
		if oldIndex > newIndex {
			mbm.p = mbm.p[:newIndex+1]
		} else {
			appended := make([]byte, newIndex-oldIndex)
			mbm.p = append(mbm.p, appended...)
		}
	}
	mbm.limit = max
}

func NewMemoryBitMap(max, offset uint64) (*MemoryBitMap, error) {
	if v, ok := normalized(max, offset); !ok {
		return nil, fmt.Errorf("非法的最大值[%d], 值必须在[%d, %d]之间", max, offset+1, offset+2^32)
	} else {
		mbm := MemoryBitMap{
			offset: offset,
			lock:   new(sync.RWMutex),
		}
		mbm.changeCapacity(v)
		return &mbm, nil
	}
}
