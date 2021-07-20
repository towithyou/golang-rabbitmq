package main

import (
	"fmt"
	"github.com/streadway/amqp"
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

	exchangeName := "test_topic_exchange"
	exchangeType := "topic"
	queueName := "test_topic_queue"
	//routineKeyName := "user.#"
	routineKeyName := "user.*"

	// 声明交换机
	if err = channel.ExchangeDeclare(exchangeName, exchangeType,
		false, false, false, false, nil); err != nil {
		fmt.Println(err)
		return
	}

	// 声明队列
	// durable 持久化，机器重启队列不会消失
	// exclusive 独占，顺序消费，
	// autoDelete 脱离交换机，自动删除
	queue, err := channel.QueueDeclare(
		queueName, false, false,
		false, false, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 建立一个绑定关系 交换机和队列的绑定关系
	if err = channel.QueueBind(queueName, routineKeyName,
		exchangeName, false, nil); err != nil {
		fmt.Println(err)
		return
	}

	deliveries, err := channel.Consume(queue.Name, "",
		true, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		for d := range deliveries {
			fmt.Printf("Received a message: %s\n", d.Body)
		}
	}()

	select {}

}
