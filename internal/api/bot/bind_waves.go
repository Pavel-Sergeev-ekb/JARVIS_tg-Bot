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

type Waves struct {
	id   int
	name string
	bind string
	text string
}

func (b *Bot) getWavesTypeInfo(chatID int64, wavesKey string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
		return err
	}

	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Printf("Ошибка подключения к БД: %v", err)
		b.SendMessage(chatID, "Произошла ошибка при подключении к базе данных", tgbotapi.ModeMarkdown)
		return err
	}
	defer db.Close(context.Background())

	query := `
    SELECT 
        id,
        name,
        bind,
        text
    FROM waves
    WHERE id = $1
    `

	var waves Waves
	err = db.QueryRow(context.Background(), query, wavesKey).Scan(
		&waves.id,
		&waves.name,
		&waves.bind,
		&waves.text,
	)

	if errors.Is(err, sql.ErrNoRows) {
		msg := fmt.Sprintf("Операция с кодом %s не найдена", wavesKey)
		b.SendMessage(chatID, msg, tgbotapi.ModeMarkdown)
		return fmt.Errorf("операция не найдена в базе данных: %w", err)
	}

	if err != nil {
		log.Printf("Ошибка при выполнении запроса: %v", err)
		b.SendMessage(chatID, "Произошла ошибка при получении данных", tgbotapi.ModeMarkdown)
		return err
	}
	message := fmt.Sprintf(
		"**Название:** %s\n"+
			"**Описание:** %s\n",
		waves.name,
		waves.text,
	)
	b.SendMessage(chatID, message, tgbotapi.ModeMarkdown)
	return nil
}

func (b *Bot) PotConsInfo(chatID int64, wavesID string) error {
	return b.getWavesTypeInfo(chatID, "1")
}
func (b *Bot) SingleInfo(chatID int64, wavesID string) error {
	return b.getWavesTypeInfo(chatID, "3")
}
func (b *Bot) StationInfo(chatID int64, wavesID string) error {
	return b.getWavesTypeInfo(chatID, "2")
}
func (b *Bot) KGTInfo(chatID int64, wavesID string) error {
	return b.getWavesTypeInfo(chatID, "4")
}

func (b *Bot) getBindInfo(chatID int64, wavesKey string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
		return err
	}

	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Printf("Ошибка подключения к БД: %v", err)
		b.SendMessage(chatID, "Произошла ошибка при подключении к базе данных", tgbotapi.ModeMarkdown)
		return err
	}
	defer db.Close(context.Background())

	query := `
    SELECT 
        id,
        name,
        bind,
        text
    FROM waves
    WHERE id = $1
    `

	var waves Waves
	err = db.QueryRow(context.Background(), query, wavesKey).Scan(
		&waves.id,
		&waves.name,
		&waves.bind,
		&waves.text,
	)

	if errors.Is(err, sql.ErrNoRows) {
		msg := fmt.Sprintf("Операция с кодом %s не найдена", wavesKey)
		b.SendMessage(chatID, msg, tgbotapi.ModeMarkdown)
		return fmt.Errorf("операция не найдена в базе данных: %w", err)
	}

	if err != nil {
		log.Printf("Ошибка при выполнении запроса: %v", err)
		b.SendMessage(chatID, "Произошла ошибка при получении данных", tgbotapi.ModeMarkdown)
		return err
	}
	message := fmt.Sprintf(
		"**Название:** %s\n"+
			"**Описание:** %s\n",
		waves.name,
		waves.bind,
	)
	b.SendMessage(chatID, message, tgbotapi.ModeMarkdown)
	return nil
}
func (b *Bot) PotConsBindInfo(chatID int64, wavesID string) error {
	return b.getBindInfo(chatID, "1")
}
func (b *Bot) SingleBindInfo(chatID int64, wavesID string) error {
	return b.getBindInfo(chatID, "3")
}
func (b *Bot) StationBindInfo(chatID int64, wavesID string) error {
	return b.getBindInfo(chatID, "2")
}
func (b *Bot) KGTBindInfo(chatID int64, wavesID string) error {
	return b.getBindInfo(chatID, "4")
}
