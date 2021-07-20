package main

import (
	"fmt"
	"github.com/streadway/amqp"
)

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

	exchangeName := "test_direct_exchange"
	exchangeType := "direct"
	queueName := "test_direct_queue"
	routineKeyName := "test.direct"

	// 声明交换机
	err = channel.ExchangeDeclare(exchangeName, exchangeType,
		true, false, false, false, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	// 声明队列
	// durable 持久化，机器重启队列不会消失
	// exclusive 独占，顺序消费，
	// autoDelete 脱离交换机，自动删除
	queue, err := channel.QueueDeclare(
		queueName, true, false,
		false, false, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 建立一个绑定关系
	err = channel.QueueBind(queueName, routineKeyName, exchangeName, false, nil)
	if err != nil {
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
