package config

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

const DotEnvFilename = ".env"

type Config struct {
	ServiceName string           `envconfig:"SERVICE_NAME" default:"post-service"`
	HTTP        HTTPServerConfig `envconfig:"HTTP"`
	Storage     StorageConfig    `envconfig:"STORAGE"`
}

type HTTPServerConfig struct {
	ListenAddr        string        `envconfig:"LISTEN_ADDR" default:":8080"`
	ReadHeaderTimeout time.Duration `envconfig:"READ_HEADER_TIMEOUT" default:"5s"`
	KeepaliveTime     time.Duration `envconfig:"KEEPALIVE_TIME" default:"60s"`
}

type StorageConfig struct {
	Driver string `envconfig:"DRIVER" default:"memory"`
	DSN    string `envconfig:"POSTGRES_DSN"`
}

func NewConfigFromEnv() (*Config, error) {
	_ = godotenv.Load(DotEnvFilename)
	cfg := &Config{}
	if err := envconfig.Process("", cfg); err != nil {
		return nil, errors.Wrap(err, "unable to process env variables")
	}
	return cfg, nil
}
