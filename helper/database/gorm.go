package database

import (
	"database/sql"
	"dating-apps/helper/config"
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type gormdb struct {
	db *gorm.DB
}

func NewGormConnect(cfg *config.DBConfig, opt *Option) (Database, error) {
	g := &gormdb{}

	var dialect gorm.Dialector
	switch cfg.Driver {
	case SqliteDBDriver:
		dialect = g.sqliteOpen(cfg)
	}

	db, err := gorm.Open(dialect, g.options(cfg))
	if err != nil {
		return nil, err
	}

	if len(opt.Migrate) > 0 {
		e := db.AutoMigrate(opt.Migrate...)
		if e != nil {
			return nil, e
		}
	}

	//assign the db to struct
	g.db = db

	//connection Pool
	_, err = g.connPool(cfg)
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (m *gormdb) Client() interface{} {
	return m.db
}

func (m *gormdb) connPool(cfg *config.DBConfig) (*sql.DB, error) {
	sqlDB, err := m.db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConnection)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConnection)
	sqlDB.SetConnMaxIdleTime(time.Second * cfg.ConnectionMaxIdleTime)
	sqlDB.SetConnMaxLifetime(time.Second * cfg.ConnectionMaxLifeTime)

	return sqlDB, nil
}

func (m *gormdb) options(cfg *config.DBConfig) *gorm.Config {
	schemaName := ""
	if cfg.SchemaName != "" {
		schemaName = cfg.SchemaName + "."
	}

	return &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			TablePrefix:   schemaName,
		},
		Logger: m.logger(cfg.LogConfig),
	}
}

func (m *gormdb) sqliteOpen(cfg *config.DBConfig) gorm.Dialector {
	return sqlite.Open(cfg.DBName)
}

func (m *gormdb) logger(cfg config.DBLogConfig) logger.Interface {
	var logLevel logger.LogLevel
	switch cfg.Level {
	case "silent":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "info":
		logLevel = logger.Info
	default:
		logLevel = logger.Warn
	}

	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             cfg.SlowThreshold * time.Millisecond,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: cfg.IgnoreNotFound,
		},
	)
}
