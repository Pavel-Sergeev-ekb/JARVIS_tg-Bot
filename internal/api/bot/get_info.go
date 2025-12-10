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

type info struct {
	id          int
	name        string
	description string
}

func (b *Bot) getBigInfo(chatID int64, infoKey string) error {
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
        name,
        description
    FROM public.data_desc
    WHERE id = $1
    `

	var info info
	err = db.QueryRow(context.Background(), query, infoKey).Scan(
		&info.id,
		&info.name,
		&info.description,
	)

	if errors.Is(err, sql.ErrNoRows) {
		msg := fmt.Sprintf("Операция с кодом %s не найдена", infoKey)
		b.SendMessage(chatID, msg, tgbotapi.ModeHTML)
		return fmt.Errorf("операция не найдена в базе данных: %w", err)
	}

	if err != nil {
		log.Printf("Ошибка при выполнении запроса: %v", err)
		b.SendMessage(chatID, "Произошла ошибка при получении данных", tgbotapi.ModeHTML)
		return err
	}
	message := fmt.Sprintf(
		"<b>Название:</b> <b>%s</b>\n"+
			"<b>Описание:</b> <b>%s</b>\n",
		info.name,
		info.description,
	)
	b.SendMessage(chatID, message, tgbotapi.ModeHTML)
	return nil
}

func (b *Bot) TimeWithInfo(chatID int64, infoKey string) error {
	return b.getBigInfo(chatID, "1")
}

func (b *Bot) TimeDoInfo(chatID int64, infoKey string) error {
	return b.getBigInfo(chatID, "2")
}

func (b *Bot) HandLetGoInfo(chatID int64, infoKey string) error {
	return b.getBigInfo(chatID, "3")
}
func (b *Bot) KPIInfo(chatID int64, infoKey string) error {
	return b.getBigInfo(chatID, "4")
}
func (b *Bot) SlaInfo(chatID int64, infoKey string) error {
	return b.getBigInfo(chatID, "5")
}
func (b *Bot) LOSTInfo(chatID int64, infoKey string) error {
	return b.getBigInfo(chatID, "6")
}
func (b *Bot) CapInfo(chatID int64, infoKey string) error {
	return b.getBigInfo(chatID, "7")
}
func (b *Bot) DayOffInfo(chatID int64, infoKey string) error {
	return b.getBigInfo(chatID, "8")
}
func (b *Bot) AlarmInfo(chatID int64, infoKey string) error {
	return b.getBigInfo(chatID, "9")
}
func (b *Bot) NPOInfo(chatID int64, infoKey string) error {
	return b.getBigInfo(chatID, "10")
}
func (b *Bot) SmoothedInfo(chatID int64, infoKey string) error {
	return b.getBigInfo(chatID, "12")
}
func (b *Bot) ShipmentProInfo(chatID int64, infoKey string) error {
	return b.getBigInfo(chatID, "13")
}
func (b *Bot) OpenProInfo(chatID int64, infoKey string) error {
	return b.getBigInfo(chatID, "14")
}
func (b *Bot) FTEInfo(chatID int64, infoKey string) error {
	return b.getBigInfo(chatID, "15")
}
func (b *Bot) GermesInfo(chatID int64, infoKey string) error {
	return b.getBigInfo(chatID, "16")
}
func (b *Bot) TicketInfo(chatID int64, infoKey string) error {
	return b.getBigInfo(chatID, "17")
}
func (b *Bot) WithdrawalInfo(chatID int64, infoKey string) error {
	return b.getBigInfo(chatID, "18")
}
func (b *Bot) ScanInfo(chatID int64, infoKey string) error {
	return b.getBigInfo(chatID, "19")
}
func (b *Bot) CutOffInfo(chatID int64, infoKey string) error {
	return b.getBigInfo(chatID, "20")
}
func (b *Bot) EkbShipInfo(chatID int64, infoKey string) error {
	return b.getBigInfo(chatID, "21")
}
func (b *Bot) TCShipInfo(chatID int64, infoKey string) error {
	return b.getBigInfo(chatID, "22")
}
func (b *Bot) TransitShipInfo(chatID int64, infoKey string) error {
	return b.getBigInfo(chatID, "23")
}
func (b *Bot) SystemShipInfo(chatID int64, infoKey string) error {
	return b.getBigInfo(chatID, "24")
}
func (b *Bot) SkillsProInfo(chatID int64, infoKey string) error {
	return b.getBigInfo(chatID, "25")
}
func (b *Bot) LockingINFO(chatID int64, infoKey string) error {
	return b.getBigInfo(chatID, "26")
}
func (b *Bot) TransactionYP(chatID int64, infoKey string) error {
	return b.getBigInfo(chatID, "31")
}
func (b *Bot) TransactionUitInfo(chatID int64, infoKey string) error {
	return b.getBigInfo(chatID, "27")
}
func (b *Bot) TransactionROV(chatID int64, infoKey string) error {
	return b.getBigInfo(chatID, "28")
}
func (b *Bot) TransactionNZN(chatID int64, infoKey string) error {
	return b.getBigInfo(chatID, "29")
}
func (b *Bot) TransactionNuance(chatID int64, infoKey string) error {
	return b.getBigInfo(chatID, "30")
}
func (b *Bot) SealDRP(chatID int64, infoKey string) error {
	return b.getBigInfo(chatID, "32")
}
func (b *Bot) SealTRP(chatID int64, infoKey string) error {
	return b.getBigInfo(chatID, "33")
}
func (b *Bot) SealNZN(chatID int64, infoKey string) error {
	return b.getBigInfo(chatID, "34")
}
func (b *Bot) SealUIT(chatID int64, infoKey string) error {
	return b.getBigInfo(chatID, "36")
}
func (b *Bot) SealP(chatID int64, infoKey string) error {
	return b.getBigInfo(chatID, "37")
}
func (b *Bot) ReplenishClientOrder(chatID int64, infoKey string) error {
	return b.getBigInfo(chatID, "38")
}
func (b *Bot) ReplenishWithdrawal(chatID int64, infoKey string) error {
	return b.getBigInfo(chatID, "39")
}
func (b *Bot) ShortBriefing(chatID int64, infoKey string) error {
	return b.getBigInfo(chatID, "40")
}
func (b *Bot) LongBriefing(chatID int64, infoKey string) error {
	return b.getBigInfo(chatID, "42")
}
func (b *Bot) LongBriefing2(chatID int64, infoKey string) error {
	return b.getBigInfo(chatID, "43")
}
func (b *Bot) FormsOT(chatID int64, infoKey string) error {
	return b.getBigInfo(chatID, "46")
}
func (b *Bot) DashInfo(chatID int64, infoKey string) error {
	return b.getBigInfo(chatID, "45")
}

func (b *Bot) LinksInfo(chatID int64, infoKey string) error {
	return b.getBigInfo(chatID, "47")
}
