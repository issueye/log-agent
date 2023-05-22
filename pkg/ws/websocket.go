package ws

import (
	"errors"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

var SMap = new(sync.Map)

type WsConn struct {
	id        string
	inChan    chan []byte
	outChan   chan []byte
	closeChan chan []byte
	isClose   bool
	mutex     sync.Mutex
	conn      *websocket.Conn
}

// 读
func (conn *WsConn) readLoop() {
	for {
		// 确定数据结构
		var (
			data []byte
			err  error
		)

		// 判断客户端是否强制断开了连接
		// 接收数据
		if _, data, err = conn.conn.ReadMessage(); err != nil {
			goto ERR
		}

		// 写入数据
		if err = conn.InChanWrite(data); err != nil {
			goto ERR
		}
	}
ERR:
	conn.CloseConn()
}

// InChanRead 读取inChan的数据
func (conn *WsConn) InChanRead() (data []byte, err error) {
	select {
	case data = <-conn.inChan:
	case <-conn.closeChan:
		err = errors.New("connection is closed")
	}
	return
}

// InChanWrite inChan写入数据
func (conn *WsConn) InChanWrite(data []byte) (err error) {
	select {
	case conn.inChan <- data:
	case <-conn.closeChan:
		err = errors.New("connection is closed")
	}
	return
}

func (conn *WsConn) OutChanRead() (data []byte, err error) {
	select {
	case data = <-conn.outChan:
	case <-conn.closeChan:
		err = errors.New("connection is closed")
	}
	return
}

// OutChanWrite inChan写入数据
func (conn *WsConn) OutChanWrite(data []byte) (err error) {
	select {
	case conn.outChan <- data:
	case <-conn.closeChan:
		err = errors.New("connection is closed")
	}
	return
}

// 写
func (conn *WsConn) writeLoop() {
	for {
		var (
			data []byte
			err  error
		)
		// 读取数据
		if data, err = conn.OutChanRead(); err != nil {
			goto ERR
		}
		// 发送数据
		if err = conn.conn.WriteMessage(1, data); err != nil {
			goto ERR
		}
	}
ERR:
	conn.CloseConn()
}

type WsConnContext struct {
	InData []byte
}

type ctxCallbackFunc func(*WsConnContext) ([]byte, error)

// Listen
func (conn *WsConn) Listen(callback ctxCallbackFunc) {
	for {

		data, err := conn.InChanRead()
		if err != nil {
			goto ERR
		}

		b, err := callback(&WsConnContext{
			InData: data,
		})
		if err != nil {
			fmt.Println("err ", err.Error())
			continue
		}

		err = conn.OutChanWrite(b)
		if err != nil {
			goto ERR
		}
	}

ERR:
	conn.CloseConn()
}

// CloseConn 关闭WebSocket连接
func (conn *WsConn) CloseConn() {
	// 关闭closeChan以控制inChan/outChan策略,仅此一次
	conn.mutex.Lock()
	if !conn.isClose {
		close(conn.closeChan)
		conn.isClose = true
	}

	SMap.Delete(conn.id)

	conn.mutex.Unlock()
	//关闭WebSocket的连接,conn.Close()是并发安全可以多次关闭
	_ = conn.conn.Close()
}

func NewConn(id string, conn *websocket.Conn) *WsConn {
	ws := &WsConn{
		id:        id,
		inChan:    make(chan []byte, 1024),
		outChan:   make(chan []byte, 1024),
		closeChan: make(chan []byte, 1024),
		conn:      conn,
	}

	go ws.readLoop()
	go ws.writeLoop()

	// 放入map中，外部在map中获取连接对象
	SMap.Store(id, ws)

	return ws
}
