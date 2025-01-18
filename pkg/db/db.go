package db

type Database interface {
	Ping() error
	Migrate(dst ...interface{}) error
}
