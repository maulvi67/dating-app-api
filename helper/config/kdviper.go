package config

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

var m = &sync.Mutex{}
var c *Config

type KDOption struct {
	ConfigPath string
	EnvPrefix  string
}

// Option is a function that modifies KDOption.
type Option func(opt *KDOption)

// WithConfigPath allows setting a custom configuration path.
func WithConfigPath(configPath string) Option {
	return func(opt *KDOption) {
		opt.ConfigPath = configPath
	}
}

// Init loads and returns the configuration using Viper.
// It applies any provided options.
func Init(opts ...Option) *Config {
	// Set default options.
	kdOption := &KDOption{
		ConfigPath: ".",
		EnvPrefix:  "KD",
	}
	// Apply any provided options.
	for _, applyOpt := range opts {
		applyOpt(kdOption)
	}

	// Create a new Viper instance.
	v := NewViper(kdOption)

	// Unmarshal the configuration into a local variable.
	var conf Config
	err := v.Unmarshal(&conf)
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	// Save the configuration in the package-level variable.
	c = &conf
	return c
}

// Get returns the singleton configuration, initializing it if necessary.
func Get(opts ...Option) *Config {
	// Use double-check locking to avoid unnecessary locks.
	if c == nil {
		m.Lock()
		defer m.Unlock()
		if c == nil {
			return Init(opts...)
		}
	}
	return c
}

// NewViper creates a new Viper instance based on the provided options.
func NewViper(opt *KDOption) *viper.Viper {
	// Determine the profile based on the environment.
	profile := "dev"
	if env := os.Getenv("KD_ENV"); env == "prd" || env == "stg" {
		profile = "prd"
	}

	// Build the config file name.
	configFileName := "config-" + profile

	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigName(configFileName)
	v.AddConfigPath(opt.ConfigPath)
	v.SetEnvPrefix(opt.EnvPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	// Read in the configuration file.
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	return v
}
