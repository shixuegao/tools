package subscribe_test

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"math/rand"
	"sxg/subscribe"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func TestBitMap(t *testing.T) {
	offset := 0 * subscribe.LimitMaxValue
	max := subscribe.LimitMaxValue
	// max := uint64(8)
	var bitMap subscribe.BitMap

	//memory
	// bitMap, err := subscribe.NewMemoryBitMap(max, offset)

	//redis
	bitMap, err := subscribe.NewRedisBitMap("bm1", max, offset, subscribe.DefaultRedisParams("127.0.0.1:6379"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer func() {
		if v, ok := bitMap.(*subscribe.RedisBitMap); ok {
			v.Close()
		}
	}()
	err = bitMap.Put(0)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = bitMap.Put(subscribe.LimitMaxValue + 1)
	if err != nil {
		fmt.Println(err.Error())
	}
	//put, existed
	printErr(bitMap.Put(1))
	printErr(bitMap.Put(max))
	fmt.Println("1: ", bitMap.Existed(1))
	fmt.Println("Max: ", bitMap.Existed(max))
	sta, end := bitMap.Range()
	fmt.Printf("Range: %v-%v\n", sta, end)
	err = bitMap.Put(max + 1)
	if err != nil {
		fmt.Println(err.Error())
	}
	//resize
	bitMap.Resize(max * 2)
	bitMap.Put(max * 2)
	fmt.Println("After resize(increase) 1: ", bitMap.Existed(1))
	fmt.Println("After resize(increase) Max: ", bitMap.Existed(max))
	fmt.Println("After resize(increase) 2*Max: ", bitMap.Existed(2*max))
	//resize
	bitMap.Resize(max / 2)
	bitMap.Put(max / 2)
	fmt.Println("After resize(decrease) 1: ", bitMap.Existed(1))
	fmt.Println("After resize(decrease) Max: ", bitMap.Existed(max))
	fmt.Println("After resize(decrease) Max/2: ", bitMap.Existed(max/2))
	//random put, existed, remove
	fmt.Println("-----------------------random-----------------------")
	bitMap.Resize(max)
	for i := 0; i < 100; i++ {
		value := uint64(rand.Int63n(int64(max)) + 1)
		bitMap.Put(value)
		fmt.Printf("(Put)%v: %v\n", value, bitMap.Existed(value))
		bitMap.Remove(value)
		fmt.Printf("(Remove)%v: %v\n", value, bitMap.Existed(value))
	}
	now := time.Now()
	fmt.Println("Count: ", bitMap.Count())
	fmt.Println("Elapsed: ", time.Since(now))
}

func printErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestMysqlBitMap(t *testing.T) {
	bitMap, err := createMysqlBitMap(0 * subscribe.MysqlLimitMaxValue)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer bitMap.RawDb().Close()
	//existed
	// fmt.Println(bitMap.Existed(0))
	// fmt.Println(bitMap.Existed(subscribe.MysqlLimitMaxValue * 2))
	// fmt.Println(bitMap.Existed(subscribe.MysqlLimitMaxValue / 2))
	// fmt.Println(bitMap.ExistedTransaction(subscribe.MysqlLimitMaxValue/4, true))
	//put
	printErr(bitMap.Put(0))
	printErr(bitMap.Put(subscribe.MysqlLimitMaxValue * 2))
	err = bitMap.Put(1)
	printErr(err)
	err = bitMap.Put(9)
	printErr(err)
	err = bitMap.Put(subscribe.MysqlLimitMaxValue)
	printErr(err)
	fmt.Println("1: ", bitMap.Existed(1))
	fmt.Println("9: ", bitMap.ExistedTransaction(9, true))
	fmt.Println("Max: ", bitMap.ExistedTransaction(subscribe.MysqlLimitMaxValue, true))

	//remove
	bitMap.Remove(1)
	bitMap.Remove(subscribe.MysqlLimitMaxValue)
	fmt.Println("1(removed): ", bitMap.ExistedTransaction(1, true))
	fmt.Println("Max(removed): ", bitMap.ExistedTransaction(subscribe.MysqlLimitMaxValue, true))
}

/*
	table: bitmap_test
	columns: id(unsignint), topic(varchar), offset(bigint), bitmap(longblob), stamp(unsignint), create_time(timestamp)
	attention: the size of blob is double to MysqlLimitMaxValue. If you want to fill the blob by using file, should change [max_packet_size] to size(2*MysqlLimitMaxValue)
*/
func createMysqlBitMap(offset uint64) (*subscribe.MysqlBitMap, error) {
	fSqlBlob := func(value uint64) (string, []interface{}) {
		return "select bitmap, stamp from bitmap_test where id=?", []interface{}{1}
	}
	fSqlPut := func(stamp uint64) (string, []interface{}) {
		return "update bitmap_test set bitmap = ?, stamp=? where id=? and stamp=?", []interface{}{stamp + 1, 1, stamp}
	}
	fSqlRemove := func(value uint64, stamp uint64) (string, []interface{}) {
		return "update bitmap_test set bitmap = ?, stamp=? where id=? and stamp=?", []interface{}{stamp + 1, 1, stamp}
	}
	fBlobFromRows := func(rows *sql.Rows) ([]byte, uint64, error) {
		var stamp int
		var blob []byte
		if !rows.Next() {
			return []byte{}, 0, nil
		}

		err := rows.Scan(&blob, &stamp)
		if err != nil {
			return nil, 0, err
		}
		return blob, uint64(stamp), nil
	}
	fConverted := func(p []byte, direction int) []byte {
		if direction == subscribe.DirectionBlobToIndexes {
			//hex -> byte
			r, _ := hex.DecodeString(string(p))
			return r
		} else {
			//byte -> hex
			return []byte(hex.EncodeToString(p))
		}
	}
	sqls := subscribe.Sqls{
		FuncSqlBlob:   fSqlBlob,
		FuncSqlPut:    fSqlPut,
		FuncSqlRemove: fSqlRemove,
	}
	//db
	db, err := sql.Open("mysql", "root:12qwaszx@/test")
	if err != nil {
		return nil, err
	}
	name := fmt.Sprintf("测试用BitMap[%v-%v]", offset+1, offset+subscribe.MysqlLimitMaxValue)
	bitMap, err := subscribe.NewMysqlBitMap(name, offset, db, &sqls, fBlobFromRows, fConverted)
	if err != nil {
		return nil, err
	}
	return bitMap, nil
}
