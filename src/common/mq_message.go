// auth: kunlun
// date: 2019-01-22
// description:
package common

import (
	"config"
	"fmt"
	"github.com/streadway/amqp"
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
				fmt.Printf("Receive message: %s\n", message.Body)
			}
		}()

		<-readChan
	}()
}

// 发送消息
func PublishMessage(exchange, queue string, message []byte) {
	config.Channel.Publish(exchange, queue, false, false, amqp.Publishing{ContentType: "text/plain", Body: message})
}

func readMQ() {
	_, err := config.Channel.Consume("go.queue3", "", true, false, false, false, nil)
	if err != nil {
		fmt.Printf("Consume Queue: %s is Fail\n", "go.queue3")
		return
	}
	fmt.Printf("Consume Queue3")
}

func kline() {
	//now:=time.Now().Unix()

}
