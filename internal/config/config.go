package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	Env      string   `yaml:"env"`
	Server   Server   `yaml:"server"`
	Postgres Postgres `yaml:"postgres"`
	Redis    Redis    `yaml:"redis"`
}

type Server struct {
	Host string `yaml:"host" env:"HOST" env-default:"localhost"`
	Port string `yaml:"port" env:"PORT" env-default:"3000"`
}

type Postgres struct {
	User     string `yaml:"user" env:"DB_USER" env-default:"postgres"`
	Password string `yaml:"password" env:"DB_PASSWORD"`
	Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port     string `yaml:"port" env:"DB_PORT" env-default:"5432"`
	Name     string `yaml:"name" env:"DB_NAME" env-default:"postgres"`
	SSLMode  bool   `yaml:"ssl_mode" env:"DB_SSL" env-default:"false"`
}

type Redis struct {
	Name     int    `yaml:"name" env:"REDIS_NAME" env-default:"0"`
	Host     string `yaml:"host" env:"REDIS_HOST" env-default:"localhost"`
	Port     string `yaml:"port" env:"REDIS_PORT" env-default:"6379"`
	User     string `yaml:"user" env:"REDIS_USER"`
	Password string `yaml:"Password" env:"REDIS_PASSWORD"`
}

func Load() *Config {
	var cfg Config

	err := cleanenv.ReadConfig("config.yml", &cfg)

	if err != nil {
		log.Fatalf("error while read config: %v", err)
	}

	return &cfg
}