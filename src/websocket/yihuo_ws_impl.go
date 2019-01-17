// auth: kunlun
// date: 2019-01-16
// description:
package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"request"
	"strings"
	"sync"
	"time"
)

type Connection struct {
	wsConnect *websocket.Conn
	inChan    chan []byte
	outChan   chan []byte
	closeChan chan byte
	mutex     sync.Mutex
	isClosed  bool
}

var GlobalOption map[string]SubOption

// 订阅参数选项
type SubOption struct {
	Conn         *websocket.Conn //连接
	LastPushTime time.Time       //最后推送时间
	Rate         int             //推送速率  默认5秒一次
	Symbol       []string        //交易对
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
		result, err := MessageHandle(conn.wsConnect, data)
		if err != nil {
			log.Panic(err)
		}
		v, err1 := json.Marshal(result)
		if err1 != nil {
			log.Panic(err1)
		}
		conn.WriteMessage(v)
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

func MessageHandle(conn *websocket.Conn, data []byte) (result *request.SubResult, err error) {
	var option = make([]SubOption, 1)
	result = &request.SubResult{}
	var sub request.SubRequest
	err = json.Unmarshal(data, &sub)
	if err != nil {
		return result.SubFailure("100010", "订阅错误"), err
	}

	if sub.OpK == "sub" {
		var str = make([]string, 0)
		str = strings.Split(sub.Opv, ".")
		option[0].Conn = conn
		option[0].Rate = sub.Rate
		option[0].Symbol = append(option[0].Symbol, str[1])
		option[0].LastPushTime = time.Now()
		GlobalOption = append(GlobalOption, option...)
		fmt.Println("全局参数: ", GlobalOption)
		return result.SubSuccess(), nil
	} else if sub.Opv == "unsub" {

	}
	return
}

func Task() {
	fmt.Println("初始化定时任务")
	task := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-task.C:
			for _, conn := range GlobalOption {
				fmt.Printf("exec task time: %d conn: %s \n", time.Now().Unix(), conn.Symbol)

			}
		}
	}
}
