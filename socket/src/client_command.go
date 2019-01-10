// auth: kunlun
// date: 2019-01-10
// description:
package main

import (
	"codec"
	"encoding/json"
	"net"
	"time"
	"utils"
)

func main() {

	conn, err := net.Dial("tcp", "0.0.0.0:12345")
	defer conn.Close()
	if err != nil {
		utils.Debug("Connect error: %v", err)
		return
	}

	for {
		data := codec.ResponseData{Dir: "bid",
			Symbol:    "USDT_BTC",
			Ts:        time.Now().UnixNano(),
			Amount:    0.2,
			Price:     0.1,
			DayVolume: 10,
			DayPrice:  0.5,
			DayHigh:   0.5,
			DayLow:    0.2}
		dataBytes, _ := json.Marshal(data)
		dataStr := string(dataBytes)
		proto := codec.NewDefaultProto(dataStr)
		_, err := conn.Write(proto.Encoder())
		if err != nil {
			return
		}
		utils.Debug("Start send message to Server: %v", conn.RemoteAddr())
		time.Sleep(200 * time.Millisecond)

	}

}
