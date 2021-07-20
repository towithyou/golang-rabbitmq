package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"time"
)

/*
# 匹配一个或多个词
* 匹配不多不少一个词

log.# 能够匹配到 log.info.oa
log.* 只能匹配到 log.error
*/

func main() {
	url := "amqp://guest:guest@127.0.0.1:5672/"
	connection, err := amqp.Dial(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer channel.Close()

	exchangeName := "test_dlx_exchange"
	routingKey := "dlx.#"
	queueName := "test_dlx_queue"

	err = channel.ExchangeDeclare(exchangeName, "topic",
		true, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 死信队列的声明， 在创建队列参数
	queue, err := channel.QueueDeclare(queueName, false,
		false, false, false, amqp.Table{
			"x-dead-letter-exchange": "dlx.exchange", // 死信队列
		})
	if err != nil {
		fmt.Println(err)
		return
	}

	// 创建死信队列
	if err = channel.ExchangeDeclare("dlx.exchange",
		"topic", true, false, false, false, nil); err != nil {
		fmt.Println(err)
		return
	}

	qu, err := channel.QueueDeclare("dlx.queue", false, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err = channel.QueueBind(qu.Name, "#", "dlx.exchange", false, nil); err != nil {
		fmt.Println(err)
		return
	}

	if err = channel.QueueBind(queue.Name, routingKey, exchangeName, false, nil); err != nil {
		fmt.Println(err)
		return
	}
	//

	// 限流方式， autoAck=false
	if err = channel.Qos(3, 0, false); err != nil {
		fmt.Println(err)
		return
	}

	deliveries, err := channel.Consume(queue.Name, "",
		false, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		for d := range deliveries {
			fmt.Printf("Received a message: %s\n", d.Body)
			if err := d.Nack(false, false); err != nil {
				fmt.Println(err)
				return
			}
			time.Sleep(3 * time.Second)
		}
	}()

	select {}

}

/*
死信队列
1. 消息被拒绝(reject/nack) 并且 requeue = false
2. 消息ttl过期
3. 队列达到最大长度

1. DLX是正常的exchange，和一般的exchange没有区别, 他能在任何的队列上被指定，实际上设置某个队列的属性
2. 当这个队列中有死信时，mq就会自动的将这个消息重新发布到设置的exchange上去， 进而被路由到另一个队列
3. 可以监听这个队列中消息做相应的处理, 这个特性弥补mq3.0 之前支持的immediate参数功能
*/
