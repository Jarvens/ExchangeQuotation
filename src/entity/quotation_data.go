// auth: kunlun
// date: 2019-01-23
// description: 消费MQ消息
package entity

import "time"

type Common struct {
	Cmd    string      `json:"cmd"`    //指令  tick depth  depth_plus
	Symbol string      `json:"symbol"` //交易对
	Ts     int64       `json:"ts"`     //成交时间戳
	Data   interface{} `json:"data"`   //成交数据
}

type DepthItem struct {
	Asks [][]float32
	Bids [][]float32
}

type TickItem struct {
	Amount float32 `json:"amount"` //成交量
	Price  float32 `json:"price"`  //成交价
	Dir    int     `json:"dir"`    //红绿
}

// 模拟Tick数据
func (comm *Common) MockTick() *Common {
	tick := &TickItem{Amount: 0.12, Price: 0.11, Dir: 1}
	return &Common{Cmd: "tick", Symbol: "BTCUSDT", Ts: time.Now().Unix(), Data: tick}

}

// 模拟Depth数据
func (comm *Common) MockDepth() *Common {
	depth := &DepthItem{Asks: [][]float32{{0.12, 1.1}, {0.13, 0.9}}, Bids: [][]float32{{0.9, 0.2}, {0.8, 1.2}}}
	return &Common{Cmd: "depth", Symbol: "BTCUSDT", Ts: time.Now().Unix(), Data: depth}
}

// 模拟DepthPlus数据
func (comm *Common) MockDepthPlus() *Common {
	depthPlus := &DepthItem{Asks: [][]float32{{0.12, 1.1}, {0.13, 0.9}}, Bids: [][]float32{{0.9, 0.2}, {0.8, 1.2}}}
	return &Common{Cmd: "depthplus", Symbol: "BTCUSDT", Ts: time.Now().Unix(), Data: depthPlus}
}
