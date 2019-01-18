// auth: kunlun
// date: 2019-01-16
// description:
package websocket

import (
	"common"
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
	UnAckCounter int         //未响应计数器
}

// 初始化连接
func InitConnection(wsCon *websocket.Conn) (conn *Connection, err error) {
	conn = &Connection{
		wsConnect: wsCon,
		inChan:    make(chan []byte, 1024),
		outChan:   make(chan []byte, 1024),
		closeChan: make(chan byte, 1)}

	// 建立连接
	addr := wsCon.RemoteAddr().String()
	subOption := SubOption{Conn: conn}
	GlobalOption[addr] = subOption
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
	//客户端下线,删除链接
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

// 消息处理器
// 订阅 & 取消订阅 消息解析
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
		Unsub(conn, sub)
	}
	return
}

// 定时任务时间间隔为1秒
// 遍历循环当前连接Map
// 检查每个连接的订阅 交易对  推送频率
// 按需推送
func Task() {
	task := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-task.C:
			for key, value := range GlobalOption {
				// go 不支持三元表达式
				var rate int
				if value.Rate == 0 {
					rate = common.DefaultRate
				} else {
					rate = value.Rate
				}

				symbols := value.Symbol
				if len(symbols) != 0 {
					//TODO 区分订阅为 1min  5min  15min
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
}

// 订阅
// 订阅信息,当前连接全局记录
func Sub(option SubOption) {
	address := option.Conn.wsConnect.RemoteAddr().String()
	// 判断是否已经订阅
	if _, ok := GlobalOption[address]; !ok {
		GlobalOption[address] = option
	} else {
		fmt.Printf("重复订阅,客户端地址: %s\n", option.Conn.wsConnect.RemoteAddr().String())
	}
}

// 取消订阅
// all 全部取消,连接还存在
// [btc,eth] 指定交易对取消 连接还存在
// 返回 SubResult
func Unsub(conn *Connection, subRequest request.SubRequest) (result request.SubResult) {
	opv := subRequest.OpK
	addr := conn.wsConnect.RemoteAddr().String()
	var opvSlice = make([]string, 0)
	opvSlice = strings.Split(opv, ".")
	if opvSlice[1] == "all" {
		//清空 订阅交易对
		option := GlobalOption[addr]
		option.Symbol = append([]string{})
		GlobalOption[addr] = option
	} else {
		//如果不是全部取消。需要遍历即将要取消的交易对数组

	}
	return result
}

// 服务器主动向客户端发送心跳包
// 频率为 5s 一次
// 客户端连续忽略2次,服务端主动断开连接
func Ping() {
	task := time.NewTicker(common.PingInterval * time.Second)
	for {
		select {
		case <-task.C:
			for k, v := range GlobalOption {
				fmt.Println("遍历连接池")
				unAck := v.UnAckCounter
				if unAck >= common.UnAck {
					fmt.Println("客户端未响应，主动断开")
					v.Conn.Close()
					continue
				}
				conn := v.Conn
				ping := common.NewPing()
				data, err := json.Marshal(ping)
				if err != nil {
					panic(err)
				}
				conn.WriteMessage(data)
				//未响应计数器 +1
				v.UnAckCounter = unAck + 1
				GlobalOption[k] = v
			}
		}
	}
}
