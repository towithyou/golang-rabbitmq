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

	exchangeName := "test_return_exchange"
	routingKey := "return.#"
	queueName := "test_return_queue"

	err = channel.ExchangeDeclare(exchangeName, "topic",
		true, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	queue, err := channel.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err = channel.QueueBind(queue.Name, routingKey, exchangeName, false, nil); err != nil {
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
