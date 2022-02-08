package coap

import (
	"net"
	"strconv"
	"sync"
)

const (
	//MethodGet Get
	MethodGet = iota
	//MethodPost Post
	MethodPost
	//MethodPut Put
	MethodPut
	//MethodDelete Delete
	MethodDelete
)

const (
	//maxPktLen udp最大缓存
	maxPktLen = 4096
)

//Connector Coap通信
type Connector struct {
	messageID      uint16
	messageIDMutex *sync.Mutex
	localIP        string
	localPort      int
	remoteIP       string
	remotePort     int
	connector      *net.UDPConn
	buf            []byte
	devType        int
	subDevType     string
	registered     bool
	statistic      *coapStatistic
}

//NewConnector 新建Coap通信
func NewConnector(devType int, localIP string, localPort int, remoteIP string, remotePort int, subDevType string) (*Connector, error) {
	addr := localIP + ":" + strconv.Itoa(localPort)
	lUdpAddr, err := net.ResolveUDPAddr("udp", addr)
	if nil != err {
		return nil, err
	}
	addr = remoteIP + ":" + strconv.Itoa(remotePort)
	rUdpAddr, err := net.ResolveUDPAddr("udp", addr)
	if nil != err {
		return nil, err
	}
	connector, err := net.DialUDP("udp", lUdpAddr, rUdpAddr)
	if nil != err {
		return nil, err
	}
	conn := Connector{
		localIP:        localIP,
		localPort:      localPort,
		remoteIP:       remoteIP,
		remotePort:     remotePort,
		messageID:      0,
		messageIDMutex: new(sync.Mutex),
		buf:            make([]byte, maxPktLen),
		connector:      connector,
		devType:        devType,
		subDevType:     subDevType,
		statistic:      newCoapStatistic(),
	}
	return &conn, nil
}

//本地端口
func (c *Connector) LocalPort() int {
	return c.localPort
}

//getMessageID 获取信息ID
func (c *Connector) getMessageID() uint16 {
	c.messageIDMutex.Lock()
	defer c.messageIDMutex.Unlock()
	if c.messageID == uint16(0xffff) {
		c.messageID = 0
	} else {
		c.messageID = c.messageID + 1
	}
	return c.messageID
}

//Send 发送数据
func (c *Connector) Send(m *Message) error {
	bytes, err := m.MarshalBinary()
	if nil != err {
		return err
	}
	//record
	c.statistic.record(m.Type, m.MessageID)
	//send
	_, err = c.connector.Write(bytes)
	if nil != err {
		return err
	}
	return nil
}

//Receive 接收数据
func (c *Connector) Receive() (Message, error) {
	len, _, err := c.connector.ReadFromUDP(c.buf)
	if nil != err {
		return Message{}, err
	}
	result, err := ParseMessage(c.buf[:len])
	if nil != err {
		return Message{}, err
	}
	//confirm
	c.statistic.confirm(result.Type, result.MessageID)
	return result, nil
}

//Close 关闭
func (c *Connector) Close() error {
	return c.connector.Close()
}
