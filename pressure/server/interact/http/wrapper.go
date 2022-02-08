package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"umx/tools/pressure/server/util"
)

const (
	maxBodySize = 200 * 1024
)

var pEmptyWrapper = new(InnerWrapper)
var wrapperType = reflect.TypeOf(pEmptyWrapper)

type InnerWrapper struct {
	pattern string
	w       http.ResponseWriter
	r       *http.Request
}

func (iw *InnerWrapper) Method() string {
	return iw.r.Method
}

func (iw *InnerWrapper) IsLegalMethod(method string) bool {
	return method == iw.r.Method
}

func (iw *InnerWrapper) FormParser() error {
	return iw.r.ParseForm()
}

func (iw *InnerWrapper) Pattern() string {
	return iw.pattern
}

func (iw *InnerWrapper) RawResponse() http.ResponseWriter {
	return iw.w
}

func (iw *InnerWrapper) RawRequest() *http.Request {
	return iw.r
}

func (iw *InnerWrapper) Header(key string) string {
	headers := iw.r.Header[key]
	if len(headers) <= 0 {
		return ""
	}
	return headers[0]
}

func (iw *InnerWrapper) ParamString(key string) string {
	values := iw.r.Form[key]
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func (iw *InnerWrapper) ParamMultiString(key string) []string {
	values := iw.r.Form[key]
	if values == nil {
		return []string{}
	}
	return values
}

func (iw *InnerWrapper) ParamInt(key string) int {
	values := iw.r.Form[key]
	if len(values) == 0 {
		return 0
	}
	i, err := strconv.Atoi(values[0])
	if err != nil {
		return 0
	}
	return i
}

func (iw *InnerWrapper) ParamMultiInt(key string) []int {
	values := iw.r.Form[key]
	if values == nil {
		return []int{}
	}
	iArray := make([]int, len(values))
	for k := 0; k < len(values); k++ {
		i, err := strconv.Atoi(values[k])
		if err != nil {
			return []int{}
		}
		iArray[k] = i
	}
	return iArray
}

//处理Body
func (iw *InnerWrapper) Body() ([]byte, error) {
	length := iw.r.ContentLength
	if length > maxBodySize {
		return nil, errors.New("数据大小超过上限")
	}
	data := make([]byte, length)
	n, err := iw.r.Body.Read(data)
	//Body.Read使用LimitedReader来封装读取。并在读取指定长度的数据后,
	//通过判断是否读取完成来使用EOF
	if nil != err && err != io.EOF {
		return nil, err
	}
	return data[:n], nil
}

//处理File
func (iw *InnerWrapper) FileContent(key string) ([]byte, error) {
	file, _, err := iw.r.FormFile(key)
	if nil != err {
		return nil, err
	}
	data := make([]byte, maxBodySize)
	n, err := file.Read(data)
	if nil != err && err != io.EOF {
		return nil, err
	}
	return data[:n], nil
}

func (iw *InnerWrapper) Write(data []byte) error {
	if _, err := iw.w.Write(data); nil != err {
		return err
	}
	return nil
}

func (iw *InnerWrapper) WriteToJson(v interface{}) error {
	if util.AssertNil(v) {
		return iw.Write([]byte{})
	}
	if data, err := json.Marshal(v); nil != err {
		return err
	} else {
		iw.RawResponse().Header().Set("Content-Type", "application/json")
		return iw.Write(data)
	}
}

func (iw *InnerWrapper) WriteError(code int) {
	iw.w.WriteHeader(code)
}

func (iw *InnerWrapper) Redirection(pattern string) {
	iw.RawResponse().Header().Add("Location", pattern)
	iw.RawResponse().WriteHeader(302)
}
