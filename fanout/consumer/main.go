package main

import (
	"fmt"
	"github.com/streadway/amqp"
)

/*
不处理理由健, 只需要简单的将队列绑定到交换机上
发送到交换机的消息都会被转发到与该交换机绑定的所有队列上
fanout 交换机转发消息是最快的
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

	exchangeName := "test_fanout_exchange"
	exchangeType := "fanout"
	queueName := "test_fanout_queue"
	//routineKeyName := "user.#"
	routineKeyName := "" // 不设置路由健

	// 声明交换机
	if err = channel.ExchangeDeclare(exchangeName, exchangeType,
		true, false, false, false, nil); err != nil {
		fmt.Println(err)
		return
	}

	// 声明队列
	// durable 持久化，机器重启队列不会消失
	// exclusive 独占，顺序消费，
	// autoDelete 脱离交换机，自动删除
	// durability 是否持久化, durable 是 transient 否
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
			fmt.Printf("Received a message: %s, headers: %v\n", d.Body, d.Headers)

		}
	}()

	select {}

}
