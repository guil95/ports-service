package config

import (
	"context"
	"log"

	_ "github.com/joho/godotenv/autoload"
	"github.com/sethvargo/go-envconfig"
)

var AppConfig MyConfig

type MyConfig struct {
	Environment string `env:"ENVIRONMENT"`
	DBHost      string `env:"DB_HOST"`
	DBPort      string `env:"DB_PORT"`
	DBUser      string `env:"DB_USER"`
	DBName      string `env:"DB_NAME"`
	DBPassword  string `env:"DB_PASSWORD"`
	DBSSLMode   string `env:"DB_SSL_MODE"`
}

func init() {
	if err := envconfig.Process(context.Background(), &AppConfig); err != nil {
		log.Fatal(err)
	}
}
