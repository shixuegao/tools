package http

import (
	"errors"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"umx/tools/pressure/server/log"
	"umx/tools/pressure/server/util"
)

const (
	root   = "/"
	static = "/static"
)

var logger = log.Logger

//http handler
type Handler interface{}

type HandlerDispatcher struct {
	rootHandler   *innerHandler
	staticHandler *innerHandler
	iHandlers     map[string]*innerHandler
}

func (hw *HandlerDispatcher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	//root
	if path == root && hw.rootHandler != nil {
		hw.rootHandler.handleHttpRequest(w, r)
		return
	}
	//static
	if strings.HasPrefix(path, static) && hw.staticHandler != nil {
		hw.staticHandler.handleHttpRequest(w, r)
		return
	}
	//cors
	hw.setCors(w)
	//handlers
	ih := hw.iHandlers[path]
	if ih != nil {
		ih.handleHttpRequest(w, r)
	}
}

func (hw *HandlerDispatcher) setCors(w http.ResponseWriter) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "POST, PUT, OPTIONS, GET, OPTIONS, DELETE")
	w.Header().Add("Access-Control-Max-Age", "3600")
	w.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
}

func (iw *HandlerDispatcher) Patterns() []string {
	count := 0
	patterns := make([]string, len(iw.iHandlers), len(iw.iHandlers)+2)
	for k := range iw.iHandlers {
		patterns[count] = k
		count++
	}
	if iw.rootHandler != nil {
		patterns = append(patterns, root)
	}
	if iw.staticHandler != nil {
		patterns = append(patterns, static)
	}
	return patterns
}

//need to implement
type funcInnerHandler func(iw *InnerWrapper)

type innerHandler struct {
	pattern string
	method  string
	f       funcInnerHandler
}

func (ih *innerHandler) handleHttpRequest(w http.ResponseWriter, r *http.Request) {
	iw := &InnerWrapper{
		pattern: ih.pattern,
		w:       w,
		r:       r,
	}
	ih.f(iw)
}

func toServerHandlers(handlers map[string]Handler) (*HandlerDispatcher, error) {
	var rootHandler *innerHandler
	var staticHandler *innerHandler
	handlerMap := make(map[string]*innerHandler)
	if len(handlers) > 0 {
		for p1, h1 := range handlers {
			if ok := isStructPtr(h1); !ok {
				return nil, errors.New("Handler不是Struct指针")
			}
			hMap := sortingFuncInnerHandler(h1)
			if len(hMap) > 0 {
				for p2, h2 := range hMap {
					path := combinePath(p1, p2)
					if path == root {
						rootHandler = h2
					} else if path == static {
						staticHandler = h2
					} else {
						handlerMap[path] = h2
					}
				}
			}
		}
	}
	if len(handlerMap) <= 0 {
		return nil, errors.New("无可用的处理器")
	}
	return &HandlerDispatcher{
		rootHandler:   rootHandler,
		staticHandler: staticHandler,
		iHandlers:     handlerMap}, nil
}

func combinePath(path1, path2 string) string {
	b1 := []byte(path1)
	b2 := []byte(path2)
	if b1[len(b1)-1] != '/' {
		b1 = append(b1, '/')
	}
	if b2[0] == '/' {
		b2 = b2[1:]
	}
	return string(b1) + string(b2)
}

//判断是否为struct指针
func isStructPtr(h Handler) bool {
	switch h.(type) {
	case nil:
		return false
	default:
	}
	t := reflect.TypeOf(h)
	if t.Kind() == reflect.Ptr {
		vP := reflect.ValueOf(h)
		v := vP.Elem()
		if v.Kind() == reflect.Struct {
			return true
		}
	}
	return false
}

//结构体中函数转handler
//1.函数必须以大写开头
//2.函数必须驼峰书写, 举例: GetTaskShow.第一个单词表示方法, 后面表示Url
func sortingFuncInnerHandler(handler Handler) map[string]*innerHandler {
	handlerMap := make(map[string]*innerHandler)
	names := make(map[int]string)
	vType := reflect.TypeOf(handler)
	for i := 0; i < vType.NumMethod(); i++ {
		m := vType.Method(i)
		mType := vType.Method(i).Type
		iCount := mType.NumIn()
		oCount := mType.NumOut()
		//iCount = self + params
		if iCount != 2 && oCount != 0 {
			continue
		}
		i0Type := mType.In(1)
		if i0Type != wrapperType {
			continue
		}
		names[i] = m.Name
	}
	vValue := reflect.ValueOf(handler)
	for i := 0; i < vValue.NumMethod(); i++ {
		if name, ok := names[i]; ok {
			m, p := splitName(name)
			if m == "" || p == "" {
				continue
			}
			vMethod := vValue.Method(i)
			handler := innerHandler{
				method:  m,
				pattern: p,
				f: funcInnerHandler(func(iw *InnerWrapper) {
					defer func() {
						if ok, err := util.AssertErr(recover()); ok {
							logger.Warnf("HTTP发生异常, Pattern: %s, Error: %s", p, err.Error())
							iw.WriteError(500)
						}
					}()
					if !iw.IsLegalMethod(m) {
						logger.Warnf("非法的方法请求, Pattern: %s, Needed Method: %s, Real Method: %s", p, m, iw.Method())
						iw.WriteError(405)
						return
					}
					v := reflect.ValueOf(iw)
					vMethod.Call([]reflect.Value{v})
				}),
			}
			handlerMap[p] = &handler
		}
	}
	return handlerMap
}

var nameCompile, _ = regexp.Compile("[A-Z][a-z0-9]+")

//分割名称，组装Url
//Get -> GET, /
//GetTaskShow -> GET, /task/show
//PostTaskAdd -> POST, /task/add
//PutTaskOrder -> PUT, /task/order
func splitName(name string) (method string, pattern string) {
	split := nameCompile.FindAllString(name, -1)
	length := len(split)
	if length == 0 {
		return
	}
	if length >= 1 {
		m := split[0]
		if m != "Get" && m != "Post" && m != "Put" && m != "Delete" {
			return
		}
		method = strings.ToUpper(m)
		if length == 1 {
			pattern = "/"
		} else {
			lowerCases := make([]string, length-1)
			for i := 1; i < length; i++ {
				lowerCases[i-1] = strings.ToLower(split[i])
			}
			pattern = "/" + strings.Join(lowerCases, "/")
		}
	}
	return
}
