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

	err = channel.Confirm(false)
	if err != nil {
		fmt.Println(err)
		return
	}

	exchangeName := "test_confirm_exchange"
	routingKey := "confirm.save"

	for i := 0; i < 5; i++ {
		err = channel.Publish(exchangeName, routingKey, false, false, amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(fmt.Sprintf("hello rabbitMQ %v", "save")),
		})

		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// 确定消息是否发布成功
	c := make(chan amqp.Confirmation)
	channel.NotifyPublish(c)

	go func() {

		for r := range c {
			fmt.Println(r.Ack, r.DeliveryTag)
			fmt.Println("ack ?")
		}

	}()

	select {}
}

/*
	Return 消息机制
Return Listener 用于处理一些不可路由的消息
我门的消息生产者，通过指定一个Exchange和routingKey，把消息送达到某一个队列中去，然后我门的消费者监听队列，进行消费处理操作
某些情况下，发送消息时，当前exchange不存在key路由不到，我门需要监听这种不可达的消息

Mandatory true, 监听器会接收到路由不可达的消息,然后进行后续处理，如果为false, 那么broker端自动删除该消息
*/
