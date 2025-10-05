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

type OperationCode struct {
	id   int
	keys string
	name string
	sms  string
}

func (b *Bot) ShowOperationInfo(chatID int64, operationCode string) error {
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
				id
        keys,
        name,
        sms
    FROM code
    WHERE keys = $1
    `

	var operation OperationCode
	err = db.QueryRow(context.Background(), query, operationCode).Scan(
		&operation.id,
		&operation.keys,
		&operation.name,
		&operation.sms,
	)

	if errors.Is(err, sql.ErrNoRows) {
		b.SendMessage(chatID, fmt.Sprintf("Операция с кодом %s не найдена", operationCode))
		return fmt.Errorf("операция не найдена в базе данных: %w", err)
	}

	if err != nil {
		log.Printf("Ошибка при выполнении запроса: %v", err)
		b.SendMessage(chatID, "Произошла ошибка при получении данных")
		return err
	}

	message := fmt.Sprintf(
		"id операции: %d\n"+
			"Код операции: %s\n"+
			"Название: %s\n\n"+
			"Описание:\n%s",
		operation.id,
		operation.keys,
		operation.name,
		operation.sms,
	)
	b.SendMessage(chatID, message)
	return nil
}
