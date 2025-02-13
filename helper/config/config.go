package config

import (
	"fmt"
	"time"
)

type Config struct {
	Url            UrlConfig      `mapstructure:"url"`
	ServerConfig   ServerConfig   `mapstructure:"server"`
	SecurityConfig SecurityConfig `mapstructure:"security"`
	DBConfig       DBConfig       `mapstructure:"database"`
	AppConfig      AppConfig      `mapstructure:"app"`
}

type DBConfig struct {
	Driver                string        `mapstructure:"driver"`
	Host                  string        `mapstructure:"host"`
	Port                  int           `mapstructure:"port"`
	Username              string        `mapstructure:"username"`
	Password              string        `mapstructure:"password"`
	DBName                string        `mapstructure:"dbname"`
	SchemaName            string        `mapstructure:"schemaname"`
	MaxIdleConnection     int           `mapstructure:"max-idle-connections" default:"20"`
	MaxOpenConnection     int           `mapstructure:"max-open-connections" default:"100"`
	ConnectionMaxLifeTime time.Duration `mapstructure:"connection-max-lifetime" default:"1200"`
	ConnectionMaxIdleTime time.Duration `mapstructure:"connection-max-idle-time" default:"1"`
	LogConfig             DBLogConfig   `mapstructure:"logger"`
}

type DBLogConfig struct {
	Level          string        `mapstructure:"level" default:"info"`
	SlowThreshold  time.Duration `mapstructure:"slow-threshold" default:"200"`
	IgnoreNotFound bool          `mapstructure:"ignore-not-found" default:"true"`
}

type UrlConfig struct {
	Basepath   string `mapstructure:"basepath"`
	Baseprefix string `mapstructure:"baseprefix"`
}

type ServerConfig struct {
	Port      int       `mapstructure:"port"`
	Env       string    `mapstructure:"env"`
	LogConfig LogConfig `mapstructure:"log"`
}

type LogConfig struct {
	Level          string `mapstructure:"level" default:"info"`
	LogOutput      string `mapstructure:"output"`
	OutputFilePath string `mapstructure:"file-path"`
}

type SecurityConfig struct {
	JwtConfig JwtConfig `mapstructure:"jwt"`
}

type JwtConfig struct {
	JwtSecret      string `mapstructure:"jwt-secret"`
	JwtExpireHours int64  `mapstructure:"jwt-expire-hours"`
}

type AppConfig struct {
	SwipeLimit int `mapstructure:"swipe-limit"`
}

func (c *Config) BasePrefix() string {
	return c.Url.Baseprefix
}

func (c *Config) UrlWithPrefix(url string) string {
	return fmt.Sprintf("%s%s", c.BasePrefix(), url)
}
