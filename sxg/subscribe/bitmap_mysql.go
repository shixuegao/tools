package subscribe

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"
)

const (
	DirectionBlobToIndexes = iota
	DirectionIndexesToBlob
)

type FuncBlobFromRows func(*sql.Rows) ([]byte, uint64, error)
type FuncConverted func([]byte, int) []byte

var MysqlLimitMaxValue uint64 = LimitMaxValue / 64 //8MB

type Sqls struct {
	//get blob and stamp
	FuncSqlBlob func(value uint64) (string, []interface{})
	//put blob bases on stamp
	FuncSqlPut func(stamp uint64) (string, []interface{})
	//remove blob bases on stamp
	FuncSqlRemove func(value uint64, stamp uint64) (string, []interface{})
}

func (sqls *Sqls) validate() error {
	if sqls.FuncSqlBlob == nil {
		return errors.New("SqlBlob函数不能为null")
	}
	if sqls.FuncSqlPut == nil {
		return errors.New("SqlPut函数不能为null")
	}
	if sqls.FuncSqlRemove == nil {
		return errors.New("SqlRemove函数不能为null")
	}
	return nil
}

/*
	1.使用MediumBlob来存储(默认最大16M)
	2.单个MysqlBitMap上限进行限制
*/
type MysqlBitMap struct {
	limit         uint64
	offset        uint64
	name          string
	sqls          *Sqls
	db            *sql.DB
	lock          *sync.RWMutex
	fBlobFromRows FuncBlobFromRows
	fConverted    FuncConverted
}

func (mbm *MysqlBitMap) doTransaction(f func(tx *sql.Tx) error) (err error) {
	tx, err := mbm.db.Begin()
	defer func() {
		if tx == nil {
			return
		}
		if err == nil {
			err = tx.Commit()
		} else {
			err = tx.Rollback()
		}
	}()
	if err != nil {
		return err
	}
	return f(tx)
}

func (mbm *MysqlBitMap) getBlob(value uint64, tx *sql.Tx) ([]byte, uint64, error) {
	strSql, args := mbm.sqls.FuncSqlBlob(value)
	if strSql == "" {
		return nil, 0, errors.New("SqlExisted函数返回值为空")
	}
	var err error
	var blob []byte
	var rows *sql.Rows
	var stamp uint64
	if tx != nil {
		rows, err = tx.Query(strSql, args...)
		if err == nil {
			defer rows.Close()
			blob, stamp, err = mbm.fBlobFromRows(rows)
		}
	} else {
		rows, err = mbm.db.Query(strSql, args...)
		if err == nil {
			defer rows.Close()
			blob, stamp, err = mbm.fBlobFromRows(rows)
		}
	}
	return blob, stamp, err
}

func (mbm *MysqlBitMap) Existed(value uint64) bool {
	return mbm.ExistedTransaction(value, false)
}

//maybe occur panic when transaction is true
func (mbm *MysqlBitMap) ExistedTransaction(value uint64, transaction bool) bool {
	if v, ok := normalizedWithLimit(value, mbm.offset, MysqlLimitMaxValue); !ok {
		return false
	} else {
		mbm.lock.RLock()
		defer mbm.lock.RUnlock()
		var blob []byte
		if transaction {
			err := mbm.doTransaction(func(tx *sql.Tx) error {
				if blob0, _, err := mbm.getBlob(value, tx); err != nil {
					return err
				} else {
					blob = blob0
					return nil
				}
			})
			if err != nil {
				return false
			}
		} else {
			if blob0, _, err0 := mbm.getBlob(value, nil); err0 != nil {
				return false
			} else {
				blob = blob0
			}
		}
		indexes := mbm.fConverted(blob, DirectionBlobToIndexes)
		index, mask := indexAndMask(v)
		if len(indexes) <= int(index) {
			return false
		}
		return (indexes[index] & mask) > 0
	}
}

func (mbm *MysqlBitMap) Put(value uint64) error {
	return mbm.PutTransaction(value, false)
}

//maybe occur panic when transaction is true
func (mbm *MysqlBitMap) PutTransaction(value uint64, transaction bool) error {
	if v, ok := normalizedWithLimit(value, mbm.offset, MysqlLimitMaxValue); !ok {
		return fmt.Errorf("非法的值[%d], 值必须在[%d, %d]之间", value, mbm.offset+1, mbm.offset+MysqlLimitMaxValue)
	} else {
		mbm.lock.Lock()
		defer mbm.lock.Unlock()
		if transaction {
			err := mbm.doTransaction(func(tx *sql.Tx) error {
				blob, stamp, err0 := mbm.getBlob(value, tx)
				if err0 != nil {
					return err0
				}
				indexes := mbm.fConverted(blob, DirectionBlobToIndexes)
				index, mask := indexAndMask(v)
				if len(indexes) <= int(index) {
					return fmt.Errorf("数据库的Blob[%v]可能被改变, 长度[%v]无法存储值[%v]", mbm.name, len(blob), value)
				}
				indexes[index] = indexes[index] | mask
				strSql, args := mbm.sqls.FuncSqlPut(stamp)
				if strSql == "" {
					return errors.New("SqlPut函数返回值为空")
				}
				blob = mbm.fConverted(indexes, DirectionIndexesToBlob)
				args0 := append([]interface{}{blob}, args...)
				r, err0 := tx.Exec(strSql, args0...)
				if err0 != nil {
					return err0
				}
				if count, _ := r.RowsAffected(); count != 1 {
					return errors.New("执行不成功, 指定的BitMap没有更新")
				}
				return nil
			})
			if err != nil {
				return err
			}
		} else {
			blob, stamp, err := mbm.getBlob(value, nil)
			if err != nil {
				return err
			}
			indexes := mbm.fConverted(blob, DirectionBlobToIndexes)
			index, mask := indexAndMask(v)
			if len(indexes) <= int(index) {
				return fmt.Errorf("数据库的Blob[%v]可能被改变, 长度[%v]无法存储值[%v]", mbm.name, len(blob), value)
			}
			indexes[index] = indexes[index] | mask
			strSql, args := mbm.sqls.FuncSqlPut(stamp)
			if strSql == "" {
				return errors.New("SqlPut函数返回值为空")
			}
			blob = mbm.fConverted(indexes, DirectionIndexesToBlob)
			args0 := append([]interface{}{blob}, args...)
			r, err := mbm.db.Exec(strSql, args0...)
			if err != nil {
				return err
			}
			if count, _ := r.RowsAffected(); count != 1 {
				return errors.New("执行不成功, 指定的BitMap没有更新")
			}
		}
		return nil
	}
}

func (mbm *MysqlBitMap) Remove(value uint64) {
	mbm.RemoveTransaction(value, false)
}

func (mbm *MysqlBitMap) RemoveTransaction(value uint64, transaction bool) {
	if v, ok := normalizedWithLimit(value, mbm.offset, MysqlLimitMaxValue); ok {
		mbm.lock.Lock()
		defer mbm.lock.Unlock()
		if transaction {
			mbm.doTransaction(func(tx *sql.Tx) error {
				blob, stamp, err := mbm.getBlob(value, tx)
				if err != nil {
					return err
				}
				indexes := mbm.fConverted(blob, DirectionBlobToIndexes)
				index, mask := indexAndMask(v)
				if len(indexes) <= int(index) {
					return fmt.Errorf("数据库的Blob[%v]可能被改变, 长度[%v]无法存储值[%v]", mbm.name, len(blob), value)
				}
				indexes[index] = indexes[index] & (^mask)
				strSql, args := mbm.sqls.FuncSqlRemove(value, stamp)
				if strSql == "" {
					return errors.New("SqlPut函数返回值为空")
				}
				blob = mbm.fConverted(indexes, DirectionIndexesToBlob)
				args0 := append([]interface{}{blob}, args...)
				_, err = tx.Exec(strSql, args0...)
				return err
			})
		} else {
			blob, stamp, err := mbm.getBlob(value, nil)
			if err != nil {
				return
			}
			indexes := mbm.fConverted(blob, DirectionBlobToIndexes)
			index, mask := indexAndMask(v)
			if len(indexes) <= int(index) {
				return
			}
			indexes[index] = indexes[index] & (^mask)
			strSql, args := mbm.sqls.FuncSqlRemove(value, stamp)
			if strSql == "" {
				return
			}
			blob = mbm.fConverted(indexes, DirectionIndexesToBlob)
			args0 := append([]interface{}{blob}, args...)
			mbm.db.Exec(strSql, args0...)
		}
	}
}

//useless
func (mbm *MysqlBitMap) Resize(value uint64) error {
	return nil
}

func (mbm *MysqlBitMap) Range() (min, max uint64) {
	mbm.lock.RLock()
	defer mbm.lock.RUnlock()
	return 1 + mbm.offset, mbm.limit + mbm.offset
}

func (mbm *MysqlBitMap) Size() int {
	mbm.lock.RLock()
	defer mbm.lock.RUnlock()
	index, _ := indexAndMask(mbm.limit)
	return int(index) + 1
}

func (mbm *MysqlBitMap) Reset() {
	mbm.lock.Lock()
	defer mbm.lock.Unlock()
	var blob []byte
	err := mbm.doTransaction(func(tx *sql.Tx) error {
		blob0, _, err := mbm.getBlob(mbm.offset+1, tx)
		if err != nil {
			return err
		}
		blob = blob0
		return nil
	})
	if err != nil || len(blob) == 0 {
		mbm.limit = 1
		return
	}
	indexes := mbm.fConverted(blob, DirectionBlobToIndexes)
	mbm.limit = uint64(len(indexes))
	if mbm.limit > MysqlLimitMaxValue {
		mbm.limit = MysqlLimitMaxValue
	}
}

func (mbm *MysqlBitMap) Count() int {
	mbm.lock.RLock()
	defer mbm.lock.RUnlock()
	var blob []byte
	err := mbm.doTransaction(func(tx *sql.Tx) error {
		blob0, _, err := mbm.getBlob(mbm.offset+1, tx)
		if err != nil {
			return err
		}
		blob = blob0
		return err
	})
	if err != nil || len(blob) == 0 {
		return 0
	}
	return caculateBitCount(blob)
}

func (mbm *MysqlBitMap) Name() string {
	return mbm.name
}

func (mbm *MysqlBitMap) RawDb() *sql.DB {
	return mbm.db
}

func NewMysqlBitMap(name string, offset uint64, db *sql.DB, sqls *Sqls, fBlobFromRows FuncBlobFromRows, fConverted FuncConverted) (*MysqlBitMap, error) {
	if name == "" {
		return nil, errors.New("名称不能为空")
	}
	if offset/MysqlLimitMaxValue != 0 {
		return nil, fmt.Errorf("Offset[%v]不是[%v]的整数倍", offset, MysqlLimitMaxValue)
	}
	if db == nil {
		return nil, errors.New("DB不能为null")
	}
	if sqls == nil {
		return nil, errors.New("Sqls不能为null")
	}
	if err := sqls.validate(); err != nil {
		return nil, err
	}
	if fBlobFromRows == nil {
		return nil, errors.New("FuncBlobFromRows函数不能为null")
	}
	//锁定为16M
	max := offset + MysqlLimitMaxValue
	value, ok := normalizedWithLimit(max, offset, MysqlLimitMaxValue)
	if !ok {
		return nil, fmt.Errorf("最大值[%d]必须在[%d]到[%d]之间", max, offset+1, offset+LimitMaxValue)
	}
	bitMap := MysqlBitMap{
		limit:         value,
		offset:        offset,
		name:          name,
		db:            db,
		sqls:          sqls,
		lock:          new(sync.RWMutex),
		fBlobFromRows: fBlobFromRows,
		fConverted:    fConverted,
	}
	return &bitMap, nil
}
