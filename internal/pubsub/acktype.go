package pubsub

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Acktype string

const (
	Ack         Acktype = "Ack"
	NackRequeue Acktype = "NackRequeue"
	NackDiscard Acktype = "NackDiscard"
)

func applyAck(msg amqp.Delivery, at Acktype) {
	var err error

	switch at {
	case Ack:
		err = msg.Ack(false)
		fmt.Println("Ack")
	case NackRequeue:
		err = msg.Nack(false, true)
		fmt.Println("NackRequeue")
	case NackDiscard:
		err = msg.Nack(false, false)
		fmt.Println("NackDiscard")
	}

	if err != nil {
		fmt.Printf("ACK/NACK ERROR (%s): %v\n", at, err)
	}
}
