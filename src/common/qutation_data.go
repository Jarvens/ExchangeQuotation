// auth: kunlun
// date: 2019-01-22
// description:
package common

import (
	"encoding/json"
)

type Quotation struct {
	Cmd    string
	Symbol string
	Data   Item
}

type Item struct {
	Asks [][]float32
	Bids [][]float32
}

// 存储结构
type PersistenceData struct {
	Id          int64  `json:"id"`
	Symbol      string `json:"symbol"`
	SymbolAlias string `json:"symbol_alias"`
	Dir         int    `json:"dir"`
	Ts          int64  `json:"ts"` //存储时间为纳秒
}

func MockData() {

	var data = Quotation{Cmd: "cmd", Symbol: "btcusdt", Data: struct {
		Asks [][]float32
		Bids [][]float32
	}{Asks: [][]float32{{1.1, 1.1}, {1.1, 1.2}}, Bids: [][]float32{{2.1, 2.1}, {2.2, 2.2}}}}
	bytes, _ := json.Marshal(data)
	PublishMessage("go.direct.exchange", "bind1", bytes)

}
