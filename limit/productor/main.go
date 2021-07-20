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

	exchangeName := "test_qos_exchange"
	routingKey := "qos.save"

	for i := 10; i < 30; i++ {
		err = channel.Publish(exchangeName, routingKey, true, false, amqp.Publishing{
			Body: []byte(fmt.Sprintf("hello rabbitmq qos test %v", i)),
		})
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

/*
	Return 消息机制
Return Listener 用于处理一些不可路由的消息
我门的消息生产者，通过指定一个Exchange和routingKey，把消息送达到某一个队列中去，然后我门的消费者监听队列，进行消费处理操作
某些情况下，发送消息时，当前exchange不存在key路由不到，我门需要监听这种不可达的消息

Mandatory true, 监听器会接收到路由不可达的消息,然后进行后续处理，如果为false, 那么broker端自动删除该消息
*/

/*
no_ack = false 才能生效
prefetchSize: 0 body大小限制
prefetchCount: 会告诉rabbitMQ 不要同时给一个消费者推送多个消息，一旦有N个消息还没ack, 则该consumer将block掉，直到消息ack
global： true,false 是否将上面设置应用于channel, 就是上面限制是channel级别的还是consumer级别的
*/
