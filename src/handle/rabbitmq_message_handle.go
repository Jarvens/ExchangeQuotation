// auth: kunlun
// date: 2019-01-24
// description:
package handle

import (
	"config"
	"fmt"
	"github.com/streadway/amqp"
	"time"
)

// 监听消息队列
func init() {
	queues := config.Queues
	for _, val := range queues {
		go func() {
			msgs, err := config.Channel.Consume(val, "", true, false, false, false, nil)
			if err != nil {
				fmt.Printf("Consume queue: %s faild\n", val)
			}

			readChan := make(chan bool)

			go func() {
				for message := range msgs {
					fmt.Printf("Receive message: %d - %s", time.Now().Unix(), message.Body)
				}
			}()

			<-readChan
		}()
	}
}

// 推送消息
func PublishMessage(exchange, queue string, message []byte) error {
	err := config.Channel.Publish(exchange, queue, false, false, amqp.Publishing{ContentType: "text/plain", Body: message})
	if err != nil {
		return err
	}
	return nil

}
