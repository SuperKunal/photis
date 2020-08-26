package config

import (
	"github.com/caarlos0/env"
	"log"
)

type Config struct {
	ConnectionString  string `env:"ConnectionString" envDefault:"mongodb://localhost:27017/?serverSelectionTimeoutMS=5000&connectTimeoutMS=10000"`
	Database  string `env:"Database" envDefault:"Photis"`
}

func NewConfig() *Config {
	config := Config{}
	if err := env.Parse(&config); err != nil {
		log.Fatal(err)
	}
	return &config
}