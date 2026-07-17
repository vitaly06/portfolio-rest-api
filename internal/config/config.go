package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	YandexKey  string
	SMTPHost   string
	SMTPPort   string
	SMTPUser   string
	SMTPPass   string
	OwnerEmail string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		fmt.Println("[CONFIG INFO] Локальный файл .env не найден, читаем переменные из окружения ОС")
	}

	return &Config{
		Port:       getEnv("PORT", "8080"),
		YandexKey:  getEnv("YANDEX_API_KEY", ""),
		SMTPHost:   getEnv("SMTP_HOST", ""),
		SMTPPort:   getEnv("SMTP_PORT", "587"),
		SMTPUser:   getEnv("SMTP_USER", ""),
		SMTPPass:   getEnv("SMTP_PASSWORD", ""),
		OwnerEmail: getEnv("OWNER_EMAIL", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
