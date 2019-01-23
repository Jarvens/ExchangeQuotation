// auth: kunlun
// date: 2019-01-23
// description: 消费MQ消息
package entity

type Tick struct {
	Cmd    string     `json:"cmd"`    //指令  tick depth  depthplus
	Symbol string     `json:"symbol"` //交易对
	Ts     int64      `json:"ts"`     //成交时间戳
	Data   []TickItem `json:"data"`   //成交数据
}

type TickItem struct {
	Amount float32 `json:"amount"` //成交量
	Price  float32 `json:"price"`  //成交价
	Dir    int     `json:"dir"`    //红绿
}
