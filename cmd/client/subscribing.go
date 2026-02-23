package main

import (
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func subPause(state *gamelogic.GameState, connection *amqp.Connection, userName string) error {
	exchange := routing.ExchangePerilDirect
	routingKey := routing.PauseKey
	queueType := pubsub.Transient
	queueName := routing.PauseKey + "." + userName
	err := pubsub.SubscribeJSON(connection, exchange, queueName, routingKey, queueType, handlerPause(state))
	return err
}
func subArmyMoves(state *gamelogic.GameState, channel *amqp.Channel, connection *amqp.Connection, userName string) error {
	exchange := routing.ExchangePerilTopic
	routingKey := routing.ArmyMovesPrefix + ".*"
	queueType := pubsub.Transient
	queueName := routing.ArmyMovesPrefix + "." + userName

	err := pubsub.SubscribeJSON(connection, exchange, queueName, routingKey, queueType, handlerMove(state, channel))
	return err
}

func subWar(state *gamelogic.GameState, channel *amqp.Channel, connection *amqp.Connection) error {
	exchange := routing.ExchangePerilTopic
	routingKey := routing.WarRecognitionsPrefix + ".*"
	queueType := pubsub.Durable
	queueName := "war"
	err := pubsub.SubscribeJSON(connection, exchange, queueName, routingKey, queueType, handlerWar(state, channel))
	return err
}
