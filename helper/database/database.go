package database

import (
	"dating-apps/helper/config"
	"errors"
)

const (
	SqliteDBDriver string = "sqlite"
)

type Option struct {
	Migrate []interface{}
}

type Database interface {
	Client() interface{}
}

func NewDatabaseConnect(cfg *config.DBConfig, opt *Option) (Database, error) {
	switch cfg.Driver {
	case SqliteDBDriver:
		return NewGormConnect(cfg, opt)
	}
	return nil, errors.New("the database driver does not support")
}
