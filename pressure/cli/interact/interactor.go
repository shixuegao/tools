package interact

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"net"
	"sync"
	"time"

	"umx/tools/pressure/cli/util"
)

const (
	maxPacketSize   = 1400
	maxResponseSize = 100 * 1024
)

type Interactor struct {
	lock *sync.Mutex
	conn *net.UDPConn
}

func NewInteractor(addr string) (*Interactor, error) {
	raddr, err := net.ResolveUDPAddr("udp", addr)
	if nil != err {
		return nil, err
	}
	conn, err := net.DialUDP("udp", nil, raddr)
	if nil != err {
		return nil, err
	}
	return &Interactor{
		lock: new(sync.Mutex),
		conn: conn,
	}, nil
}

func (i *Interactor) Exchange(path string, v interface{}) (*Response, error) {
	i.lock.Lock()
	defer i.lock.Unlock()
	data, err := toWrapperBytes(path, v)
	if nil != err {
		return nil, err
	}
	//write
	split := util.SplitToFixedByteArray(data, maxPacketSize)
	for _, v := range split {
		_, err := i.conn.Write(v)
		if nil != err {
			return nil, err
		}
	}
	//read
	size := 0
	buffer := make([]byte, maxPacketSize)
	for {
		err := i.conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		if nil != err {
			return nil, err
		}
		n, err := i.conn.Read(buffer[size:])
		if nil != err {
			break
		}
		size += n
		if size >= len(buffer) {
			if total := size + maxPacketSize; total > maxResponseSize {
				return nil, errors.New("overlong response length")
			}
			appended := make([]byte, maxPacketSize)
			buffer = append(buffer, appended...)
		}
	}
	return toResponse(buffer[:size])
}

func (i *Interactor) Close() error {
	return i.conn.Close()
}

func toWrapperBytes(path string, v interface{}) ([]byte, error) {
	c, err := json.Marshal(v)
	if nil != err {
		return nil, err
	}
	pLen := len(path)
	cLen := len(c)
	wrapper := Wrapper{
		Version: 0,
		PLength: uint16(pLen),
		Path:    []byte(path),
		CLength: uint32(cLen),
		Content: []byte(c),
	}
	size := 1 + 2 + pLen + 4 + cLen
	buffer := bytes.NewBuffer(make([]byte, 0, size))
	//version
	binary.Write(buffer, binary.BigEndian, wrapper.Version)
	//path length
	binary.Write(buffer, binary.BigEndian, wrapper.PLength)
	//path
	binary.Write(buffer, binary.BigEndian, wrapper.Path)
	//content length
	binary.Write(buffer, binary.BigEndian, wrapper.CLength)
	//content
	binary.Write(buffer, binary.BigEndian, wrapper.Content)
	bytes := make([]byte, size)
	_, err = buffer.Read(bytes)
	if nil != err {
		return nil, err
	}
	return bytes, nil
}

func toResponse(data []byte) (*Response, error) {
	resp := Response{}
	err := json.Unmarshal(data, &resp)
	if nil != err {
		return nil, err
	}
	return &resp, nil
}
