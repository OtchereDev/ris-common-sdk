package sqlite

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func (d *SqlLiteDB) Ping() (err error) {
	sqlDB, err := d.Client.DB()

	if err != nil {
		return err
	}
	err = sqlDB.Ping()

	return
}

func (d *SqlLiteDB) Migrate(dst ...interface{}) (err error) {

	err = d.Client.AutoMigrate(
		dst...,
	)

	return
}

func Connect(p SqliteConnectionParam) (d *SqlLiteDB, err error) {

	d = &SqlLiteDB{}

	d.Client, err = gorm.Open(
		sqlite.Open(p.Name),
		&gorm.Config{},
	)

	return
}

type SqlLiteDB struct {
	Client *gorm.DB
}

type SqliteConnectionParam struct {
	Name string
}
