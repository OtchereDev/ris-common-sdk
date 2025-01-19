package nats

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
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

func (n Nat) Subscribe(q string, subjects []string, handler func(eventType string, data map[string]interface{})) (err error) {

	for _, subject := range subjects {
		subject := subject
		sub, err := n.Jet.QueueSubscribe(subject, q, func(m *nats.Msg) {
			var event struct {
				Type string                 `json:"type"`
				Data map[string]interface{} `json:"data"`
			}

			if err := json.Unmarshal(m.Data, &event); err != nil {
				log.Printf("Error parsing message on subject %s: %v", m.Subject, err)
				return
			}

			fmt.Printf("Received by %s from %s: %s\n",
				q, m.Subject, string(m.Data))

			// Call the provided handler function
			handler(event.Type, event.Data)

			m.Ack()
		}, nats.Durable(fmt.Sprintf("durable-%s", subject)), nats.ManualAck())
		if err != nil {
			log.Fatalf("Error subscribing to subject %s: %v", subject, err)
			return err
		}
		defer sub.Unsubscribe()
	}

	select {}
}

func (n Nat) Publish(subject string, data map[string]interface{}) (err error) {
	message, err := json.Marshal(data)
	if err != nil {
		return
	}

	_, err = n.Jet.Publish(subject, message)
	if err != nil {
		return
	}

	return
}
