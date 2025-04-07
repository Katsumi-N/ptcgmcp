package config

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Server      Server
	DB          DBConfig
	MeiliConfig MeiliConfig
}

type DBConfig struct {
	Name     string `envconfig:"DB_DATABASE"`
	User     string `envconfig:"DB_USER"`
	Password string `envconfig:"DB_PASS"`
	Port     string `envconfig:"DB_PORT"`
	Host     string `envconfig:"DB_HOST"`
}

type Server struct {
	Address string `envconfig:"ADDRESS"`
	Port    string `envconfig:"PORT"`
}

// MeiliConfig Meilisearchの設定
type MeiliConfig struct {
	Host     string `envconfig:"MEILI_HOST"`
	Port     string `envconfig:"MEILI_PORT"`
	Protocol string `envconfig:"MEILI_PROTOCOL"`
	ApiKey   string `envconfig:"MEILI_API_KEY"`
}

var (
	once   sync.Once
	config Config
)

func GetConfig() *Config {
	if os.Getenv("GO_ENV") == "" {
		os.Setenv("GO_ENV", "development")
	}

	if err := godotenv.Load(fmt.Sprintf(".env.%s", os.Getenv("GO_ENV"))); err != nil {
		log.Printf("No .env file found")
	}

	once.Do(func() {
		if err := envconfig.Process("", &config); err != nil {
			panic(err)
		}
	})
	return &config
}
