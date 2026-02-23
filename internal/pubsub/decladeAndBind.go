package pubsub

import (
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType string

const (
	Durable   SimpleQueueType = "durable"
	Transient SimpleQueueType = "transient"
)

// SimpleQueueType is an "enum" type I made to represent "durable" or "transient"
func DeclareAndBind(conn *amqp.Connection, exchange, queueName, key string, queueType SimpleQueueType) (*amqp.Channel, amqp.Queue, error) {
	durable := false
	autoDelete := false
	exclusive := false
	noWait := false
	tableArgs := amqp.Table{
		"x-dead-letter-exchange": "peril_dlx",
	}

	channel, err := conn.Channel()
	if err != nil {
		return &amqp.Channel{}, amqp.Queue{}, err
	}
	switch queueType {
	case "durable":
		durable = true
	case "transient":
		autoDelete = true
		exclusive = true
	default:
		return &amqp.Channel{}, amqp.Queue{}, errors.New("incompatible queueType")
	}

	queue, err := channel.QueueDeclare(queueName, durable, autoDelete, exclusive, noWait, tableArgs)
	if err != nil {

		return &amqp.Channel{}, amqp.Queue{}, err
	}
	if err := channel.QueueBind(queue.Name, key, exchange, noWait, nil); err != nil {
		return &amqp.Channel{}, amqp.Queue{}, err
	}
	return channel, queue, nil
}
