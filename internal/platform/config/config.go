package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName         string
	AppEnv          string
	AppPort         string
	DatabaseURL     string
	RedisAddr       string
	RedisPassword   string
	RedisDB         int
	RabbitMQURL     string
	ShutdownTimeout time.Duration
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		AppName:         getEnv("APP_NAME", "backend-boiler"),
		AppEnv:          getEnv("APP_ENV", "development"),
		AppPort:         getEnv("APP_PORT", "8080"),
		DatabaseURL:     getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/backend_boiler?sslmode=disable"),
		RedisAddr:       getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:   getEnv("REDIS_PASSWORD", ""),
		RedisDB:         getEnvInt("REDIS_DB", 0),
		RabbitMQURL:     getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		ShutdownTimeout: getEnvDuration("SHUTDOWN_TIMEOUT", 10*time.Second),
	}
}

func (c Config) ServerAddress() string {
	return ":" + c.AppPort
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return parsed
}
