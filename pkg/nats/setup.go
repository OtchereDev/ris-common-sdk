// package nats

// import (
// 	"errors"
// 	"fmt"
// 	"regexp"
// 	"time"

// 	"github.com/nats-io/nats.go"
// 	"google.golang.org/protobuf/proto"
// )

// type Nat struct {
// 	Client *nats.Conn
// 	Jet    nats.JetStreamContext
// }

// type NatSubjects struct {
// 	Subject string
// 	Handler func(eventType string, m *nats.Msg)
// }

// type NatSetupOptions struct {
// 	Queue    string
// 	Subjects []string
// }

// func Connect(conn string) (n Nat, err error) {
// 	nc, err := nats.Connect(conn)
// 	if err != nil {
// 		return
// 	}

// 	js, err := nc.JetStream()
// 	if err != nil {
// 		return
// 	}

// 	n = Nat{
// 		Client: nc,
// 		Jet:    js,
// 	}

// 	return
// }

// func Setup(conn string, options []NatSetupOptions) (n Nat, err error) {
// 	nc, err := nats.Connect(conn)
// 	if err != nil {
// 		return
// 	}

// 	js, err := nc.JetStream()
// 	if err != nil {
// 		return
// 	}

// 	n = Nat{
// 		Client: nc,
// 		Jet:    js,
// 	}

// 	for _, o := range options {

// 		streamInfo, e := js.StreamInfo(o.Queue)
// 		if e != nil && e != nats.ErrStreamNotFound {
// 			return n, e
// 		}

// 		if streamInfo != nil {
// 			return n, e
// 		} else {
// 			// Add the stream if it doesn't exist
// 			_, err = js.AddStream(&nats.StreamConfig{
// 				Name:      o.Queue,
// 				Subjects:  o.Subjects,
// 				Retention: nats.LimitsPolicy,
// 				Storage:   nats.FileStorage,
// 				MaxAge:    time.Nanosecond * 2592000000000000,
// 			})
// 			if err != nil {
// 				return
// 			}
// 		}
// 	}

// 	return
// }

// func sanitizeConsumerName(name string) string {
// 	re := regexp.MustCompile(`[^a-zA-Z0-9-_]`)
// 	sanitized := re.ReplaceAllString(name, "-")
// 	return sanitized
// }

// func (n Nat) Subscribe(q string, subjects []NatSubjects) (err error) {

// 	var subs []*nats.Subscription = []*nats.Subscription{}

// 	for _, subject := range subjects {
// 		s := subject.Subject
// 		sub, err := n.Jet.QueueSubscribe(s, q, func(m *nats.Msg) {
// 			subject.Handler(s, m)
// 		},
// 			nats.Durable(fmt.Sprintf("durable-%s", sanitizeConsumerName(s))),
// 			nats.ManualAck(),
// 			nats.DeliverAll(),
// 			nats.MaxDeliver(5),
// 			nats.AckWait(30*time.Second),
// 		)

// 		if err != nil {
// 			// Clean up any subscriptions already created before returning error
// 			for _, createdSub := range subs {
// 				createdSub.Unsubscribe()
// 			}
// 			return fmt.Errorf("error subscribing to subject %s: %v", s, err)
// 		}

// 		subs = append(subs, sub)
// 	}

// 	defer func() {
// 		for _, sub := range subs {
// 			sub.Unsubscribe()
// 		}
// 	}()

// 	select {}
// }

// func (n Nat) Publish(subject string, protoMsg proto.Message) (err error) {
// 	if protoMsg == nil {
// 		return errors.New("proto message cannot be nil")
// 	}

// 	message, err := proto.Marshal(protoMsg)
// 	if err != nil {
// 		return
// 	}

// 	_, err = n.Jet.Publish(subject, message)
// 	if err != nil {
// 		return
// 	}

// 	return
// }

package nats

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type Nat struct {
	Client *nats.Conn
	Jet    nats.JetStreamContext
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

type NatSubjects struct {
	Subject string
	Handler func(eventType string, m *nats.Msg)
}

type StreamConfig struct {
	Name     string
	Subjects []string
}

type SubscribeOptions struct {
	ServiceName   string
	Subjects      []NatSubjects
	MaxDeliver    int
	AckWait       time.Duration
	DeliverPolicy nats.DeliverPolicy
}

func Connect(conn string) (n *Nat, err error) {
	nc, err := nats.Connect(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	n = &Nat{
		Client: nc,
		Jet:    js,
		ctx:    ctx,
		cancel: cancel,
	}

	return n, nil
}

func Setup(conn string, streams []StreamConfig) (n *Nat, err error) {
	n, err = Connect(conn)
	if err != nil {
		return nil, err
	}

	for _, stream := range streams {
		if err := n.CreateStream(stream); err != nil {
			n.Close()
			return nil, fmt.Errorf("failed to create stream %s: %w", stream.Name, err)
		}
	}

	return n, nil
}

func (n *Nat) CreateStream(config StreamConfig) error {
	if config.Name == "" {
		return errors.New("stream name cannot be empty")
	}

	if len(config.Subjects) == 0 {
		return errors.New("stream must have at least one subject")
	}

	// Check if stream already exists
	streamInfo, err := n.Jet.StreamInfo(config.Name)
	if err != nil && err != nats.ErrStreamNotFound {
		return fmt.Errorf("error checking stream info: %w", err)
	}

	// If stream exists, check if subjects match
	if streamInfo != nil {
		existingSubjects := streamInfo.Config.Subjects
		if !subjectsEqual(existingSubjects, config.Subjects) {
			log.Printf("Stream %s exists with different subjects. Existing: %v, Requested: %v",
				config.Name, existingSubjects, config.Subjects)
		}
		return nil
	}

	// Create new stream
	_, err = n.Jet.AddStream(&nats.StreamConfig{
		Name:      config.Name,
		Subjects:  config.Subjects,
		Retention: nats.LimitsPolicy,
		Storage:   nats.FileStorage,
		MaxAge:    30 * 24 * time.Hour, // 30 days
		Replicas:  1,
	})

	if err != nil {
		return fmt.Errorf("failed to add stream: %w", err)
	}

	log.Printf("Created stream: %s with subjects: %v", config.Name, config.Subjects)
	return nil
}

func (n *Nat) Subscribe(options SubscribeOptions) error {
	if options.ServiceName == "" {
		return errors.New("service name cannot be empty")
	}

	if len(options.Subjects) == 0 {
		return errors.New("must provide at least one subject to subscribe to")
	}

	// Set defaults
	if options.MaxDeliver == 0 {
		options.MaxDeliver = 5
	}
	if options.AckWait == 0 {
		options.AckWait = 30 * time.Second
	}
	if options.DeliverPolicy == 0 {
		options.DeliverPolicy = nats.DeliverAllPolicy
	}

	// Group subjects by stream
	streamGroups, err := n.groupSubjectsByStream(options.Subjects)
	if err != nil {
		return fmt.Errorf("failed to group subjects by stream: %w", err)
	}

	// Create consumers for each stream
	for streamName, subjects := range streamGroups {
		n.wg.Add(1)
		go func(sName string, sSubjects []NatSubjects) {
			defer n.wg.Done()
			if err := n.subscribeToStream(sName, options.ServiceName, sSubjects, options); err != nil {
				log.Printf("Error subscribing to stream %s: %v", sName, err)
			}
		}(streamName, subjects)
	}

	return nil
}

// Add this enhanced version of subscribeToStream to your nats package

func (n *Nat) subscribeToStream(streamName, serviceName string, subjects []NatSubjects, options SubscribeOptions) error {
	consumerName := fmt.Sprintf("%s-%s-consumer", serviceName, streamName)
	queueGroup := serviceName

	log.Printf("[INFO] Creating consumer %s for stream %s with %d subjects", consumerName, streamName, len(subjects))

	var subs []*nats.Subscription

	for idx, subject := range subjects {
		durableName := consumerName + "-" + sanitizeConsumerName(subject.Subject)

		log.Printf("[INFO] [%d/%d] Subscribing to subject: %s (durable: %s)",
			idx+1, len(subjects), subject.Subject, durableName)

		// Build subscription options
		subOpts := []nats.SubOpt{
			nats.Durable(durableName),
			nats.ManualAck(),
			nats.MaxDeliver(options.MaxDeliver),
			nats.AckWait(options.AckWait),
		}

		// Add deliver policy option
		switch options.DeliverPolicy {
		case nats.DeliverAllPolicy:
			subOpts = append(subOpts, nats.DeliverAll())
		case nats.DeliverLastPolicy:
			subOpts = append(subOpts, nats.DeliverLast())
		case nats.DeliverNewPolicy:
			subOpts = append(subOpts, nats.DeliverNew())
		case nats.DeliverLastPerSubjectPolicy:
			subOpts = append(subOpts, nats.DeliverLastPerSubject())
		default:
			subOpts = append(subOpts, nats.DeliverAll())
		}

		// Wrap the handler with logging
		wrappedHandler := func(subj string, handler func(string, *nats.Msg)) func(*nats.Msg) {
			return func(m *nats.Msg) {
				log.Printf("[DEBUG] Message received - Subject: %s, Size: %d bytes, Reply: %s",
					m.Subject, len(m.Data), m.Reply)

				// Call the original handler
				handler(subj, m)
			}
		}(subject.Subject, subject.Handler)

		sub, err := n.Jet.QueueSubscribe(
			subject.Subject,
			queueGroup,
			wrappedHandler,
			subOpts...,
		)

		if err != nil {
			log.Printf("[ERROR] Failed to subscribe to %s: %v", subject.Subject, err)
			// Clean up any subscriptions already created
			for _, createdSub := range subs {
				if unsubErr := createdSub.Unsubscribe(); unsubErr != nil {
					log.Printf("[ERROR] Failed to unsubscribe during cleanup: %v", unsubErr)
				}
			}
			return fmt.Errorf("error subscribing to subject %s: %w", subject.Subject, err)
		}

		subs = append(subs, sub)
		log.Printf("[SUCCESS] Subscribed to subject: %s in stream: %s", subject.Subject, streamName)
	}

	log.Printf("[INFO] Successfully created %d subscriptions for stream %s", len(subs), streamName)

	// Keep subscriptions alive until context is cancelled
	<-n.ctx.Done()

	log.Printf("[INFO] Shutting down subscriptions for stream %s", streamName)

	// Clean up subscriptions
	for _, sub := range subs {
		if err := sub.Unsubscribe(); err != nil {
			log.Printf("[ERROR] Error unsubscribing: %v", err)
		}
	}

	return nil
}

func (n *Nat) groupSubjectsByStream(subjects []NatSubjects) (map[string][]NatSubjects, error) {
	streamGroups := make(map[string][]NatSubjects)

	for _, subject := range subjects {
		streamName, err := n.findStreamForSubject(subject.Subject)
		if err != nil {
			return nil, fmt.Errorf("no stream found for subject %s: %w", subject.Subject, err)
		}

		streamGroups[streamName] = append(streamGroups[streamName], subject)
	}

	return streamGroups, nil
}

func (n *Nat) findStreamForSubject(subject string) (string, error) {
	streams := n.Jet.StreamNames()

	for streamName := range streams {
		streamInfo, err := n.Jet.StreamInfo(streamName)
		if err != nil {
			continue
		}

		// Check if subject matches any of the stream's subjects
		for _, streamSubject := range streamInfo.Config.Subjects {
			if n.subjectMatches(subject, streamSubject) {
				return streamName, nil
			}
		}
	}

	return "", fmt.Errorf("no stream found for subject: %s", subject)
}

func (n *Nat) subjectMatches(subject, pattern string) bool {
	if pattern == subject {
		return true
	}

	// Handle wildcards like "user.*"
	if strings.HasSuffix(pattern, ".*") {
		prefix := strings.TrimSuffix(pattern, ".*")
		return strings.HasPrefix(subject, prefix+".")
	}

	// Handle single level wildcards like "user.>"
	if strings.HasSuffix(pattern, ".>") {
		prefix := strings.TrimSuffix(pattern, ".>")
		return strings.HasPrefix(subject, prefix+".")
	}

	return false
}

func (n *Nat) Publish(subject string, protoMsg proto.Message) error {
	if protoMsg == nil {
		return errors.New("proto message cannot be nil")
	}

	if subject == "" {
		return errors.New("subject cannot be empty")
	}

	message, err := proto.Marshal(protoMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal proto message: %w", err)
	}

	_, err = n.Jet.Publish(subject, message)
	if err != nil {
		return fmt.Errorf("failed to publish message to subject %s: %w", subject, err)
	}

	log.Printf("Published message to subject: %s", subject)
	return nil
}

func (n *Nat) Close() error {
	if n.cancel != nil {
		n.cancel()
	}

	// Wait for all goroutines to finish
	n.wg.Wait()

	if n.Client != nil {
		n.Client.Close()
	}

	return nil
}

func (n *Nat) HealthCheck() error {
	if n.Client == nil || n.Client.Status() != nats.CONNECTED {
		return errors.New("NATS connection is not healthy")
	}

	// Try to get JetStream account info
	_, err := n.Jet.AccountInfo()
	if err != nil {
		return fmt.Errorf("JetStream not available: %w", err)
	}

	return nil
}

func sanitizeConsumerName(name string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9-_]`)
	sanitized := re.ReplaceAllString(name, "-")
	return sanitized
}

func subjectsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	subjectMap := make(map[string]bool)
	for _, subject := range a {
		subjectMap[subject] = true
	}

	for _, subject := range b {
		if !subjectMap[subject] {
			return false
		}
	}

	return true
}
