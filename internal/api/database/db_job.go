package database

import (
	"context"
	"log"

	"github.com/Pavel-Sergeev-ekb/JARVIS_tg-Bot/internal/config"
)

func SaveUser(chatID int64, userName string) error {
	// Подключение к БД внутри функции (лучше использовать пул подключений)
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Ошибка загрузки конфигурации: %v", err)
		return err
	}

	db, err := ConnectDB(cfg)
	if err != nil {
		log.Printf("Ошибка подключения к БД: %v", err)
		return err
	}
	defer db.Close(context.Background())

	query := `
        INSERT INTO users (chat_id, user_name)
        VALUES ($1, $2)
        ON CONFLICT (chat_id) DO NOTHING
    `

	_, err = db.Exec(context.Background(), query, chatID, userName)
	if err != nil {
		log.Printf("Ошибка при сохранении пользователя chat_id=%d, user_name=%s: %v", chatID, userName, err)
		return err
	}

	log.Printf("Пользователь успешно сохранён: chat_id=%d, user_name=%s", chatID, userName)
	return nil
}
