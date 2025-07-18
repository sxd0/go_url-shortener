package configs

import (
	"os"
)

type Config struct {
	Db   *Db
	Auth AuthConfig
}

type AuthConfig struct {
	Secret string
}

func LoadConfig() *Config {
	return &Config{
		Db: &Db{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
		},
		Auth: AuthConfig{
			Secret: os.Getenv("SECRET"),
		},
	}
}
