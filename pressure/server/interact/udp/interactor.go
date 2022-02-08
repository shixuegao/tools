package udp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
	"runtime/debug"
	"umx/tools/pressure/server/util"
)

const (
	maxPacketSize  = 1400
	maxWrapperSize = 200 * 1024
)

type Interactor struct {
	conn     *net.UDPConn
	wrappers map[string]*wrapperBuffer
	closed   chan struct{}
}

type wrapperBuffer struct {
	wrapper *Wrapper
	rest    int
}

func (w *wrapperBuffer) read(data []byte) (ok bool, err error) {
	defer func() {
		r := recover()
		if ok, e := util.AssertErr(r); ok {
			err = e
		}
	}()
	if nil == w.wrapper {
		w.wrapper = &Wrapper{}
		buffer := bytes.NewBuffer(data)
		//Version
		w.wrapper.Version, _ = buffer.ReadByte()
		//PLength
		binary.Read(buffer, binary.BigEndian, &w.wrapper.PLength)
		//Path
		w.wrapper.Path = make([]byte, w.wrapper.PLength)
		buffer.Read(w.wrapper.Path)
		//CLength
		binary.Read(buffer, binary.BigEndian, &w.wrapper.CLength)
		if w.wrapper.CLength > maxWrapperSize {
			return false, errors.New("overlong content length")
		}
		//Content
		w.wrapper.Content = make([]byte, w.wrapper.CLength)
		n, _ := buffer.Read(w.wrapper.Content)
		w.rest = int(w.wrapper.CLength) - n
	} else {
		//Content
		l := len(data)
		if l > w.rest {
			l = w.rest
		}
		n := copy(w.wrapper.Content[int(w.wrapper.CLength)-w.rest:], data[:l])
		w.rest = w.rest - n
	}
	return w.rest <= 0, nil
}

func NewInteractor(addr string) (*Interactor, error) {
	laddr, err := net.ResolveUDPAddr("udp", addr)
	if nil != err {
		return nil, err
	}
	udpConn, err := net.ListenUDP("udp", laddr)
	if nil != err {
		return nil, err
	}
	return &Interactor{
		conn:     udpConn,
		wrappers: make(map[string]*wrapperBuffer),
		closed:   make(chan struct{}),
	}, nil
}

func (i *Interactor) IsRunning() bool {
	select {
	case <-i.closed:
		return true
	default:
		return false
	}
}

func (i *Interactor) Close() error {
	select {
	case <-i.closed:
		return nil
	default:
		close(i.closed)
		return i.conn.Close()
	}
}

func (i *Interactor) StartAsync(fWrapper func(*Wrapper) []byte, fErr func(error, []byte)) {
	go func() {
		buffer := make([]byte, maxPacketSize)
		for {
			n, raddr, err := i.conn.ReadFromUDP(buffer)
			if nil != err {
				fErr(err, debug.Stack())
				return
			}
			w, err := i.handlePacket(n, raddr, buffer)
			if nil != err {
				fErr(err, debug.Stack())
			} else if nil == w {
				continue
			} else {
				data := fWrapper(w)
				result := util.SplitToFixedByteArray(data, maxPacketSize)
				for _, v := range result {
					_, err = i.conn.WriteToUDP(v, raddr)
					if nil != err {
						fErr(err, debug.Stack())
					}
				}
			}
			select {
			case <-i.closed:
				return
			default:
			}
		}
	}()
}

func (i *Interactor) handlePacket(size int, raddr *net.UDPAddr, b []byte) (wrapper *Wrapper, err error) {
	raddrStr := raddr.String()
	wb := i.wrappers[raddrStr]
	if nil == wb {
		wb = &wrapperBuffer{}
		i.wrappers[raddrStr] = wb
	}
	ok, e := wb.read(b[:size])
	if nil != e {
		err = e
		return
	}
	if ok {
		wrapper = wb.wrapper
		wb.wrapper = nil
	}
	return
}
