package bot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/Pavel-Sergeev-ekb/JARVIS_tg-Bot/internal/api/database"
	"github.com/Pavel-Sergeev-ekb/JARVIS_tg-Bot/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type OperationCode struct {
	id   int
	keys string
	name string
	sms  string
}

func (b *Bot) getOperationInfo(chatID int64, operationKey string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
		return err
	}

	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Printf("Ошибка подключения к БД: %v", err)
		b.SendMessage(chatID, "Произошла ошибка при подключении к базе данных", tgbotapi.ModeHTML)
		return err
	}
	defer db.Close(context.Background())

	query := `
    SELECT 
        id,
        keys,
        name,
        sms
    FROM code
    WHERE keys = $1
    `

	var operation OperationCode
	err = db.QueryRow(context.Background(), query, operationKey).Scan(
		&operation.id,
		&operation.keys,
		&operation.name,
		&operation.sms,
	)

	if errors.Is(err, sql.ErrNoRows) {
		msg := fmt.Sprintf("Операция с кодом %s не найдена", operationKey)
		b.SendMessage(chatID, msg, tgbotapi.ModeHTML)
		return fmt.Errorf("операция не найдена в базе данных: %w", err)
	}

	if err != nil {
		log.Printf("Ошибка при выполнении запроса: %v", err)
		b.SendMessage(chatID, "Произошла ошибка при получении данных", tgbotapi.ModeHTML)
		return err
	}

	message := fmt.Sprintf(
		"<b>Код операции в WMS</b>: <b>%s</b>\n"+
			"<b>Название</b>: <b>%s</b>\n\n"+
			"<b>Описание инструмента</b>:\n%s",

		operation.keys,
		operation.name,
		operation.sms,
	)
	b.SendMessage(chatID, message, tgbotapi.ModeHTML)
	return nil
}

// Конкретные функции для каждой операции
func (b *Bot) ControlPoints(chatID int64, operationCode string) error {
	return b.getOperationInfo(chatID, "1-1")
}

func (b *Bot) CanceledOrders(chatID int64, operationCode string) error {
	return b.getOperationInfo(chatID, "2-7")
}

// Аналогично для остальных операций:
func (b *Bot) Selection20(chatID int64, operationCode string) error {
	return b.getOperationInfo(chatID, "3-3-8")
}

func (b *Bot) AutoSelection20(chatID int64, operationCode string) error {
	return b.getOperationInfo(chatID, "3-3-9")
}
func (b *Bot) Door(chatID int64, operationCode string) error {
	return b.getOperationInfo(chatID, "3-7-2")
}
func (b *Bot) Shipped(chatID int64, operationCode string) error {
	return b.getOperationInfo(chatID, "3-7-3")
}
func (b *Bot) Packing(chatID int64, operationCode string) error {
	return b.getOperationInfo(chatID, "3-10")
}
func (b *Bot) WaveDetails(chatID int64, operationCode string) error {
	return b.getOperationInfo(chatID, "4-1-2-1-2")
}
func (b *Bot) Dispatcher(chatID int64, operationCode string) error {
	return b.getOperationInfo(chatID, "4-1-2-3-3")
}
func (b *Bot) HandStartInfo(chatID int64, operationCode string) error {
	return b.getOperationInfo(chatID, "4-1-2-8-1-1")
}
func (b *Bot) TransactionSKU(chatID int64, operationCode string) error {
	return b.getOperationInfo(chatID, "4-9-1")
}
func (b *Bot) TransactionUIT(chatID int64, operationCode string) error {
	return b.getOperationInfo(chatID, "4-9-2")
}
func (b *Bot) Balance(chatID int64, operationCode string) error {
	return b.getOperationInfo(chatID, "5-1")
}
func (b *Bot) BalanceUIT(chatID int64, operationCode string) error {
	return b.getOperationInfo(chatID, "5-2-2")
}

func (b *Bot) CheckOrders(chatID int64, operationCode string) error {
	return b.getOperationInfo(chatID, "5-3")
}
func (b *Bot) ControlPointsWMS(chatID int64, operationCode string) error {
	return b.getOperationInfo(chatID, "7-2-1")
}
func (b *Bot) Locking(chatID int64, operationCode string) error {
	return b.getOperationInfo(chatID, "9-5")
}
func (b *Bot) AutoStart(chatID int64, operationCode string) error {
	return b.getOperationInfo(chatID, "4-1-2-1-1")
}
func (b *Bot) Action(chatID int64, operationCode string) error {
	return b.getOperationInfo(chatID, "9-7")
}
func (b *Bot) Skills(chatID int64, operationCode string) error {
	return b.getOperationInfo(chatID, "9-2")
}
func (b *Bot) AssignedWork(chatID int64, operationCode string) error {
	return b.getOperationInfo(chatID, "4-11")
}
func (b *Bot) Hiring(chatID int64, operationCode string) error {
	return b.getOperationInfo(chatID, "9-6")
}
