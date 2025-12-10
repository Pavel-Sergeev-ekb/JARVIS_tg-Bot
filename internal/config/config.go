package config

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken   string
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
}

var Cfg Config

func LoadConfig() (Config, error) {

	_, err := os.Stat(".env")
	if err == nil {
		// Файл .env найден — загружаем его
		err = godotenv.Load(".env")
		if err != nil {
			return Config{}, fmt.Errorf("ошибка загрузки .env: %v", err)
		}
	} else if !os.IsNotExist(err) {
		// Если ошибка не «файл не найден», а другая (например, права) — сообщаем
		return Config{}, fmt.Errorf("ошибка проверки .env: %v", err)
	}

	cfg := Config{
		BotToken:   os.Getenv("TELEGRAM_BOT_TOKEN"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBName:     os.Getenv("DB_NAME"),
	}
	if cfg.DBHost == "" {
		cfg.DBHost = "localhost"
	}
	if cfg.DBPort == "" {
		cfg.DBPort = "5432"
	}
	if cfg.DBName == "" {
		cfg.DBName = "memory"
	}

	if cfg.BotToken == "" {
		log.Fatal("Токен не найден в переменных окружения")
		return Config{}, errors.New("missing telegram bot token")
	}

	if cfg.DBUser == "" || cfg.DBPassword == "" || cfg.DBName == "" {
		log.Fatal("Не заданы обязательные параметры БД")
		return Config{}, fmt.Errorf("missing database credentials")
	}
	return cfg, nil

}
