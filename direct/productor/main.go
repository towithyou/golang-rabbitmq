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
	routineKeyName := "test.direct"

	for i := 0; i <= 5; i++ {
		err = channel.Publish(exchangeName, routineKeyName, false, false, amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(fmt.Sprintf("hello rabbitMQ %v", i)),
		})

		if err != nil {
			fmt.Println(err)
			return
		}
	}

}
