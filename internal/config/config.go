package config

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken string
}

func LoadConfig() (Config, error) {

	_, err := os.Stat(".env")
	if os.IsNotExist(err) {
		return Config{}, fmt.Errorf(".env файл не найден. Создайте его на основе .env.example")
	}

	err = godotenv.Load(".env")

	if err != nil {
		log.Fatalf("ошибка загрузки файла .env: %v", err)
		return Config{}, err
	}

	cfg := Config{
		BotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
	}

	if cfg.BotToken == "" {
		log.Fatal("Токен не найден в переменных окружения")
		return Config{}, errors.New("missing telegram bot token")
	}
	return cfg, nil

}
