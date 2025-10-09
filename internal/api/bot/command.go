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

func (b *Bot) HandleCallbackKeyboard(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID

	switch callback.Data {

	case "wms":
		// При нажатии на основную кнопку "Автостарт" переходим в меню
		msg := tgbotapi.NewMessage(chatID, "Выберите раздел: ")
		msg.ReplyMarkup = NewWMSMenu()
		b.BotAPI.Send(msg)

	case "hourlyReport":

		b.HourlyReport(chatID, "")

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

	case "back_to_autostart":
		b.SendMessage(chatID, "Выберите раздел меню:")
		msg := tgbotapi.NewMessage(chatID, "")
		msg.ReplyMarkup = NewWMSMenu()
		b.BotAPI.Send(msg)

	case "back":
		msg := tgbotapi.NewMessage(chatID, "Командир! Чем могу помочь?")
		msg.ReplyMarkup = NewMainKeyboard()
		b.BotAPI.Send(msg)

	default:

	}

}
func (b *Bot) HourlyReport(chatID int64, t string) {
	text := "Я уже учусь часовому отчету, пока что можешь покидать мне Excel файлы с точками контроля WMS"
	b.SendMessage(chatID, text)
}

func (b *Bot) HandleCallbackOperation(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID

	switch callback.Data {

	case "operations":
		msg := tgbotapi.NewMessage(chatID, "Выбери операцию:")
		msg.ReplyMarkup = NewWMSOperationsKeyboard()
		b.BotAPI.Send(msg)

	case "controlPoints":
		b.ControlPoints(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Окей, босс...\nС этим разобрались\nО чем-то еще хочешь узнать?\n")
		msg.ReplyMarkup = NewWMSMenu()
		b.BotAPI.Send(msg)

	case "canceled":
		b.CanceledOrders(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Командир, двигаем дальше?")
		msg.ReplyMarkup = NewWMSMenu()
		b.BotAPI.Send(msg)

	case "selection":
		b.Selection20(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Командир, ну тут все просто!\nО чем-то еще хочешь узнать?")
		msg.ReplyMarkup = NewWMSMenu()
		b.BotAPI.Send(msg)

	case "autoSelection":
		b.AutoSelection20(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Да, это было просто\nО чем-то еще хочешь узнать?")
		msg.ReplyMarkup = NewWMSMenu()
		b.BotAPI.Send(msg)

	case "door":
		b.Door(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Да, командир, это база!\nО чем еще поговорим?")
		msg.ReplyMarkup = NewWMSMenu()
		b.BotAPI.Send(msg)

	case "shipped":
		b.Shipped(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Пу-пу-пу\n Погнали дальше?")
		msg.ReplyMarkup = NewWMSMenu()
		b.BotAPI.Send(msg)

	case "packing":
		b.Packing(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Подсказка: используй везде, где нужна посылочная этикетка\nЛадно, давай еще что-нибудь глянем!")
		msg.ReplyMarkup = NewWMSMenu()
		b.BotAPI.Send(msg)

	case "waveDetails":
		b.WaveDetails(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Тут можно целый курс выпустить как работать с этим инструментом, но мне кажется не много опыта и будет легко\nЛадненько давай еще что-нибудь изучим")
		msg.ReplyMarkup = NewWMSMenu()
		b.BotAPI.Send(msg)

	case "dispatcher":
		b.Dispatcher(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "В день на ДПД документами занимается диспетчер, а вот ночью 5POST сами\nКстати говоря, если нет ТС - то и документы не нужны.")
		msg.ReplyMarkup = NewWMSMenu()
		b.BotAPI.Send(msg)

	case "handStartInfo":
		b.HandStartInfo(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Ручником пользоваться нужно аккуратно - часто возникают проблемы, приоритетнее автостарт в 99% случаев\nХорошо, я готов к следующему вопросу")
		msg.ReplyMarkup = NewWMSMenu()
		b.BotAPI.Send(msg)

	case "transactionSKU":
		b.TransactionSKU(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Для чего это понятно\nКак пользоваться - расскажу в туториале\nА теперь давай вернемся к ознакомлению ")
		msg.ReplyMarkup = NewWMSMenu()
		b.BotAPI.Send(msg)

	case "transactionUIT":
		b.TransactionUIT(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Для чего это понятно\nКак пользоваться - расскажу в туториале\nА теперь давай вернемся к ознакомлению ")
		msg.ReplyMarkup = NewWMSMenu()
		b.BotAPI.Send(msg)

	case "balance":
		b.Balance(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Балансы очень полезны, есть несколько видов, нужно понимать какие когда использовать\nЧто-то еще рассказать?")
		msg.ReplyMarkup = NewWMSMenu()
		b.BotAPI.Send(msg)

	case "balanceUIT":
		b.BalanceUIT(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Балансы очень полезны, есть несколько видов, нужно понимать какие когда использовать\nЧто-то еще рассказать?")
		msg.ReplyMarkup = NewWMSMenu()
		b.BotAPI.Send(msg)

	case "checkOrders":
		b.CheckOrders(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Еще многое предстоит узнать, но не будем торопиться, я если что тут, можешь нажать кнопочку и я поведаю еще какие нибудь тайны")
		msg.ReplyMarkup = NewWMSMenu()
		b.BotAPI.Send(msg)

	case "controlPointsWMS":
		b.ControlPointsWMS(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Сколько раз обновил?\nМного, очень много...\nВкладка с точками открыта всегда, а кнопка обновления не отжимается не на секунду\nYep,yep - light weight baby - погнали дальше")
		msg.ReplyMarkup = NewWMSMenu()
		b.BotAPI.Send(msg)

	case "locking":
		b.Locking(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Много преимуществ сулит это обновление: подсвечивание в админке процессов на которых был сотрудник, его средний производ, где топчик, где среднячок, а куда лучше не ставить\nТакже сотрудник не сможет сам менять операцию, если вдруг работать на текущей ему не нравится или наскучило\nОкей, босс, пройдем-те дальше")
		msg.ReplyMarkup = NewWMSMenu()
		b.BotAPI.Send(msg)

	case "autoStart":
		b.AutoStart(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Есть отдельная менюшка, где будем детально смотреть, что за что отвечает\nОстались вопросы?")
		msg.ReplyMarkup = NewWMSMenu()
		b.BotAPI.Send(msg)

	case "action":
		b.Action(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Справедливость и отсутствие фрода - залог успеха\nДавай, командир, я только размялся, спроси еще что-нибудь")
		msg.ReplyMarkup = NewWMSMenu()
		b.BotAPI.Send(msg)

	case "skills":
		b.Skills(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Подробнее разберем этот вопрос в туториал\nЕще что-то хочешь узнать?")
		msg.ReplyMarkup = NewWMSMenu()
		b.BotAPI.Send(msg)

	case "AssignedWork":
		b.AssignedWork(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Подсказка: Назначенная работа повторяет 'Выпущено', но не 100%, определенные заказы уже собраны, но не сброшены в буфер\nЕще что-то хочешь узнать?")
		msg.ReplyMarkup = NewWMSMenu()
		b.BotAPI.Send(msg)

	case "hiring":
		b.Hiring(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Подробнее разберем этот вопрос в туториал\nЕще что-то хочешь узнать?")
		msg.ReplyMarkup = NewWMSMenu()
		b.BotAPI.Send(msg)

	case "backToWMS":
		msg := tgbotapi.NewMessage(chatID, "Выберите раздел:")
		msg.ReplyMarkup = NewWMSMenu()
		b.BotAPI.Send(msg)

	}

}

func (b *Bot) HandleCallbackOrders(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID

	switch callback.Data {
	case "createdExt":
		b.createdExt(chatID)
		msg := tgbotapi.NewMessage(chatID, "Босс, еще чем то могу помочь?")
		msg.ReplyMarkup = NewWMSMenu()
		b.BotAPI.Send(msg)

	case "orderStatuses":
		msg := tgbotapi.NewMessage(chatID, "Выберите статус заказа:")
		msg.ReplyMarkup = NewOrderStatusesKeyboard() // добавьте эту функцию
		b.BotAPI.Send(msg)

	case "released":
		b.released(chatID)
		msg := tgbotapi.NewMessage(chatID, "Босс, еще чем то могу помочь?")
		msg.ReplyMarkup = NewWMSMenu()
		b.BotAPI.Send(msg)

	case "sorted":
		b.sorted(chatID)
		msg := tgbotapi.NewMessage(chatID, "Босс, еще чем то могу помочь?")
		msg.ReplyMarkup = NewWMSMenu()
		b.BotAPI.Send(msg)

	case "packed":
		b.packed(chatID)
		msg := tgbotapi.NewMessage(chatID, "Босс, еще чем то могу помочь?")
		msg.ReplyMarkup = NewWMSMenu()
		b.BotAPI.Send(msg)

	case "sortedSD":
		b.sortedSD(chatID)
		msg := tgbotapi.NewMessage(chatID, "Босс, еще чем то могу помочь?")
		msg.ReplyMarkup = NewWMSMenu()
		b.BotAPI.Send(msg)

	case "shipment":
		b.shipment(chatID)
		msg := tgbotapi.NewMessage(chatID, "Босс, еще чем то могу помочь?")
		msg.ReplyMarkup = NewWMSMenu()
		b.BotAPI.Send(msg)

	case "KIZ":
		b.KIZ(chatID)
		msg := tgbotapi.NewMessage(chatID, "Босс, еще чем то могу помочь?")
		msg.ReplyMarkup = NewWMSMenu()
		b.BotAPI.Send(msg)

	case "selectionCompleted":
		b.selectionCompleted(chatID)
		msg := tgbotapi.NewMessage(chatID, "Босс, еще чем то могу помочь?")
		msg.ReplyMarkup = NewWMSMenu()
		b.BotAPI.Send(msg)

	case "unknown":
		b.Unknown(chatID)
		msg := tgbotapi.NewMessage(chatID, "Босс, еще чем то могу помочь?")
		msg.ReplyMarkup = NewWMSMenu()
		b.BotAPI.Send(msg)
	}

}

func (b *Bot) HandleCallbackTutorial(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID

	switch callback.Data {

	case "tutorial":
		msg := tgbotapi.NewMessage(chatID, "ТУТОРИАЛ ")
		msg.ReplyMarkup = NewTutorialKeyboard()
		b.BotAPI.Send(msg)
	case "primer":
		msg := tgbotapi.NewMessage(chatID, "БУКВАРЬ ")
		msg.ReplyMarkup = NewPrimerKeyboard()
		b.BotAPI.Send(msg)

	case "backToTutorial":
		msg := tgbotapi.NewMessage(chatID, "ТУТОРИАЛ")
		msg.ReplyMarkup = NewTutorialKeyboard()
		b.BotAPI.Send(msg)

	case "shippedIns":
		msg := tgbotapi.NewMessage(chatID, "ОТГРУЗКА")
		msg.ReplyMarkup = NewShippedInsKeyboard()
		b.BotAPI.Send(msg)

	case "fact":
		msg := tgbotapi.NewMessage(chatID, "ФИЗИЧЕСКАЯ ОТГРУЗКА")
		msg.ReplyMarkup = NewFactInsKeyboard()
		b.BotAPI.Send(msg)

	case "transactionIns":
		msg := tgbotapi.NewMessage(chatID, "ТРАНЗАКЦИИ")
		msg.ReplyMarkup = NewTransactionInsMenu()
		b.BotAPI.Send(msg)

	case "seal":
		msg := tgbotapi.NewMessage(chatID, "ПЕЧАТЬ ЭТИКЕТОК")
		msg.ReplyMarkup = NewSealKeyboard()
		b.BotAPI.Send(msg)

	case "backToShip":
		msg := tgbotapi.NewMessage(chatID, "ОТГРУЗКА")
		msg.ReplyMarkup = NewShippedInsKeyboard()
		b.BotAPI.Send(msg)

	case "Replenishments":
		msg := tgbotapi.NewMessage(chatID, "ПОПОЛНЕНИЯ")
		msg.ReplyMarkup = NewReplenishKeyboard()
		b.BotAPI.Send(msg)

	}
}

func (b *Bot) HandleCallbackKPI(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID

	switch callback.Data {

	case "kpi":
		msg := tgbotapi.NewMessage(chatID, "Выбери должность")
		msg.ReplyMarkup = NewKPIKeyboard()
		b.BotAPI.Send(msg)

	case "zru":
		msg := tgbotapi.NewMessage(chatID, "KPI Заместителя Руководителя Участков")
		msg.ReplyMarkup = NewZruKpiKeyboard()
		b.BotAPI.Send(msg)

	case "brigadir":
		msg := tgbotapi.NewMessage(chatID, "KPI Бригадира/Руководителя Участков")
		msg.ReplyMarkup = NewBrigadirKpiKeyboard()
		b.BotAPI.Send(msg)

	}
}
