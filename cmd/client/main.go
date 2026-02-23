package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {

	fmt.Println("Starting Peril client...")

	const connStr = "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(connStr)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	fmt.Println("Connected to RabbitMQ successfully!")

	userName, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Didn't recieve a valid user name: %v", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("could not create server channel: %v", err)
	}
	warPubCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("could not create war publish channel: %v", err)
	}
	defer warPubCh.Close()
	logCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("could not create log channel: %v", err)
	}
	defer logCh.Close()

	state := gamelogic.NewGameState(userName)
	if err := subPause(state, conn, userName); err != nil {
		log.Fatalf("subPause failed: %v", err)
	}
	if err := subArmyMoves(state, warPubCh, conn, userName); err != nil {
		log.Fatalf("subArmyMoves failed: %v", err)
	}
	if err := subWar(state, logCh, conn); err != nil {
		log.Fatalf("subWar failed: %v", err)
	}
outerloop:
	for {
		inputs := gamelogic.GetInput()

		switch inputs[0] {
		case "spawn":
			err := state.CommandSpawn(inputs)
			if err != nil {
				fmt.Printf("unable to spawn properly: %v", err)
			}

		case "move":
			mv, err := state.CommandMove(inputs)
			if err != nil {
				fmt.Printf("unable to move properly: %v", err)
			} else {
				key := routing.ArmyMovesPrefix + "." + mv.Player.Username
				pubsub.PublishJSON(channel, routing.ExchangePerilTopic, key, mv)
				for _, unit := range mv.Units {
					fmt.Printf("Moved unit %d to %s\n", unit.ID, mv.ToLocation)
				}
			}
		case "status":
			state.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			if len(inputs) < 2 {
				fmt.Println("Incorrect spam format: spam <Number of time you want to spam>")
				continue
			}
			amount, err := strconv.Atoi(inputs[1])
			if err != nil {
				fmt.Println("Incorrect spam format: spam <Number of time you want to spam>")
				continue
			}
			for i := 0; i < amount; i++ {
				msg := gamelogic.GetMaliciousLog()
				pubGameLog(logCh, userName, msg)
			}

		case "quit":
			gamelogic.PrintQuit()
			break outerloop
		default:
			fmt.Println("Unknown command given!")
		}

	}
	// Wait for Ctrl+C to shut down the server
	/*
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

		<-sigCh
	*/
	fmt.Println("Shutting down...")
}
