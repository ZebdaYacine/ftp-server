package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	Port      string
	Host      string
	UploadDir string
	BaseURL   string
}

func Load() (Env, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using system environment variables")
	}

	env := Env{
		Port:      lookup("SERVER_PORT", ""),
		Host:      lookup("SERVER_HOST", ""),
		UploadDir: lookup("UPLOAD_DIR", "/var/www/aisha"),
		BaseURL:   lookup("BASE_URL", "http://localhost:8080/aisha"),
	}

	return env, nil
}

func lookup(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
