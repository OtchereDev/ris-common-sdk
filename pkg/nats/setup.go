package nats

import (
	"errors"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type Nat struct {
	Client *nats.Conn
	Jet    nats.JetStreamContext
}

func Connect(conn, queue string, subject []string) (n Nat, err error) {
	nc, err := nats.Connect(conn)
	if err != nil {
		return
	}

	js, err := nc.JetStream()
	if err != nil {
		return
	}

	n = Nat{
		Client: nc,
		Jet:    js,
	}

	streamInfo, err := js.StreamInfo(queue)
	if err != nil && err != nats.ErrStreamNotFound {
		return
	}

	if streamInfo != nil {
		return
	} else {
		// Add the stream if it doesn't exist
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     queue,
			Subjects: subject,
		})
		if err != nil {
			return
		}
	}

	return
}

func (n Nat) Subscribe(q string, subjects []string, handler func(eventType string, m *nats.Msg)) (err error) {

	for _, subject := range subjects {
		subject := subject
		sub, err := n.Jet.QueueSubscribe(subject, q, func(m *nats.Msg) {
			// Call the provided handler function
			handler(subject, m)
		}, nats.Durable(fmt.Sprintf("durable-%s", subject)), nats.ManualAck())
		if err != nil {
			log.Fatalf("Error subscribing to subject %s: %v", subject, err)
			return err
		}
		defer sub.Unsubscribe()
	}

	select {}
}

func (n Nat) Publish(subject string, protoMsg proto.Message) (err error) {
	if protoMsg == nil {
		return errors.New("proto message cannot be nil")
	}

	message, err := proto.Marshal(protoMsg)
	if err != nil {
		return
	}

	_, err = n.Jet.Publish(subject, message)
	if err != nil {
		return
	}

	return
}
