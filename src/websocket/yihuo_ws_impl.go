// auth: kunlun
// date: 2019-01-16
// description:
package websocket

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
)

type Connection struct {
	wsConnect *websocket.Conn
	inChan    chan []byte
	outChan   chan []byte
	closeChan chan byte
	mutex     sync.Mutex
	isClosed  bool
}

// 初始化连接
func InitConnection(wsCon *websocket.Conn) (conn *Connection, err error) {
	conn = &Connection{
		wsConnect: wsCon,
		inChan:    make(chan []byte, 1024),
		outChan:   make(chan []byte, 1024),
		closeChan: make(chan byte, 1)}
	go conn.readLoop()
	go conn.writeLoop()
	return
}

// 读取消息
func (conn *Connection) ReadMessage() (data []byte, err error) {
	select {
	case data = <-conn.inChan:
		conn.WriteMessage([]byte("你好"))
	case <-conn.closeChan:
		err = errors.New("connection is closed")
		fmt.Println("Connection is closed")
	}
	return
}

// 写消息
func (conn *Connection) WriteMessage(data []byte) (err error) {
	select {
	case conn.outChan <- data:
	case <-conn.closeChan:
		err = errors.New("connection is closed")
		fmt.Println("Connection is closed")
	}
	return
}

// 关闭连接
func (conn *Connection) Close() {
	conn.wsConnect.Close()
	conn.mutex.Lock()
	if !conn.isClosed {
		close(conn.closeChan)
		conn.isClosed = true
	}
	conn.mutex.Unlock()
}

// 循环读取
func (conn *Connection) readLoop() {
	var (
		data []byte
		err  error
	)

	for {
		if _, data, err = conn.wsConnect.ReadMessage(); err != nil {
			goto ERR
		}

		select {
		case conn.inChan <- data:
		case <-conn.closeChan:
			goto ERR
		}
	}

ERR:
	conn.Close()
}

// 循环写入
func (conn *Connection) writeLoop() {
	var (
		data []byte
		err  error
	)

	for {
		select {
		case data = <-conn.outChan:
		case <-conn.closeChan:
			goto ERR
		}
		if err = conn.wsConnect.WriteMessage(websocket.TextMessage, data); err != nil {
			goto ERR
		}
	}
ERR:
	conn.Close()
}
