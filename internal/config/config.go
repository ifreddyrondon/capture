package config

import (
	"io"
	"log"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

const defaultAddr = "127.0.0.1:8080"

type Constants struct {
	ADDR string
	PG   string
}

// Source set the configuration source in case you aren't allowed to read a file.
func Source(source io.Reader) func(*Config) {
	return func(cfg *Config) {
		cfg.source = source
	}
}

type Config struct {
	Constants
	Database *gorm.DB
	source   io.Reader
}

// NewConfig is used to generate a configuration instance which will be passed around the codebase
func New(opts ...func(*Config)) (*Config, error) {
	var cfg Config

	for _, opt := range opts {
		opt(&cfg)
	}

	constants, err := initViper(&cfg)
	cfg.Constants = constants
	if err != nil {
		return &cfg, err
	}
	db, err := gorm.Open("postgres", cfg.Constants.PG)
	if err != nil {
		log.Panic(err)
	}
	cfg.Database = db
	return &cfg, err
}

// OnShutdown is executed as graceful shutdown.
func (cfg *Config) OnShutdown() {
	log.Printf("[finalizer:data source] closing the main session")
	if err := cfg.Database.Close(); err != nil {
		log.Fatal(err)
	}
}

func initViper(cfg *Config) (Constants, error) {
	viper.SetDefault("ADDR", defaultAddr)

	var err error
	if cfg.source != nil {
		viper.SetConfigType("toml")
		err = viper.ReadConfig(cfg.source)
	} else {
		viper.SetConfigName("config")
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
