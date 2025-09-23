package sqlite

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/OtchereDev/ris-common-sdk/pkg/db"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SqliteDB wraps the GORM database client for SQLite
type SqliteDB struct {
	Client *gorm.DB
}

// SqliteConnectionParam holds all SQLite connection parameters
type SqliteConnectionParam struct {
	Name                 string            // Database file path or ":memory:" for in-memory
	MaxIdleConns         int               // Maximum number of idle connections
	MaxOpenConns         int               // Maximum number of open connections
	ConnMaxLifetime      time.Duration     // Maximum connection lifetime
	ConnMaxIdleTime      time.Duration     // Maximum connection idle time
	CreateDirIfNotExists bool              // Create directory if it doesn't exist
	PragmaSettings       map[string]string // SQLite PRAGMA settings
}

// Ping tests the database connection
func (d *SqliteDB) Ping(ctx context.Context) error {
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
func (d *SqliteDB) Migrate(dst ...interface{}) error {
	if err := d.Client.AutoMigrate(dst...); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	return nil
}

// Close closes the database connection
func (d *SqliteDB) Close() error {
	sqlDB, err := d.Client.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	return nil
}

// Connect establishes a connection to SQLite database with connection pooling
func Connect(ctx context.Context, p SqliteConnectionParam) (*SqliteDB, error) {
	d := &SqliteDB{}

	// Validate connection parameters
	if p.Name == "" {
		return nil, errors.New("database name/path is required")
	}

	// Create directory if needed (for file-based databases)
	if p.Name != ":memory:" && p.CreateDirIfNotExists {
		dir := filepath.Dir(p.Name)
		if dir != "." && dir != "" {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return nil, fmt.Errorf("failed to create database directory: %w", err)
			}
		}
	}

	// Prepare GORM config
	config := &gorm.Config{}

	// Open database connection
	client, err := gorm.Open(sqlite.Open(p.Name), config)
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
	// Note: SQLite doesn't benefit from high connection counts due to file locking
	maxIdle := p.MaxIdleConns
	if maxIdle == 0 {
		maxIdle = 1 // SQLite works best with single connection for writes
	}

	maxOpen := p.MaxOpenConns
	if maxOpen == 0 {
		maxOpen = 1 // SQLite handles concurrent reads but single writer
	}

	connLifetime := p.ConnMaxLifetime
	if connLifetime == 0 {
		connLifetime = time.Hour
	}

	connIdleTime := p.ConnMaxIdleTime
	if connIdleTime == 0 {
		connIdleTime = 10 * time.Minute
	}

	// Set connection pool parameters
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetConnMaxLifetime(connLifetime)
	sqlDB.SetConnMaxIdleTime(connIdleTime)

	// Apply PRAGMA settings if provided
	if len(p.PragmaSettings) > 0 {
		for pragma, value := range p.PragmaSettings {
			if err := d.Client.Exec(fmt.Sprintf("PRAGMA %s = %s", pragma, value)).Error; err != nil {
				return nil, fmt.Errorf("failed to set PRAGMA %s: %w", pragma, err)
			}
		}
	} else {
		// Apply default PRAGMA settings for better performance and reliability
		defaultPragmas := map[string]string{
			"journal_mode": "WAL",    // Write-Ahead Logging for better concurrency
			"synchronous":  "NORMAL", // Balance between safety and performance
			"cache_size":   "-64000", // 64MB cache
			"foreign_keys": "ON",     // Enable foreign key constraints
			"temp_store":   "MEMORY", // Store temporary tables in memory
		}

		for pragma, value := range defaultPragmas {
			if err := d.Client.Exec(fmt.Sprintf("PRAGMA %s = %s", pragma, value)).Error; err != nil {
				return nil, fmt.Errorf("failed to set default PRAGMA %s: %w", pragma, err)
			}
		}
	}

	return d, nil
}

// ConnectInMemory creates an in-memory SQLite database
func ConnectInMemory(ctx context.Context) (*SqliteDB, error) {
	params := SqliteConnectionParam{
		Name: ":memory:",
	}
	return Connect(ctx, params)
}

// ConnectFile connects to a file-based SQLite database
func ConnectFile(ctx context.Context, filepath string) (*SqliteDB, error) {
	params := SqliteConnectionParam{
		Name:                 filepath,
		CreateDirIfNotExists: true,
	}
	return Connect(ctx, params)
}

// ConnectWithWAL connects to SQLite with WAL mode enabled for better concurrency
func ConnectWithWAL(ctx context.Context, filepath string, maxOpenConns int) (*SqliteDB, error) {
	if maxOpenConns == 0 {
		maxOpenConns = 5 // Allow more connections with WAL mode
	}

	params := SqliteConnectionParam{
		Name:                 filepath,
		MaxOpenConns:         maxOpenConns,
		MaxIdleConns:         2,
		CreateDirIfNotExists: true,
		PragmaSettings: map[string]string{
			"journal_mode":   "WAL",
			"synchronous":    "NORMAL",
			"cache_size":     "-64000",
			"foreign_keys":   "ON",
			"temp_store":     "MEMORY",
			"wal_checkpoint": "TRUNCATE",
		},
	}
	return Connect(ctx, params)
}

// GetStats returns database connection statistics
func (d *SqliteDB) GetStats() (*db.ConnectionStats, error) {
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

// GetPragma retrieves a SQLite PRAGMA setting
func (d *SqliteDB) GetPragma(pragma string) (string, error) {
	var value string
	if err := d.Client.Raw(fmt.Sprintf("PRAGMA %s", pragma)).Scan(&value).Error; err != nil {
		return "", fmt.Errorf("failed to get PRAGMA %s: %w", pragma, err)
	}
	return value, nil
}

// SetPragma sets a SQLite PRAGMA setting
func (d *SqliteDB) SetPragma(pragma, value string) error {
	if err := d.Client.Exec(fmt.Sprintf("PRAGMA %s = %s", pragma, value)).Error; err != nil {
		return fmt.Errorf("failed to set PRAGMA %s: %w", pragma, err)
	}
	return nil
}

// Vacuum runs VACUUM command to reclaim space and defragment the database
func (d *SqliteDB) Vacuum() error {
	if err := d.Client.Exec("VACUUM").Error; err != nil {
		return fmt.Errorf("failed to vacuum database: %w", err)
	}
	return nil
}

// Analyze runs ANALYZE command to update query planner statistics
func (d *SqliteDB) Analyze() error {
	if err := d.Client.Exec("ANALYZE").Error; err != nil {
		return fmt.Errorf("failed to analyze database: %w", err)
	}
	return nil
}

// ConnectionStats holds database connection pool statistics
