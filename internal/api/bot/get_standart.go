package bot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Pavel-Sergeev-ekb/JARVIS_tg-Bot/internal/api/database"
	"github.com/Pavel-Sergeev-ekb/JARVIS_tg-Bot/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Stand struct {
	id          int
	name        string
	indicator   int64
	user_chatid int
	data_at     time.Time
}

func (b *Bot) getStand(chatID int64, id string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Ошибка загрузки конфигурации: %v", err)
		b.SendMessage(chatID, "Произошла ошибка при загрузке конфигурации", tgbotapi.ModeHTML)
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
        SELECT s.*, u.user_name
        FROM standart s
        JOIN users u ON s.user_chatid = u.chat_id
        WHERE s.id = $1
    `

	var S Stand
	var username string
	err = db.QueryRow(context.Background(), query, id).Scan(
		&S.id,
		&S.name,
		&S.indicator,
		&S.user_chatid,
		&S.data_at,
		&username,
	)

	if errors.Is(err, sql.ErrNoRows) {
		msg := fmt.Sprintf("Показатель %s не найден", id)
		b.SendMessage(chatID, msg, tgbotapi.ModeHTML)
		return fmt.Errorf("показатель не найден в базе данных: %w", err)
	}

	if err != nil {
		log.Printf("Ошибка при выполнении запроса: %v", err)
		b.SendMessage(chatID, "Произошла ошибка при получении данных", tgbotapi.ModeHTML)
		return err
	}

	message := formatStanMessage(S, username)
	b.SendMessage(chatID, message, tgbotapi.ModeHTML)
	return nil
}

func formatStanMessage(S Stand, username string) string {
	return fmt.Sprintf(
		"<b>Название:</b> <b>%s</b>\n"+
			"<b>Показатель:</b> <b>%d</b>\n"+
			"<b>Обновил:</b> <b>%s</b>\n"+
			"<b>Последнее обновление:</b> <b>%s</b>\n",
		S.name,
		S.indicator,
		username,
		S.data_at.Format("02.01.2006 15:04:05"),
	)
}

func (b *Bot) Mez1Info(chatID int64, id string) error {
	return b.getStand(chatID, "1")
}
func (b *Bot) Mez2Info(chatID int64, id string) error {
	return b.getStand(chatID, "2")
}
func (b *Bot) Mez3Info(chatID int64, id string) error {
	return b.getStand(chatID, "3")
}
func (b *Bot) Mez4Info(chatID int64, id string) error {
	return b.getStand(chatID, "4")
}
func (b *Bot) Mez5Info(chatID int64, id string) error {
	return b.getStand(chatID, "5")
}
func (b *Bot) SpecInfo(chatID int64, id string) error {
	return b.getStand(chatID, "9")
}
func (b *Bot) SelectionKGTInfo(chatID int64, id string) error {
	return b.getStand(chatID, "7")
}
func (b *Bot) SelectionSGTInfo(chatID int64, id string) error {
	return b.getStand(chatID, "8")
}
func (b *Bot) SelectionIZInfo(chatID int64, id string) error {
	return b.getStand(chatID, "6")
}
func (b *Bot) SortKGTInfo(chatID int64, id string) error {
	return b.getStand(chatID, "10")
}
func (b *Bot) SortSGTInfo(chatID int64, id string) error {
	return b.getStand(chatID, "11")
}
func (b *Bot) SortLotInfo(chatID int64, id string) error {
	return b.getStand(chatID, "12")
}
func (b *Bot) PackPotInfo(chatID int64, id string) error {
	return b.getStand(chatID, "13")
}
func (b *Bot) ShipNormaInfo(chatID int64, id string) error {
	return b.getStand(chatID, "14")
}
func (b *Bot) BufNormaInfo(chatID int64, id string) error {
	return b.getStand(chatID, "15")
}
func (b *Bot) SelectionBalkon(chatID int64, id string) error {
	return b.getStand(chatID, "16")
}
