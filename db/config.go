package db

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type ClientConfigFormat struct {
	Driver             string        `json:"driver" yaml:"driver" validate:"required"`
	Address            string        `json:"address" yaml:"address" validate:"required"`
	Source             string        `json:"source" yaml:"source" validate:"required"`
	MaxIdle            int           `json:"max_idle" yaml:"max_idle" validate:"required"`
	MaxOpen            int           `json:"max_open" yaml:"max_open" validate:"required"`
	MaxLifetimeSeconds time.Duration `json:"max_lifetime_seconds" yaml:"max_lifetime_seconds" validate:"required"`
}

func NewDB(config *ClientConfigFormat) (sqlxDB *sqlx.DB, err error) {
	sqlxDB, err = sqlx.Connect(config.Driver, config.Source)
	if nil != err {
		logrus.Errorf("connect to %s:%s failed. %s", config.Driver, config.Source, err)
		return
	}

	sqlxDB.SetMaxIdleConns(config.MaxIdle)
	sqlxDB.SetMaxOpenConns(config.MaxOpen)
	sqlxDB.SetConnMaxLifetime(config.MaxLifetimeSeconds * time.Second)

	err = sqlxDB.Ping()
	if nil != err {
		logrus.Errorf("ping to db failed. %s.", err)
		return
	}

	return
}
