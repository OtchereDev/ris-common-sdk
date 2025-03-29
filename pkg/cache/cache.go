package cache

import (
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type Cache interface {
	UpdateStatus(event proto.Message, role string)
	IsUnSafe(id uint32, role string) bool
	UpdateDeleted(event proto.Message, role string)
	LoadHistoricalUserStatus(js nats.JetStreamContext, service string) error
	LoadHistoricalUserDeleted(js nats.JetStreamContext, service string) error
	LoadHistoricalDoctorStatus(js nats.JetStreamContext, service string) error
	LoadHistoricalDoctorDeleted(js nats.JetStreamContext, service string) error
}
