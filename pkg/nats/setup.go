package nats

import (
	"errors"
	"fmt"
	"regexp"
	"time"

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

type NatSetupOptions struct {
	Queue    string
	Subjects []string
}

func Connect(conn string) (n Nat, err error) {
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

	return
}

func Setup(conn string, options []NatSetupOptions) (n Nat, err error) {
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

	for _, o := range options {

		streamInfo, e := js.StreamInfo(o.Queue)
		if e != nil && e != nats.ErrStreamNotFound {
			return n, e
		}

		if streamInfo != nil {
			return n, e
		} else {
			// Add the stream if it doesn't exist
			_, err = js.AddStream(&nats.StreamConfig{
				Name:      o.Queue,
				Subjects:  o.Subjects,
				Retention: nats.LimitsPolicy,
				Storage:   nats.FileStorage,
				MaxAge:    time.Nanosecond * 2592000000000000,
			})
			if err != nil {
				return
			}
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
		},
			nats.Durable(fmt.Sprintf("durable-%s", sanitizeConsumerName(s))),
			nats.ManualAck(),
			nats.DeliverAll(),
			nats.MaxDeliver(5),
			nats.AckWait(30*time.Second),
		)

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
