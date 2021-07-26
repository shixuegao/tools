package subscribe

import (
	"strconv"
	"sync/atomic"
	"time"
)

const (
	PrefixRedisKey = "PreferredVideo:"
)

var userID uint64 = 0
var videoID uint64 = 0

func PreferredRedisKey(userId uint64) string {
	return PrefixRedisKey + strconv.Itoa(int(userId))
}

type User struct {
	id             uint64
	name           string
	preferredVideo BitMap
}

func (u *User) ID() uint64 {
	return u.id
}

func (u *User) Name() string {
	return u.name
}

func (u *User) BitMap() BitMap {
	return u.preferredVideo
}

func (u *User) SetBitMap(bitMap BitMap) {
	u.preferredVideo = bitMap
}

type Video struct {
	id           uint64
	name         string
	watchedCount int64
	likedCount   int64
}

func (v *Video) ID() uint64 {
	return v.id
}

func (v *Video) Name() string {
	return v.name
}

func (v *Video) WatchedCount() int64 {
	return atomic.LoadInt64(&v.watchedCount)
}

func (v *Video) LikedCount() int64 {
	return atomic.LoadInt64(&v.likedCount)
}

func (v *Video) UpdateWatchedCount() {
	atomic.AddInt64(&v.watchedCount, 1)
}

func (v *Video) UpdatelikedCount(negative bool) bool {
	//try 10 times
	for i := 0; i < 10; i++ {
		oldValue := atomic.LoadInt64(&v.likedCount)
		if oldValue == 0 && negative {
			return true
		} else if oldValue >= 0 && !negative {
			if atomic.CompareAndSwapInt64(&v.likedCount, oldValue, oldValue+1) {
				return true
			}
		} else {
			if atomic.CompareAndSwapInt64(&v.likedCount, oldValue, oldValue-1) {
				return true
			}
		}
		//sleep 10ms
		time.Sleep(10 * time.Millisecond)
	}
	return false
}

func NewUser() *User {
	id := atomic.AddUint64(&userID, 1)
	return &User{
		id:   id,
		name: strconv.Itoa(int(id)),
	}
}

func NewVideo() *Video {
	id := atomic.AddUint64(&videoID, 1)
	return &Video{
		id:   id,
		name: strconv.Itoa(int(id)),
	}
}
