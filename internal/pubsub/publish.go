package pubsub

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	bytes, err := json.Marshal(val)
	mandatory := false
	immediate := false
	if err != nil {
		return fmt.Errorf("Failled to Marshal the value %v", err)
	}

	///option 1
	err = ch.PublishWithContext(context.Background(), exchange, key, mandatory, immediate, amqp.Publishing{
		ContentType: "application/json",
		Body:        bytes,
	})
	if err != nil {
		return fmt.Errorf("failled to publish: %v", err)
	}

	return nil
}

func PublishGob[T any](ch *amqp.Channel, exchange, key string, val T) error {
	var buff bytes.Buffer
	err := gob.NewEncoder(&buff).Encode(val)

	body := buff.Bytes()
	mandatory := false
	immediate := false
	if err != nil {
		return fmt.Errorf("Failled to Encode Gob the value %v", err)
	}

	err = ch.PublishWithContext(context.Background(), exchange, key, mandatory, immediate, amqp.Publishing{
		ContentType: "application/gob",
		Body:        body,
	})
	if err != nil {
		return fmt.Errorf("failled to publish: %v", err)
	}
	fmt.Printf("hello Exchange:%s  with key %s executed", exchange, key)

	return nil
}
