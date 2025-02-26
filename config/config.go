package config

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig   `envPrefix:"APP_"`
	Database DatabaseConfig `envPrefix:"DB_"`
	JWT      JWTConfig      `envPrefix:"JWT_"`
}

type ServerConfig struct {
	Port string `env:"PORT"`
}

type DatabaseConfig struct {
	Host     string `env:"HOST"`
	Port     string `env:"PORT"`
	Username string `env:"USER"`
	Password string `env:"PASSWORD"`
	Database string `env:"NAME"`
}

type JWTConfig struct {
	Secret     string `env:"SECRET"`
	Expiration int    `env:"EXPIRATION"`
}

var Cfg Config

func Load() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	if err := env.Parse(&Cfg); err != nil {
		log.Fatal("Error parsing environment variables: ", err)
	}
}

func GetConfig() *Config {
	return &Cfg
}
