package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type JWTConfig struct {
	Secret     string `yaml:"secret"`
	Expiration string `yaml:"expiration"`
}

var Cfg Config

func Load() {
	file, err := os.ReadFile("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	err = yaml.Unmarshal(file, &Cfg)
	if err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	log.Println("Config loaded successfully")
}

func GetConfig() *Config {
	return &Cfg
}
