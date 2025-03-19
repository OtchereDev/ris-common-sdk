package cache

import (
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type Cache interface {
	UpdateStatus(event proto.Message)
	IsUnSafe(id uint32) bool
	UpdateDeleted(event proto.Message)
	LoadHistoricalStatus(js nats.JetStreamContext, service string) error
	LoadHistoricalDeleted(js nats.JetStreamContext, service string) error
}
