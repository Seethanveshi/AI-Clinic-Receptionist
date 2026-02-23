package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl string
	Port string
}

func Load() *Config {
	_ = godotenv.Load()

	dbURL := "postgres://" + 
	 	os.Getenv("DB_USER") + ":" +
		os.Getenv("DB_PASSWORD") + "@" + 
		os.Getenv("DB_HOST") + ":" +
		os.Getenv("DB_PORT") + "/" + 
		os.Getenv("DB_NAME") + 
		"?sslmode=" + os.Getenv("DB_SSLMODE")

	if dbURL == "" {
		log.Fatal("DB configuration missing")
	}

	return &Config{
		DBUrl: dbURL,
		Port: os.Getenv("PORT"),
	}
}