package pubsub

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type decodeFn[T any] func([]byte, *T) error

func subscribe[T any](conn *amqp.Connection, exchange, queueName, key string, queueType SimpleQueueType, decode decodeFn[T], handler func(T) Acktype) error {
	consumer := ""
	autoAck := false
	exclusive := false
	noLocal := false
	noWait := false

	channel, queue, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return fmt.Errorf("was unable to declare and bind: %w", err)
	}
	channel.Qos(10, 0, false)
	deliveryCh, err := channel.Consume(queue.Name, consumer, autoAck, exclusive, noLocal, noWait, nil)
	if err != nil {
		_ = channel.Close()
		return fmt.Errorf("invoking consume channel went wrong: %w", err)
	}

	go func() {
		defer channel.Close()

		for message := range deliveryCh {
			var msg T
			if err := decode(message.Body, &msg); err != nil {
				fmt.Printf("Failed to decode data: %v\n", err)
				// malformed payload is a logical error → discard
				_ = message.Nack(false, false)
				continue
			}

			at := handler(msg)
			applyAck(message, at)
		}
	}()

	return nil
}

// --- Concrete decoders ---

func decodeJSON[T any](b []byte, out *T) error {
	return json.Unmarshal(b, out)
}

func decodeGob[T any](b []byte, out *T) error {
	dec := gob.NewDecoder(bytes.NewReader(b))
	return dec.Decode(out)
}

// --- Public API wrappers ---

func SubscribeJSON[T any](conn *amqp.Connection, exchange, queueName, key string, queueType SimpleQueueType, handler func(T) Acktype) error {
	return subscribe(conn, exchange, queueName, key, queueType, decodeJSON[T], handler)
}

func SubscribeGob[T any](conn *amqp.Connection, exchange, queueName, key string, queueType SimpleQueueType, handler func(T) Acktype) error {
	return subscribe(conn, exchange, queueName, key, queueType, decodeGob[T], handler)
}
