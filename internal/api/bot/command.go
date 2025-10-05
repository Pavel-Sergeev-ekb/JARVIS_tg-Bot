package bot

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	BotAPI *tgbotapi.BotAPI
}

func NewOneBot(botAPI *tgbotapi.BotAPI) *Bot {
	return &Bot{
		BotAPI: botAPI,
	}
}

func (b *Bot) HandleUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}
	chatID := update.Message.Chat.ID
	text := update.Message.Text

	b.SendGreetingKeyboard(chatID)

	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case "start":
			b.StartCommand(chatID, text)
		case "hourlyReport":
			b.HourlyReport(chatID, text)
		default:
			b.DefaultCommand(chatID, text)
		}
	} else {
		b.HandleTextMessage(chatID, text)
	}
}

func (b *Bot) SendGreetingKeyboard(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Выберите действие:")
	keyboard := NewMainKeyboard()
	msg.ReplyMarkup = keyboard

	_, err := b.BotAPI.Send(msg)
	if err != nil {
		log.Printf("Ошибка при отправке клавиатуры: %v", err)
	}
}

func (b *Bot) SendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := b.BotAPI.Send(msg)
	if err != nil {
		log.Printf("Ошибка при отправке сообщения: %v", err)
	}
}

func (b *Bot) StartCommand(chatID int64, t string) {
	text := "Привет, командир!\n" + "Я - J.A.R.V.I.S(Джарвис)...\nДа, когда-то давно я помогал Железному Человеку строить империю, теперь я на пенсии и хочу помогать бригадирам в ЯМ\nПриступим?"
	b.SendMessage(chatID, text)

}

func (b *Bot) DefaultCommand(chatID int64, t string) {
	text := "Ты что - то нажмякал, ничего не разобрать, напиши /help, чтобы увидеть доступные варианты"
	b.SendMessage(chatID, text)
}

func (b *Bot) HandleTextMessage(chatID int64, text string) {
	text = strings.ToLower(text)
	text = strings.TrimSpace(text)

	switch text {
	case "привет", "здарова", "хей", "хай", "хелло":
		b.SendMessage(chatID, "Привет, коллега, чем могу помочь?")

	case "как дела?", "как ты?":
		b.SendMessage(chatID, "О, все здорово! А у вас как дела?")

	default:
		b.SendMessage(chatID, "Пока что я только учусь, не понимаю о чем ты говоришь...")
	}
}

func (b *Bot) HandleCallback(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID

	if strings.HasPrefix(callback.Data, "status:") {
		statusID := strings.Split(callback.Data, ":")[1]
		b.CreatedExt(chatID, statusID)
		return
	}

	if strings.HasPrefix(callback.Data, "op_") {
		operationCode := strings.Split(callback.Data, "_")[1]
		b.ShowOperationInfo(chatID, operationCode)
		return
	}

	switch callback.Data {

	case "wms":
		// При нажатии на основную кнопку "Автостарт" переходим в меню
		msg := tgbotapi.NewMessage(chatID, "Выберите раздел: ")
		msg.ReplyMarkup = NewWMSMenu()
		b.BotAPI.Send(msg)

	case "hourlyReport":

		b.HourlyReport(chatID, "")

	case "tutorial":
		b.Tutorial(chatID, "")

	case "launch":
		msg := tgbotapi.NewMessage(chatID, "Выберите раздел: ")
		msg.ReplyMarkup = NewLaunchKeyboard()
		b.BotAPI.Send(msg)

	case "autopilot":
		msg := tgbotapi.NewMessage(chatID, "Выберите раздел: ")
		msg.ReplyMarkup = NewAutostartKeyboard()
		b.BotAPI.Send(msg)

	case "typeWaves":
		msg := tgbotapi.NewMessage(chatID, "Выберите раздел: ")
		msg.ReplyMarkup = NewWavesKeyboard()
		b.BotAPI.Send(msg)
	case "binding":
		msg := tgbotapi.NewMessage(chatID, "Выберите раздел: ")
		msg.ReplyMarkup = NewWavesKeyboard()
		b.BotAPI.Send(msg)
	case "hours":
		msg := tgbotapi.NewMessage(chatID, "Выберите раздел: ")
		msg.ReplyMarkup = NewHoursKeyboard()
		b.BotAPI.Send(msg)

	case "operations":
		msg := tgbotapi.NewMessage(chatID, "Выберите операцию:")
		msg.ReplyMarkup = NewWMSOperationsKeyboard()
		b.BotAPI.Send(msg)

	case "controlPoints":
		b.SendMessage(chatID, "в разработке")

	case "canceled":
		b.SendMessage(chatID, "в разработке")

	case "selection":
		b.SendMessage(chatID, "в разработке")

	case "autoSelection":
		b.SendMessage(chatID, "в разработке")

	case "door":
		b.SendMessage(chatID, "в разработке")

	case "shipped":
		b.SendMessage(chatID, "в разработке")

	case "packing":
		b.SendMessage(chatID, "в разработке")

	case "back_to_autostart":
		b.SendMessage(chatID, "Выберите раздел меню:")
		msg := tgbotapi.NewMessage(chatID, "")
		msg.ReplyMarkup = NewWMSMenu()
		b.BotAPI.Send(msg)

	case "orderStatuses":
		msg := tgbotapi.NewMessage(chatID, "Выберите статус заказа:")
		msg.ReplyMarkup = NewOrderStatusesKeyboard() // добавьте эту функцию
		b.BotAPI.Send(msg)

	case "back":
		msg := tgbotapi.NewMessage(chatID, "Привет, командир! Чем могу помочь?")
		msg.ReplyMarkup = NewMainKeyboard()
		b.BotAPI.Send(msg)

	default:
		b.SendMessage(chatID, "Неизвестная команда")
	}
	//нужно найти решение
	//	config := tgbotapi.CallbackConfig{
	//	CallbackQueryID: callback.ID,
	//	Text:            "✅ Команда выполнена",
	//	ShowAlert:       false,
	//}
	//_, err := b.BotAPI.AnswerCallbackQuery(config)
	//if err != nil {
	//	log.Printf("Ошибка при подтверждении callback: %v", err)
	//	}
}

func (b *Bot) HourlyReport(chatID int64, t string) {
	text := "Я уже учусь часовому отчету, пока что можешь покидать мне Excel файлы с точками контроля WMS"
	b.SendMessage(chatID, text)
}

func (b *Bot) Tutorial(chatID int64, t string) {
	text := "Тут будет руководство по работе бригадира исходящего потока"
	b.SendMessage(chatID, text)
}
