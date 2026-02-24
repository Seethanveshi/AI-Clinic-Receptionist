package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl string
	Port  string
}

func Load() *Config {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error while load environment variables.", err)
	}

	dbURL := "postgres://" +
		getEnv("DB_USER") + ":" +
		getEnv("DB_PASSWORD") + "@" +
		getEnv("DB_HOST") + ":" +
		getEnv("DB_PORT") + "/" +
		getEnv("DB_NAME") +
		"?sslmode=" + getEnv("DB_SSLMODE")
		
	return &Config{
		DBUrl: dbURL,
		Port:  getEnv("PORT"),
	}
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Missing required environment variable: %s", key)
	}
	return value
}
