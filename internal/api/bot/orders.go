package bot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/Pavel-Sergeev-ekb/JARVIS_tg-Bot/internal/api/database"
	"github.com/Pavel-Sergeev-ekb/JARVIS_tg-Bot/internal/config"
)

type OrderStatus struct {
	statusID   int
	statusCode int
	statusName string
	statusText string
}

func (b *Bot) CreatedExt(chatID int64, statusID string) error {

	if statusID == "" {
		b.SendMessage(chatID, "Ошибка: не указан идентификатор статуса")
		return fmt.Errorf("пустой идентификатор статуса")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
		return err
	}

	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Printf("Ошибка подключения к БД: %v", err)
		b.SendMessage(chatID, "Произошла ошибка при подключении к базе данных")
		return err
	}
	defer db.Close(context.Background())

	query := `
    SELECT 
    status_id,
    status_code,
    status_name,
    status_text
    FROM order_status
    WHERE status_id = $1
    `

	var data OrderStatus
	err = db.QueryRow(context.Background(), query, statusID).Scan(
		&data.statusID,
		&data.statusCode,
		&data.statusName,
		&data.statusText,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			b.SendMessage(chatID, "Данные не найдены")
			return fmt.Errorf("данные не найдены: %w", err)
		}
		log.Printf("Ошибка при выполнении запроса: %v", err)
		b.SendMessage(chatID, "Произошла ошибка при получении данных")
		return err
	}

	message := fmt.Sprintf(
		"Код статуса: %d\n"+
			"Название: %s\n"+
			"Описание: %s\n",
		data.statusCode,
		data.statusName,
		data.statusText,
	)
	b.SendMessage(chatID, message)
	return nil
}
