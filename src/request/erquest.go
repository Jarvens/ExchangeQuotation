// auth: kunlun
// date: 2019-01-16
// description:
package request

// 订阅返回
type SubResult struct {
	Id         string `json:"id"`          //客户端编号，暂时作为保留属性
	Status     string `json:"status"`      //请求状态
	ErrCode    string `json:"err_code"`    //错误编码
	ErrMessage string `json:"err_message"` //错误信息
	Timestamp  uint32 `json:"timestamp"`   //时间戳
}

// 订阅请求
type SubRequest struct {
	Id  string `json:"id"`  //客户端标识
	Sub string `json:"sub"` //请求
}

// 成交数据
type Tick struct {
	Amount float32 `json:"amount"` //成交价
	Open   float32 `json:"open"`   //开盘价
	Close  float32 `json:"close"`  //收盘价
	High   float32 `json:"high"`   //最高价
	Ts     int64   `json:"ts"`     //时间戳
	Id     int64   `json:"id"`     //id
	Count  int32   `json:"count"`  //成交量
	Low    float32 `json:"low"`    //最低价
	Vol    float32 `json:"vol"`    //日成交量
}
