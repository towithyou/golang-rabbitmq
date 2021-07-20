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

	exchangeName := "test_topic_exchange"
	routingKey1 := "user.save"
	routingKey2 := "user.update"
	routingKey3 := "user.delete.tom"

	err = channel.Publish(exchangeName, routingKey1, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(fmt.Sprintf("hello rabbitMQ %v", "save")),
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	err = channel.Publish(exchangeName, routingKey2, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(fmt.Sprintf("hello rabbitMQ %v", "update")),
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	err = channel.Publish(exchangeName, routingKey3, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(fmt.Sprintf("hello rabbitMQ %v", "delete")),
	})

	if err != nil {
		fmt.Println(err)
		return
	}

}
