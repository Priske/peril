package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func handlerMove(gs *gamelogic.GameState, channel *amqp.Channel) func(gamelogic.ArmyMove) pubsub.Acktype {

	return func(mv gamelogic.ArmyMove) pubsub.Acktype {
		defer fmt.Print("> ")
		outcome := gs.HandleMove(mv)
		fmt.Printf("\n[DEBUG move] user=%s outcome=%v movePlayer=%s to=%s\n",
			gs.GetUsername(), outcome, mv.Player.Username, mv.ToLocation)
		switch outcome {
		case gamelogic.MoveOutComeSafe:
			return pubsub.Ack
		case gamelogic.MoveOutcomeMakeWar:
			routingKey := routing.WarRecognitionsPrefix + "." + gs.GetUsername()

			battle := gamelogic.RecognitionOfWar{
				Attacker: mv.Player,
				Defender: gs.GetPlayerSnap(),
			}
			err := pubsub.PublishJSON(channel, routing.ExchangePerilTopic, routingKey, battle)
			if err != nil {
				fmt.Printf("failed to publish war message: %v\n", err)
				return pubsub.NackRequeue // or requeue; but discard makes debugging easier
			}
			return pubsub.Ack

		default:
			return pubsub.NackDiscard
		}

	}
}
