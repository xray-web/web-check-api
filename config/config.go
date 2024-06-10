package config

import (
	"cmp"
	"fmt"
	"os"
)

type Config struct {
	Host          string
	Port          string
	AllowedOrigin string
}

func New() Config {
	host := getEnvDefault("HOST", "0.0.0.0")
	port := getEnvDefault("PORT", "8080")
	return Config{
		Host:          host,
		Port:          port,
		AllowedOrigin: getEnvDefault("ALLOWED_ORIGINS", fmt.Sprintf("http://%s:%s", host, port)),
	}
}

func getEnvDefault(key, def string) string {
	return cmp.Or(os.Getenv(key), def)
}
