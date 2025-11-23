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

type KPI struct {
	id          int
	name        string
	indicator   float64
	weight      int64
	user_chatid int
	data        time.Time
}

func (b *Bot) getKPI(chatID int64, id string) error {
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
        SELECT k.*, u.user_name
        FROM kpi k
        JOIN users u ON k.user_chatid = u.chat_id
        WHERE k.id = $1
    `

	var kpi KPI
	var username string
	err = db.QueryRow(context.Background(), query, id).Scan(
		&kpi.id,
		&kpi.name,
		&kpi.indicator,
		&kpi.weight,
		&kpi.user_chatid,
		&kpi.data,
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

	message := formatKPIMessage(kpi, username)
	b.SendMessage(chatID, message, tgbotapi.ModeHTML)
	return nil
}

func formatKPIMessage(kpi KPI, username string) string {
	return fmt.Sprintf(
		"<b>Название:</b> <b>%s</b>\n"+
			"<b>Показатель:</b> <b>%.2f</b>\n"+
			"<b>Вес:</b> <b>%d%%</b>\n"+
			"<b>Обновил:</b> <b>%s</b>\n"+
			"<b>Последнее обновление:</b> <b>%s</b>\n",
		kpi.name,
		kpi.indicator,
		kpi.weight,
		username,
		kpi.data.Format("02.01.2006 15:04:05"),
	)
}
func (b *Bot) DispatchInfoZRU(chatID int64, id string) error {
	return b.getKPI(chatID, "1")
}
func (b *Bot) DispatchInfoRU(chatID int64, id string) error {
	return b.getKPI(chatID, "9")
}
func (b *Bot) NpoKPI(chatID int64, id string) error {
	return b.getKPI(chatID, "2")
}
func (b *Bot) SmoothedKpiZRU(chatID int64, id string) error {
	return b.getKPI(chatID, "3")
}
func (b *Bot) SmoothedKpiRU(chatID int64, id string) error {
	return b.getKPI(chatID, "10")
}
func (b *Bot) OutPotokKPI(chatID int64, id string) error {
	return b.getKPI(chatID, "7")
}
func (b *Bot) TimelinessKpi(chatID int64, id string) error {
	return b.getKPI(chatID, "5")
}
func (b *Bot) LostKpi(chatID int64, id string) error {
	return b.getKPI(chatID, "8")
}
func (b *Bot) CancelKpi(chatID int64, id string) error {
	return b.getKPI(chatID, "4")
}
func (b *Bot) VpKpi(chatID int64, id string) error {
	return b.getKPI(chatID, "6")
}
