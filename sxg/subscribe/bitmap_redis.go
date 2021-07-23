package subscribe

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
)

const (
	RO_AUTH     = "AUTH"
	RO_SELECT   = "SELECT"
	RO_GETBIT   = "GETBIT"
	RO_SETBIT   = "SETBIT"
	RO_BITCOUNT = "BITCOUNT"
)

type RedisParams struct {
	Addr           string
	Password       string
	Db             int
	ConnectTimeout time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
}

func (rp RedisParams) Validate() error {
	if rp.Addr == "" {
		return errors.New("服务地址不能为空")
	}
	if rp.Db < 0 {
		return errors.New("DB必须大于等于0")
	}
	if rp.ConnectTimeout < 0 {
		rp.ConnectTimeout = 30 * time.Millisecond
	}
	if rp.ReadTimeout < 0 {
		rp.ReadTimeout = 20 * time.Millisecond
	}
	if rp.WriteTimeout < 0 {
		rp.WriteTimeout = 20 * time.Millisecond
	}
	return nil
}

func (rp RedisParams) String() string {
	return fmt.Sprintf("服务地址: %v, DB: %v, ConnectionTimeout: %v, ReadTimeout: %v, WriteTimeout: %v",
		rp.Addr, rp.Db, rp.ConnectTimeout, rp.ReadTimeout, rp.WriteTimeout)
}

func DefaultRedisParams(addr string) RedisParams {
	return RedisParams{
		Addr:           addr,
		Password:       "",
		Db:             0,
		ConnectTimeout: 30 * time.Second,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
	}
}

type RedisBitMap struct {
	rp     RedisParams
	f      func() redis.Conn
	pool   *redis.Pool
	key    string
	limit  uint64
	offset uint64
	lock   *sync.RWMutex
}

func (rbm *RedisBitMap) getConn() redis.Conn {
	var conn redis.Conn = nil
	if rbm.f != nil {
		conn = rbm.f()
	} else if rbm.pool != nil {
		conn = rbm.pool.Get()
	}
	if conn != nil {
		return conn
	}
	panic(errors.New("无法获取Redis连接"))
}

func (rbm *RedisBitMap) Params() string {
	return rbm.rp.String()
}

func (rbm *RedisBitMap) Close() error {
	if rbm.pool != nil {
		return rbm.pool.Close()
	}
	return nil
}

func (rbm *RedisBitMap) Existed(value uint64) bool {
	if v, ok := normalized(value, rbm.offset); !ok {
		return false
	} else {
		rbm.lock.RLock()
		defer rbm.lock.RUnlock()
		if v > rbm.limit {
			return false
		}
		v, err := rbm.getConn().Do(RO_GETBIT, rbm.key, v-1)
		if err != nil {
			return false
		}
		return AssertInt(v) == 1
	}
}

func (rbm *RedisBitMap) Put(value uint64) error {
	if v, ok := normalized(value, rbm.offset); !ok {
		return fmt.Errorf("非法的值[%d], 值必须在[%d, %d]之间", value, rbm.offset+1, rbm.offset+LimitMaxValue)
	} else {
		rbm.lock.Lock()
		defer rbm.lock.Unlock()
		if v > rbm.limit {
			return fmt.Errorf("值[%d]大于当前最大值[%d], 请先扩容再插入值", value, rbm.limit+rbm.offset)
		}
		if _, err := rbm.getConn().Do(RO_SETBIT, rbm.key, v-1, 1); err != nil {
			return err
		}
	}
	return nil
}

func (rbm *RedisBitMap) Remove(value uint64) {
	if v, ok := normalized(value, rbm.offset); ok {
		rbm.lock.Lock()
		defer rbm.lock.Unlock()
		if v > rbm.limit {
			return
		}
		rbm.getConn().Do(RO_SETBIT, rbm.key, v-1, 0)
	}
}

//only can increase the capacity
func (rbm *RedisBitMap) Resize(value uint64) error {
	if v, ok := normalized(value, rbm.offset); !ok {
		return fmt.Errorf("非法的值[%d], 值必须在[%d, %d]之间", value, rbm.offset+1, rbm.offset+LimitMaxValue)
	} else {
		rbm.lock.Lock()
		defer rbm.lock.Unlock()
		if rbm.limit < v {
			_, err := rbm.getConn().Do(RO_SETBIT, rbm.key, v-1, 0)
			if err != nil {
				return err
			}
			rbm.limit = v
		}
		return nil
	}
}

func (rbm *RedisBitMap) Range() (min, max uint64) {
	rbm.lock.RLock()
	defer rbm.lock.RUnlock()
	return 1 + rbm.offset, rbm.limit + rbm.offset
}

func (rbm *RedisBitMap) Size() int {
	rbm.lock.RLock()
	defer rbm.lock.RUnlock()
	index, _ := indexAndMask(rbm.limit)
	return int(index) + 1
}

func (rbm *RedisBitMap) Count() int {
	if v, err := rbm.getConn().Do(RO_BITCOUNT, rbm.key); err != nil {
		return -1
	} else {
		return AssertInt(v)
	}
}

func NewRedisBitMap(key string, max, offset uint64, rp RedisParams) (*RedisBitMap, error) {
	if key == "" {
		return nil, errors.New("键值不能为空")
	}
	value, ok := normalized(max, offset)
	if !ok {
		return nil, fmt.Errorf("最大值[%d]必须在[%d]到[%d]之间", max, offset+1, offset+LimitMaxValue)
	}
	if err := rp.Validate(); err != nil {
		return nil, err
	}
	bitMap := RedisBitMap{
		rp:     rp,
		key:    key,
		limit:  value,
		offset: offset,
		lock:   new(sync.RWMutex),
		pool: &redis.Pool{
			Dial: func() (redis.Conn, error) {
				conn, err := redis.Dial("tcp", rp.Addr, redis.DialConnectTimeout(rp.ConnectTimeout),
					redis.DialReadTimeout(rp.ReadTimeout), redis.DialWriteTimeout(rp.WriteTimeout))
				if err != nil {
					return nil, err
				}
				if rp.Password != "" {
					_, err = conn.Do(RO_AUTH, rp.Password)
					if err != nil {
						conn.Close()
						return nil, err
					}
				}
				_, err = conn.Do(RO_SELECT, rp.Db)
				if err != nil {
					conn.Close()
					return nil, err
				}
				return conn, nil
			},
		},
	}
	return &bitMap, nil
}

func NewRedisBitMapWithFunc(key string, max, offset uint64, f func() redis.Conn) (*RedisBitMap, error) {
	if key == "" {
		return nil, errors.New("键值不能为空")
	}
	if f == nil {
		return nil, errors.New("Redis连接获取函数不能为空")
	}
	value, ok := normalized(max, offset)
	if !ok {
		return nil, fmt.Errorf("最大值[%d]必须在[%d]到[%d]之间", max, offset+1, offset+LimitMaxValue)
	}
	bitMap := RedisBitMap{
		key:    key,
		limit:  value,
		offset: offset,
		lock:   new(sync.RWMutex),
		f:      f,
	}
	return &bitMap, nil
}
