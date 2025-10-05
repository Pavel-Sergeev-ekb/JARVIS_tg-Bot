package database

import (
	"context"
	"log"
)

func SaveUser(chatID int64, userName string) {
	_, err := db.Exec(
		context.Background(),
		"INSERT INTO users (chat_id, user_name) VALUES ($1, $2) ON CONFLICT (chat_id) DO NOTHING",
		chatID,
		userName,
	)
	if err != nil {
		log.Println("Ошибка при сохранении пользователя:", err)
	} else {
		log.Println("Пользователь сохранен:", userName)
	}
}
