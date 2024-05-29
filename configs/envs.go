package configs

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	DBUser     string
	DBPassword string
	DBHost     string
	DBName     string
	DBPort     string
}

var Envs = load()

func load() Config {
	godotenv.Load()
	return Config{
		Port:       getEnv("API_PORT", "5000"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBName:     getEnv("DB_NAME", "datasil"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}

		return i
	}
	return fallback
}
