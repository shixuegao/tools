package util

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func WriteTxtFile(prefix, content string) (err error) {
	filename := prefix + time.Now().Format("20060102150405") + ".txt"
	file, e := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 066)
	if nil != e {
		err = e
		return
	}
	defer func() {
		if file != nil {
			file.Close()
		}
		if ok, e := AssertErr(recover()); ok {
			err = e
		}
	}()
	//write
	_, e = file.WriteString(content)
	if nil != e {
		err = e
	}
	return
}

func AssertErr(t interface{}) (bool, error) {
	switch v := t.(type) {
	case error:
		return true, v
	default:
	}
	return false, nil
}

func GetContentFromTxtFile(filename string) (string, error) {
	if filename == "" {
		return "", errors.New("文件名称不能为空")
	}
	//读取文件内容
	if !strings.HasSuffix(filename, ".txt") {
		return "", errors.New("文件格式不正确, 只接收以.txt结尾的文件")
	}
	content, err := ioutil.ReadFile(filename)
	if nil != err {
		return "", fmt.Errorf("读取文件%s失败-->%s", filename, err.Error())
	}
	return string(content), nil
}

//将原始byte数组分割成指定大小的多个byte数组
func SplitToFixedByteArray(src []byte, size int) [][]byte {
	if len(src) <= size {
		return [][]byte{src}
	} else {
		count := 0
		length := len(src)
		if length%size == 0 {
			count = length / size
		} else {
			count = length/size + 1
		}
		result := make([][]byte, count)
		for i := 0; i < count; i++ {
			sta := i * size
			end := (i + 1) * size
			if end > length {
				end = length
			}
			result[i] = src[sta:end]
		}
		return result
	}
}
