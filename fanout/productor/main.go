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

	exchangeName := "test_fanout_exchange"

	err = channel.Publish(exchangeName, "", false, false, amqp.Publishing{
		DeliveryMode:    2,       // 持久化
		ContentEncoding: "UTF-8", // 字符集
		ContentType:     "text/plain",
		Expiration:      "10000", // 10秒过期
		Headers:         amqp.Table{"my1": "111", "my2": "222"},
		Body:            []byte(fmt.Sprintf("hello rabbitMQ %v", "save")),
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	err = channel.Publish(exchangeName, "dd", false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(fmt.Sprintf("hello rabbitMQ %v", "update")),
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	err = channel.Publish(exchangeName, "yy", false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(fmt.Sprintf("hello rabbitMQ %v", "delete")),
	})

	if err != nil {
		fmt.Println(err)
		return
	}

}
