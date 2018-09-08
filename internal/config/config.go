package config

import (
	"io"
	"log"

	"github.com/sarulabs/di"
	"github.com/spf13/viper"
)

const defaultAddr = "127.0.0.1:8080"

type Constants struct {
	ADDR string
	PG   string
}

// Source set the configuration source in case you aren't allowed to read a file.
func Source(source io.Reader) func(*configOpts) {
	return func(cfg *configOpts) {
		cfg.source = source
	}
}

type configOpts struct {
	source io.Reader
}

// Config represent the global app configuration
type Config struct {
	Constants
	Container di.Container
}

// NewConfig is used to generate a configuration instance which will be passed around the codebase
func New(opts ...func(*configOpts)) (*Config, error) {
	var cfgOpts configOpts
	for _, opt := range opts {
		opt(&cfgOpts)
	}

	var cfg Config
	constants, err := initViper(&cfgOpts)
	cfg.Constants = constants
	if err != nil {
		return &cfg, err
	}
	cfg.Container = getResources(&cfg)

	return &cfg, err
}

// OnShutdown is executed as graceful shutdown.
func (cfg *Config) OnShutdown() {
	log.Printf("[finalizer:resources] deleting resources")
	cfg.Container.Delete()
}

func initViper(cfg *configOpts) (Constants, error) {
	viper.SetDefault("ADDR", defaultAddr)

	var err error
	if cfg.source != nil {
		viper.SetConfigType("toml")
		err = viper.ReadConfig(cfg.source)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath("./config")
		viper.AddConfigPath(".")
		err = viper.ReadInConfig()
		viper.AutomaticEnv()
	}

	if err != nil {
		return Constants{}, err
	}

	var constants Constants
	err = viper.Unmarshal(&constants)
	return constants, err
}
