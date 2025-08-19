package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv  string
	AppPort string

	DBHost string
	DBPort int
	DBUser string
	DBPass string
	DBName string
	DBSSL  string
}

func Load() *Config {
	_ = godotenv.Load() // ignore error in local; use real env in prod
	port, _ := strconv.Atoi(getenv("DB_PORT", "5432"))

	return &Config{
		AppEnv:  getenv("APP_ENV", "local"),
		AppPort: getenv("APP_PORT", "8080"),

		DBHost: getenv("DB_HOST", "localhost"),
		DBPort: port,
		DBUser: getenv("DB_USER", "todo_user"),
		DBPass: getenv("DB_PASS", "todo_pass"),
		DBName: getenv("DB_NAME", "todo_db"),
		DBSSL:  getenv("DB_SSLMODE", "disable"),
	}
}

func getenv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}

func MustLoad() *Config {
	cfg := Load()
	if cfg == nil {
		log.Fatal("failed to load config")
	}
	return cfg
}
