package main

import (
	"context"
	"log"

	"github.com/Pavel-Sergeev-ekb/JARVIS_tg-Bot/internal/api/bot"
	"github.com/Pavel-Sergeev-ekb/JARVIS_tg-Bot/internal/api/database"
	"github.com/Pavel-Sergeev-ekb/JARVIS_tg-Bot/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Printf("Ошибка подключения к БД: %v", err)
		return
	}
	defer db.Close(context.Background())

	botApi, err := tgbotapi.NewBotAPI(cfg.BotToken)

	if err != nil {
		log.Panic(err)
	}

	botApi.Debug = true

	log.Printf("Авторизован на аккаунте %s", botApi.Self.UserName)

	// инициализация

	botInstance := bot.NewOneBot(botApi)

	//запрос на получение обновлений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// получение канала обновлений

	updates := botApi.GetUpdatesChan(u)

	// обработка полученных обновлений

	for update := range updates {
		switch {
		case update.Message != nil:
			botInstance.HandleUpdate(update)
		case update.CallbackQuery != nil:
			botInstance.HandleCallbackKeyboard(update.CallbackQuery)
			botInstance.HandleCallbackOperation(update.CallbackQuery)
			botInstance.HandleCallbackOrders(update.CallbackQuery)
			botInstance.HandleCallbackTutorial(update.CallbackQuery)
			botInstance.HandleCallbackKPI(update.CallbackQuery)
		}
	}
}
