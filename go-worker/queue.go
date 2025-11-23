package main

import (
	amqp "github.com/streadway/amqp"
)

type RabbitQueue struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	Queue   amqp.Queue
}

func initQueue(url string) (*RabbitQueue, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		"jobs",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, err
	}

	return &RabbitQueue{
		Conn:    conn,
		Channel: ch,
		Queue:   q,
	}, nil

}

func (rq *RabbitQueue) Consume() (<-chan amqp.Delivery, error) {

	return rq.Channel.Consume(
		rq.Queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	// channel.consume returns error as second argument ; implicit
}
