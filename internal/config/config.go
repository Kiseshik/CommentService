package config

import (
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

const (
	DotEnvFilename  = ".env"
	StorageMemory   = "memory"
	StoragePostgres = "postgres"
)

type Config struct {
	ServiceName       string        `envconfig:"SERVICE_NAME" default:"post-service"`
	ListenAddr        string        `envconfig:"HTTP_LISTEN_ADDR" default:":8080"`
	ReadHeaderTimeout time.Duration `envconfig:"HTTP_READ_HEADER_TIMEOUT" default:"5s"`
	KeepaliveTime     time.Duration `envconfig:"HTTP_KEEPALIVE_TIME" default:"60s"`
	StorageDriver     string        `envconfig:"STORAGE_DRIVER" default:"memory"`
	PostgresDSN       string        `envconfig:"POSTGRES_DSN"`
}

func NewConfigFromEnv() (*Config, error) {
	_ = godotenv.Load(DotEnvFilename) //игнорим чтобы приложуха не упала в проде с ошибкой если там не будет энв файла
	cfg := &Config{}

	if err := envconfig.Process("", cfg); err != nil {
		return nil, errors.Wrap(err, "unable to process env variables")
	}

	cfg.StorageDriver = strings.ToLower(cfg.StorageDriver)
	switch cfg.StorageDriver {
	case StorageMemory:
	case StoragePostgres:
		if cfg.PostgresDSN == "" {
			return nil, errors.New("POSTGRES_DSN is required when STORAGE_DRIVER=postgres")
		}
	default:
		return nil, errors.New("invalid STORAGE_DRIVER (allowed: memory, postgres)")
	}

	return cfg, nil
}

// для сборки в руте
func (c *Config) IsMemoryStorage() bool {
	return c.StorageDriver == StorageMemory
}

func (c *Config) IsPostgresStorage() bool {
	return c.StorageDriver == StoragePostgres
}
