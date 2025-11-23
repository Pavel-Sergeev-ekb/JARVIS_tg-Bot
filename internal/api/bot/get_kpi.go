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
		b.SendMessage(chatID, "Произошла ошибка при загрузке конфигурации")
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
		b.SendMessage(chatID, fmt.Sprintf("Показатель %s не найден", id))
		return fmt.Errorf("показатель не найден в базе данных: %w", err)
	}

	if err != nil {
		log.Printf("Ошибка при выполнении запроса: %v", err)
		b.SendMessage(chatID, "Произошла ошибка при получении данных")
		return err
	}

	message := formatKPIMessage(kpi, username)
	b.SendMessage(chatID, message)
	return nil
}

func formatKPIMessage(kpi KPI, username string) string {
	return fmt.Sprintf(
		"Название: %s\n"+
			"Показатель: %.2f\n"+
			"Вес: %d%%\n"+
			"Обновил: %s\n"+
			"Последнее обновление: %s\n",
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
