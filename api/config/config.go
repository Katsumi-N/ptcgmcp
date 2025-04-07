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
	Server         Server
	DB             DBConfig
	ESConfig       ESConfig
	FrontendConfig FrontendConfig
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

// Removed the leftover 'type JWT struct {' line

type FrontendConfig struct {
	BaseUrl string `envconfig:"FRONTEND_BASE_URL"`
}

// Restoring the ESConfig struct definition
type ESConfig struct {
	EsHost     string `envconfig:"ES_HOST"`
	EsPort     string `envconfig:"ES_PORT"`
	EsProtocol string `envconfig:"ES_PROTOCOL"`
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
