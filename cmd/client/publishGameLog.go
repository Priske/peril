package main

import (
	"time"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func pubGameLog(channel *amqp.Channel, username string, msg string) error {
	exchange := routing.ExchangePerilTopic
	routingKey := routing.GameLogSlug + "." + username
	gl := routing.GameLog{
		CurrentTime: time.Now(),
		Message:     msg,
		Username:    username,
	}

	if err := pubsub.PublishGob(channel, exchange, routingKey, gl); err != nil {
		return err
	}
	return nil
}
