package db

import (
	"context"
	"time"
)

type ConnectionStats struct {
	MaxOpenConnections int
	OpenConnections    int
	InUse              int
	Idle               int
	WaitCount          int64
	WaitDuration       time.Duration
	MaxIdleClosed      int64
	MaxLifetimeClosed  int64
}

type Database interface {
	Ping(ctx context.Context) error
	Migrate(dst ...interface{}) error
	Close() error
	GetStats() (*ConnectionStats, error)
}
