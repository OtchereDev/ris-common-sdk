package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/OtchereDev/ris-common-sdk/pkg/db"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PgDB wraps the GORM database client
type PgDB struct {
	Client *gorm.DB
}

// PgConnectionParam holds all database connection parameters
type PgConnectionParam struct {
	Name            string
	Host            string
	Port            string
	User            string
	Password        string
	URL             string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// Ping tests the database connection
func (d *PgDB) Ping(ctx context.Context) error {
	sqlDB, err := d.Client.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

// Migrate runs GORM auto-migration for the given models
func (d *PgDB) Migrate(dst ...interface{}) error {
	if err := d.Client.AutoMigrate(dst...); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	return nil
}

// Close closes the database connection
func (d *PgDB) Close() error {
	sqlDB, err := d.Client.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	return nil
}

// Connect establishes a connection to PostgreSQL database with connection pooling
func Connect(ctx context.Context, p PgConnectionParam) (*PgDB, error) {
	d := &PgDB{}

	// Validate connection parameters
	if p.URL == "" {
		if p.Host == "" || p.User == "" || p.Password == "" || p.Port == "" || p.Name == "" {
			return nil, errors.New("provide the correct database credentials or connection URL")
		}
	}

	var dsn string
	if p.URL != "" {
		dsn = p.URL
	} else {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=prefer",
			p.Host, p.User, p.Password, p.Name, p.Port)
	}

	// Open database connection
	client, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	d.Client = client

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := d.Client.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Test the connection
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Configure connection pool with defaults
	maxIdle := p.MaxIdleConns
	if maxIdle == 0 {
		maxIdle = 10
	}

	maxOpen := p.MaxOpenConns
	if maxOpen == 0 {
		maxOpen = 100
	}

	connLifetime := p.ConnMaxLifetime
	if connLifetime == 0 {
		connLifetime = time.Hour
	}

	connIdleTime := p.ConnMaxIdleTime
	if connIdleTime == 0 {
		connIdleTime = 30 * time.Minute
	}

	// Set connection pool parameters
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetConnMaxLifetime(connLifetime)
	sqlDB.SetConnMaxIdleTime(connIdleTime)

	return d, nil
}

// ConnectWithDefaults connects to PostgreSQL with default connection pool settings
func ConnectWithDefaults(ctx context.Context, host, port, user, password, dbname string) (*PgDB, error) {
	params := PgConnectionParam{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Name:     dbname,
	}
	return Connect(ctx, params)
}

// ConnectWithURL connects to PostgreSQL using a connection URL
func ConnectWithURL(ctx context.Context, url string) (*PgDB, error) {
	params := PgConnectionParam{
		URL: url,
	}
	return Connect(ctx, params)
}

// GetStats returns database connection statistics
func (d *PgDB) GetStats() (*db.ConnectionStats, error) {
	sqlDB, err := d.Client.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	stats := sqlDB.Stats()
	return &db.ConnectionStats{
		MaxOpenConnections: stats.MaxOpenConnections,
		OpenConnections:    stats.OpenConnections,
		InUse:              stats.InUse,
		Idle:               stats.Idle,
		WaitCount:          stats.WaitCount,
		WaitDuration:       stats.WaitDuration,
		MaxIdleClosed:      stats.MaxIdleClosed,
		MaxLifetimeClosed:  stats.MaxLifetimeClosed,
	}, nil
}
