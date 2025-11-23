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

type OrderStatus struct {
	statusID   int
	statusCode string
	statusName string
	statusText string
}

// Общая функция для получения информации о статусе заказа
func (b *Bot) getOrderStatusInfo(chatID int64, statusID string) error {
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

	if errors.Is(err, sql.ErrNoRows) {
		b.SendMessage(chatID, "Данные не найдены", tgbotapi.ModeHTML)
		return fmt.Errorf("данные не найдены: %w", err)
	}

	if err != nil {
		log.Printf("Ошибка при выполнении запроса: %v", err)
		b.SendMessage(chatID, "Произошла ошибка при получении данных", tgbotapi.ModeHTML)
		return err
	}

	message := fmt.Sprintf(
		"<b>Код статуса:</b> <b>%s</b>\n"+
			"<b>Название:</b> <b>%s</b>\n"+
			"<b>Описание:</b> %s\n",
		data.statusCode,
		data.statusName,
		data.statusText,
	)
	b.SendMessage(chatID, message, tgbotapi.ModeHTML)
	return nil

}

func (b *Bot) createdExt(chatID int64) error {
	return b.getOrderStatusInfo(chatID, "3")
}
func (b *Bot) released(chatID int64) error {
	return b.getOrderStatusInfo(chatID, "4")
}
func (b *Bot) sorted(chatID int64) error {
	return b.getOrderStatusInfo(chatID, "6")
}
func (b *Bot) packed(chatID int64) error {
	return b.getOrderStatusInfo(chatID, "7")
}
func (b *Bot) sortedSD(chatID int64) error {
	return b.getOrderStatusInfo(chatID, "8")
}

func (b *Bot) shipment(chatID int64) error {
	return b.getOrderStatusInfo(chatID, "9")
}
func (b *Bot) KIZ(chatID int64) error {
	return b.getOrderStatusInfo(chatID, "2")
}

func (b *Bot) selectionCompleted(chatID int64) error {
	return b.getOrderStatusInfo(chatID, "5")
}

func (b *Bot) Unknown(chatID int64) error {
	return b.getOrderStatusInfo(chatID, "1")
}
