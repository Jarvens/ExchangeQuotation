// auth: kunlun
// date: 2019-01-24
// description:
package handle

import (
	"config"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"time"
)

// 监听消息队列
func init() {
	queues := config.Queues
	for _, val := range queues {
		go func() {
			messages, err := config.Channel.Consume(val, "", true, false, false, false, nil)
			if err != nil {
				fmt.Printf("Consume queue: %s faild\n", val)
			}

			readChan := make(chan bool)

			go func() {
				for message := range messages {
					fmt.Printf("Receive message: %d - %s", time.Now().Unix(), message.Body)
					var data interface{}
					err = json.Unmarshal(message.Body, &data)
					if err != nil {
						fmt.Printf("Marshal mq message faild")
					}

					m := data.(map[string]interface{})
					for k, v := range m {
						if k == "cmd" {
							if v == "tick" {
								TickHandle(message.Body)
							}
						}
					}
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
