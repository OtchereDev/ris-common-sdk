package nats

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type Nat struct {
	Client *nats.Conn
	Jet    nats.JetStreamContext
}

type NatSubjects struct {
	Subject string
	Handler func(eventType string, m *nats.Msg)
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

func sanitizeConsumerName(name string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9-_]`)
	sanitized := re.ReplaceAllString(name, "-")
	return sanitized
}

func (n Nat) Subscribe(q string, subjects []NatSubjects) (err error) {

	var subs []*nats.Subscription = []*nats.Subscription{}

	for _, subject := range subjects {
		s := subject.Subject
		sub, err := n.Jet.QueueSubscribe(s, q, func(m *nats.Msg) {
			subject.Handler(s, m)
		}, nats.Durable(fmt.Sprintf("durable-%s", sanitizeConsumerName(s))), nats.ManualAck())

		if err != nil {
			// Clean up any subscriptions already created before returning error
			for _, createdSub := range subs {
				createdSub.Unsubscribe()
			}
			return fmt.Errorf("error subscribing to subject %s: %v", s, err)
		}

		subs = append(subs, sub)
	}

	defer func() {
		for _, sub := range subs {
			sub.Unsubscribe()
		}
	}()

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
