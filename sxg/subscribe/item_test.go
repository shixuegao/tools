package subscribe_test

import (
	"fmt"
	"math/rand"
	"sxg/subscribe"
	"sync"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
)

type Videos struct {
	count  int
	videos []*subscribe.Video
}

type Users struct {
	count int
	users []*subscribe.User
}

func Test0(t *testing.T) {
	videos := Videos{}
	users := Users{}
	redisPool := redisPool("127.0.0.1:6379")
	defer redisPool.Close()
	//init
	start := time.Now()
	videos.count = 200
	videos.videos = make([]*subscribe.Video, videos.count)
	for i := 0; i < videos.count; i++ {
		videos.videos[i] = subscribe.NewVideo()
	}
	users.count = 1000000
	users.users = make([]*subscribe.User, users.count)
	for i := 0; i < users.count; i++ {
		users.users[i] = subscribe.NewUser()
		key := subscribe.PreferredRedisKey(users.users[i].ID())
		bitMap, _ := subscribe.NewRedisBitMapWithFunc(key, 16777216 /*2^24*/, 0, func() redis.Conn {
			return redisPool.Get()
		})
		users.users[i].SetBitMap(bitMap)
	}
	fmt.Printf("初始化共花费: %vms\n", time.Since(start)/time.Millisecond)
	wait := new(sync.WaitGroup)
	exit := make(chan byte)
	interval(wait, &users, &videos, exit)
	go func() {
		//run 300s
		time.Sleep(300 * time.Second)
		close(exit)
	}()
	wait.Wait()
}

func redisPool(addr string) *redis.Pool {
	rp := subscribe.DefaultRedisParams(addr)
	return &redis.Pool{
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", rp.Addr, redis.DialConnectTimeout(rp.ConnectTimeout),
				redis.DialReadTimeout(rp.ReadTimeout), redis.DialWriteTimeout(rp.WriteTimeout))
			if err != nil {
				return nil, err
			}
			if rp.Password != "" {
				_, err = conn.Do(subscribe.RO_AUTH, rp.Password)
				if err != nil {
					conn.Close()
					return nil, err
				}
			}
			_, err = conn.Do(subscribe.RO_SELECT, rp.Db)
			if err != nil {
				conn.Close()
				return nil, err
			}
			return conn, nil
		},
	}
}

func interval(wait *sync.WaitGroup, users *Users, videos *Videos, exit <-chan byte) {
	f := func() {
		userIndex := rand.Intn(users.count)
		videoIndex := rand.Intn(videos.count)
		user := users.users[userIndex]
		video := videos.videos[videoIndex]
		//watch
		video.UpdateWatchedCount()
		//like
		doLiked := rand.Intn(2) == 0
		liked := user.BitMap().Existed(video.ID())
		if !doLiked && liked {
			if video.UpdatelikedCount(true) {
				user.BitMap().Remove(video.ID())
			} else {
				fmt.Println("Cancel thumb failed...")
			}
		} else if doLiked && !liked {
			if video.UpdatelikedCount(false) {
				user.BitMap().Put(video.ID())
			} else {
				fmt.Println("Thumb failed...")
			}
		}
	}
	//10 goroutine
	for i := 0; i < 10; i++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			for {
				f()

				select {
				case <-exit:
					return
				default:
					time.Sleep(time.Millisecond)
				}
			}
		}()
	}
	//print
	wait.Add(1)
	go func() {
		defer wait.Done()
		for {
			fmt.Println("--------------------------")
			for i := 0; i < videos.count; i++ {
				video := videos.videos[i]
				watched := video.WatchedCount()
				liked := video.LikedCount()
				if watched+liked > 0 {
					fmt.Printf("视频[%v], 观看数[%v], 喜欢数[%v]\n", video.Name(), watched, liked)
				}
			}

			select {
			case <-exit:
				return
			default:
				time.Sleep(10 * time.Second)
			}
		}
	}()
}
