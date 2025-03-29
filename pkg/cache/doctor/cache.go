package doctor

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/OtchereDev/ris-common-sdk/pkg/proto/referral"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

// UserStatusCache maintains an in-memory cache of disabled users
type UserStatusCache struct {
	disabledUsers map[string]int64
	deletedUsers  map[string]int64
	mu            sync.RWMutex
}

// NewUserStatusCache creates a new cache
func NewUserStatusCache() *UserStatusCache {
	return &UserStatusCache{
		disabledUsers: make(map[string]int64),
		deletedUsers:  make(map[string]int64),
	}
}

// UpdateStatus updates the cache based on an event
func (c *UserStatusCache) UpdateStatus(event *referral.ReferringDoctor) {
	c.mu.Lock()
	defer c.mu.Unlock()

	id := strconv.Itoa(int(event.Id))
	if event.IsDisabled {
		c.disabledUsers[id] = event.CreatedAt.Seconds
	} else {
		delete(c.disabledUsers, id)
	}
}

// IsUnSafe checks if a user is disabled or deleted
func (c *UserStatusCache) IsUnSafe(userID uint32) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	id := strconv.Itoa(int(userID))
	_, disabled := c.disabledUsers[id]
	_, deleted := c.deletedUsers[id]
	return disabled || deleted
}

// UpdateDeleted updates the cache based on an event
func (c *UserStatusCache) UpdateDeleted(event *referral.ReferringDoctor) {
	c.mu.Lock()
	defer c.mu.Unlock()

	id := strconv.Itoa(int(event.Id))
	if event.IsDeleted {
		c.deletedUsers[id] = event.CreatedAt.Seconds
	} else {
		delete(c.deletedUsers, id)
	}
}

// LoadHistoricalStatus loads all historical user status events from JetStream
func (c *UserStatusCache) LoadHistoricalStatus(js nats.JetStreamContext, service string) error {
	// Create a consumer that starts from the beginning of the stream
	sub, err := js.Subscribe("referral.doctor.disabled", func(msg *nats.Msg) {
		event := &referral.ReferringDoctor{}
		if err := proto.Unmarshal(msg.Data, event); err != nil {
			log.Printf("Error unmarshalling historical event: %v", err)
			return
		}

		c.UpdateStatus(event)
		msg.Ack()
	},
		nats.DeliverAll(),
		nats.ManualAck(),
		nats.Durable(fmt.Sprintf("%s-historical-status-loader", service)),
	)

	if err != nil {
		return err
	}

	time.Sleep(5 * time.Second)
	return sub.Unsubscribe()
}

// LoadHistoricalDeleted loads all historical user deleted events from JetStream
func (c *UserStatusCache) LoadHistoricalDeleted(js nats.JetStreamContext, service string) error {
	// Create a consumer that starts from the beginning of the stream
	sub, err := js.Subscribe("referral.doctor.deleted", func(msg *nats.Msg) {
		event := &referral.ReferringDoctor{}
		if err := proto.Unmarshal(msg.Data, event); err != nil {
			log.Printf("Error unmarshalling historical event: %v", err)
			return
		}

		c.UpdateDeleted(event)
		msg.Ack()
	},
		nats.DeliverAll(),
		nats.ManualAck(),
		nats.Durable(fmt.Sprintf("%s-historical-delete-loader", service)),
	)

	if err != nil {
		return err
	}

	time.Sleep(5 * time.Second)
	return sub.Unsubscribe()
}
