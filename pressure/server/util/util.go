package util

import (
	"net/http"
	"regexp"
	"strconv"
	"sync/atomic"
	"time"

	uuid "github.com/satori/go.uuid"
)

const (
	//regexIpv4 IPV4正则表达式
	regexIpv4 = "^(1\\d{2}|2[0-4]\\d|25[0-5]|[1-9]\\d|[1-9])\\.(1\\d{2}|2[0-4]\\d|25[0-5]|[1-9]\\d|\\d)\\.(1\\d{2}|2[0-4]\\d|25[0-5]|[1-9]\\d|\\d)\\.(1\\d{2}|2[0-4]\\d|25[0-5]|[1-9]\\d|\\d)$"
)

//编号
var number int32 = 0

//IsLegalIpv4 检查IPV4是否合法
func IsLegalIpv4(ip string) bool {
	if ip == "localhost" {
		return true
	}
	matched, _ := regexp.Match(regexIpv4, []byte(ip))
	return matched
}

//IsLegalPort 检查端口是否合法
func IsLegalPort(port int) bool {
	return port > 1000 && port <= 65535
}

//NewSignature 获取一个新的签名
func NewSignature(unit string) string {
	value := atomic.AddInt32(&number, 1)
	return unit + "_" + strconv.Itoa(int(value))
}

//CurrentTime 当前时间字符串, 格式为2006-01-02 15:04:05
func CurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

//TightTime
func TightTime() string {
	return time.Now().Format("20060102150405")
}

//NowMillisecond 返回当前毫秒时间戳
func NowMillisecond() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

//UUID
func UUID() string {
	return uuid.NewV4().String()
}

//WriteToHttpResponse 写txt文件到Http中
func WriteToHttpResponse(filename string, contentType string, content []byte, response http.ResponseWriter) error {
	//tell browser dont do cache
	response.Header().Set("Pragma", "No-cache")
	response.Header().Set("Expires", "0")
	response.Header().Set("Cache-Control", "No-cache")
	//txt character
	response.Header().Set("Content-Type", contentType)
	response.Header().Set("Content-Disposition", "attachment; filename*=utf-8''"+filename)
	response.Header().Set("Character-Encoding", "utf-8")
	length := len(content)
	response.Header().Set("Content-Length", strconv.Itoa(length))
	//write to http
	_, err := response.Write(content)
	if nil != err {
		return err
	}
	return nil
}

//错误断言
func AssertErr(t interface{}) (bool, error) {
	switch v := t.(type) {
	case error:
		return true, v
	default:
	}
	return false, nil
}

//nil断言
func AssertNil(t interface{}) bool {
	switch t.(type) {
	case nil:
		return true
	default:
		return false
	}
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

func ResetIntArray(p []int, v int) {
	for i := 0; i < len(p); i++ {
		p[i] = v
	}
}
