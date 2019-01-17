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

var GlobalOption = make(map[string]SubOption)

// 订阅参数选项
type SubOption struct {
	Conn         *Connection //连接
	LastPushTime time.Time   //最后推送时间
	Rate         int         //推送速率  默认5秒一次
	Symbol       []string    //交易对
	Type         string      //推送类型
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
		result, err := MessageHandle(conn, data)
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
	//客户端下线  删除链接
	addr := conn.wsConnect.RemoteAddr().String()
	delete(GlobalOption, addr)
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

func MessageHandle(conn *Connection, data []byte) (result *request.SubResult, err error) {
	var option SubOption
	result = &request.SubResult{}
	var sub request.SubRequest
	err = json.Unmarshal(data, &sub)
	if err != nil {
		return result.SubFailure("100010", "订阅错误"), err
	}

	if sub.OpK == "sub" {
		var str = make([]string, 0)
		str = strings.Split(sub.Opv, ".")
		option.Conn = conn
		option.Rate = sub.Rate
		option.Symbol = append(option.Symbol, str[1])
		option.LastPushTime = time.Now()
		Sub(option)
		return result.SubSuccess(), nil
	} else if sub.Opv == "unsub" {

	}
	return
}

func Task() {
	task := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-task.C:
			for key, value := range GlobalOption {
				rate := value.Rate
				if int(time.Now().Unix()-value.LastPushTime.Unix()) >= rate {
					fmt.Printf("时间差值已达到，开始推送: %s\n", value.Conn.wsConnect.RemoteAddr().String())
					result := request.NewTick()
					data, err := json.Marshal(result)
					if err != nil {
						panic(err)
					}
					value.Conn.WriteMessage(data)
					value.LastPushTime = time.Now()
					GlobalOption[key] = value
				}
			}
		}
	}
}

// 订阅
func Sub(option SubOption) {
	address := option.Conn.wsConnect.RemoteAddr().String()
	// 判断是否已经订阅
	if _, ok := GlobalOption[address]; !ok {
		GlobalOption[address] = option
	} else {
		fmt.Println("重复订阅")
	}
	for _, v := range GlobalOption {
		fmt.Println(v.Symbol)
	}
}

// 取消订阅
//func Unsub(conn *websocket.Conn,subRequest request.SubRequest)request.SubResult  {
//	opv:=subRequest.OpK
//
//	var opvSlice=make([]string,0)
//	opvSlice=strings.Split(opv,".")
//
//}
