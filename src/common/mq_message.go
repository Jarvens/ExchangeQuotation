// auth: kunlun
// date: 2019-01-22
// description:
package common

import (
	"config"
	"fmt"
	"time"
	//"time"
)

func init() {
	go func() {
		msgs, err := config.Channel.Consume("go.queue1", "", true, false, false, false, nil)
		if err != nil {
			fmt.Printf("Consume Queue: %s is Fail\n", "go.queue1")
			return
		}

		readChan := make(chan bool)

		go func() {
			for message := range msgs {
				fmt.Printf("\r Receive message: %d %s", time.Now().UnixNano(), message.Body)
			}
		}()

		<-readChan
	}()
}

// 发送消息

func kline() {
	//now:=time.Now().Unix()

}
