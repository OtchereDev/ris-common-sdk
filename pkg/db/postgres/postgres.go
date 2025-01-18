package postgres

import (
	"errors"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func (d *PgDB) Ping() (err error) {
	sqlDB, err := d.Client.DB()

	if err != nil {
		return err
	}
	err = sqlDB.Ping()

	return
}

func (d *PgDB) Migrate(dst ...interface{}) (err error) {

	err = d.Client.AutoMigrate(
		dst...,
	)

	return
}

func Connect(p PgConnectionParam) (d *PgDB, err error) {

	d = &PgDB{}

	if p.URL == "" {

		if p.Host == "" || p.User == "" ||
			p.Password == "" || p.Port == "" ||
			p.Name == "" {
			d.Client, err = nil,
				errors.New("provide the correct database credentials")
			return
		}
	}

	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=prefer",
		p.Host, p.User, p.Password, p.Name, p.Port)

	if p.URL != "" {
		d.Client, err = gorm.Open(postgres.Open(p.URL), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: false,
		})
	} else {
		d.Client, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: false,
		})
	}

	return
}

type PgDB struct {
	Client *gorm.DB
}

type PgConnectionParam struct {
	Name     string
	Host     string
	Port     string
	User     string
	Password string
	URL      string
}
