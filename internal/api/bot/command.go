package bot

import (
	"fmt"
	"log"
	"strings"

	"github.com/Pavel-Sergeev-ekb/JARVIS_tg-Bot/internal/api/database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	BotAPI               *tgbotapi.BotAPI
	waitingForStandInput map[int64]string
	waitingForKPIInput   map[int64]string
	waitingForFirstFile  map[int64]bool
	waitingForSecondFile map[int64]bool
	tempFilePaths        map[int64]map[int]string
}

const TargetChatID int64 = -5016795698

func NewOneBot(botAPI *tgbotapi.BotAPI) *Bot {
	return &Bot{
		BotAPI:               botAPI,
		waitingForKPIInput:   make(map[int64]string),
		waitingForStandInput: make(map[int64]string),
		waitingForFirstFile:  make(map[int64]bool),
		waitingForSecondFile: make(map[int64]bool),
		tempFilePaths:        make(map[int64]map[int]string),
	}
}

func (b *Bot) HandleUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID

	if update.Message.Document != nil {
		b.HandleFileUpload(chatID, update.Message.Document)
		return
	}

	text := update.Message.Text
	kpiID, ok := b.waitingForKPIInput[chatID]

	if ok {
		delete(b.waitingForKPIInput, chatID)
		b.handleUpdateInput(chatID, update.Message.Text, kpiID)
		return // Завершаем метод после обработки ожидаемого ввода
	}

	standartID, ok := b.waitingForStandInput[chatID]

	if ok {
		delete(b.waitingForStandInput, chatID)
		b.handleUpdateInputStand(chatID, update.Message.Text, standartID)
		return
	}

	if update.Message.IsCommand() {
		b.handleCommand(chatID, update.Message)
		return // Завершаем метод после обработки команды
	}

	b.SendGreetingKeyboard(chatID)
	b.HandleTextMessage(chatID, text)

}

func (b *Bot) handleCommand(chatID int64, msg *tgbotapi.Message) {
	switch msg.Command() {
	case "start":
		b.StartCommand(msg)
	case "hourlyReport":
		b.waitingForFirstFile[chatID] = true   // Ждём первый файл
		b.waitingForSecondFile[chatID] = false // Сбрасываем ожидание второго
		b.SendMessage(chatID, "Отправьте первый Excel‑файл Точек Контроля")
	default:
		b.DefaultCommand(chatID, msg.Text)
	}
}

func (b *Bot) HandleFile(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	doc := update.Message.Document

	if b.waitingForFirstFile[chatID] {
		// Инициализируем внутреннюю мапу, если её нет
		if b.tempFilePaths[chatID] == nil {
			b.tempFilePaths[chatID] = make(map[int]string)
		}

		err := b.downloadTelegramFile(chatID, doc, 0)
		if err != nil {
			b.SendMessage(chatID, "Ошибка сохранения первого файла")
			return
		}

		b.waitingForFirstFile[chatID] = false
		b.waitingForSecondFile[chatID] = true
		b.SendMessage(chatID, "Отправьте второй Excel‑файл Точек Контроля")

	} else if b.waitingForSecondFile[chatID] {
		// Инициализируем внутреннюю мапу, если её нет
		if b.tempFilePaths[chatID] == nil {
			b.tempFilePaths[chatID] = make(map[int]string)
		}

		err := b.downloadTelegramFile(chatID, doc, 1)
		if err != nil {
			b.SendMessage(chatID, "Ошибка сохранения второго файла")
			return
		}

		err = b.sendReport(TargetChatID)
		if err != nil {
			b.SendMessage(chatID, fmt.Sprintf("Ошибка отправки отчёта: %v", err))
		} else {
			b.SendMessage(chatID, "Отчёт сформирован и отправлен в общий чат!")
		}

		b.waitingForFirstFile[chatID] = false
		b.waitingForSecondFile[chatID] = false
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

func (b *Bot) SendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := b.BotAPI.Send(msg)
	if err != nil {
		log.Printf("Ошибка при отправке сообщения: %v", err)
	}
	return err
}

func (b *Bot) StartCommand(msg *tgbotapi.Message) {
	// Формируем приветственное сообщение
	text := "Привет, командир!\n" +
		"Я —  бот J.A.R.V.I.S (Джарвис)...\n" +
		"Я создан для помощи в обучении и адаптации, помогу сформировать часовой отчет, расскажу о текущих нормативах и ключевых показателях и еще очень и очень много всего полезного \n" +
		"Приступим?"

	// Отправляем приветственное сообщение
	b.SendMessage(msg.Chat.ID, text)

	// Сохраняем пользователя в БД
	err := database.SaveUser(msg.Chat.ID, msg.From.UserName)
	if err != nil {
		// Логируем ошибку (можно также отправить уведомление пользователю)
		log.Printf("Ошибка сохранения пользователя chat_id=%d, user_name=%s: %v",
			msg.Chat.ID, msg.From.UserName, err)

		// Опционально: уведомляем пользователя о проблеме
		b.SendMessage(msg.Chat.ID, "Произошла ошибка при сохранении ваших данных. Попробуйте позже.")
	}
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
		msg := tgbotapi.NewMessage(chatID, "Выбери раздел: ")
		msg.ReplyMarkup = NewWMSMenu()
		b.BotAPI.Send(msg)

	case "launch":
		msg := tgbotapi.NewMessage(chatID, "Выбери раздел: ")
		msg.ReplyMarkup = NewLaunchKeyboard()
		b.BotAPI.Send(msg)

	case "autopilot":
		msg := tgbotapi.NewMessage(chatID, "Выбери раздел: ")
		msg.ReplyMarkup = NewAutostartKeyboard()
		b.BotAPI.Send(msg)

	case "typeWaves":
		msg := tgbotapi.NewMessage(chatID, "Выбери раздел: ")
		msg.ReplyMarkup = NewWavesKeyboard()
		b.BotAPI.Send(msg)
	case "binding":
		msg := tgbotapi.NewMessage(chatID, "Выбери раздел: ")
		msg.ReplyMarkup = NewBindKeyboard()
		b.BotAPI.Send(msg)
	case "hours":
		msg := tgbotapi.NewMessage(chatID, "Выбери раздел: ")
		msg.ReplyMarkup = NewHoursKeyboard()
		b.BotAPI.Send(msg)

	case "back_to_autostart":
		b.SendMessage(chatID, "Выбери раздел меню:")
		msg := tgbotapi.NewMessage(chatID, "")
		msg.ReplyMarkup = NewWMSMenu()
		b.BotAPI.Send(msg)

	case "hourlyReport":
		b.waitingForFirstFile[chatID] = true   // Ждём первый файл
		b.waitingForSecondFile[chatID] = false // Сбрасываем ожидание второго
		b.SendMessage(chatID, "Отправьте первый Excel‑файл Точек Контроля")

	case "back":
		msg := tgbotapi.NewMessage(chatID, "Командир! Чем могу помочь?")
		msg.ReplyMarkup = NewMainKeyboard()
		b.BotAPI.Send(msg)

	case "potCons1":
		b.PotConsInfo(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Выбери раздел")
		msg.ReplyMarkup = NewWavesKeyboard()
		b.BotAPI.Send(msg)

	case "potCons2":
		b.PotConsBindInfo(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Выбери раздел")
		msg.ReplyMarkup = NewBindKeyboard()
		b.BotAPI.Send(msg)

	case "kgt1":
		b.KGTInfo(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Выбери раздел")
		msg.ReplyMarkup = NewWavesKeyboard()
		b.BotAPI.Send(msg)

	case "kgt2":
		b.KGTBindInfo(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Выбери раздел")
		msg.ReplyMarkup = NewBindKeyboard()
		b.BotAPI.Send(msg)

	case "single1":
		b.SingleInfo(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Выбери раздел")
		msg.ReplyMarkup = NewWavesKeyboard()
		b.BotAPI.Send(msg)

	case "single2":
		b.SingleBindInfo(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Выбери раздел")
		msg.ReplyMarkup = NewBindKeyboard()
		b.BotAPI.Send(msg)

	case "stage1":
		b.StationInfo(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Выбери раздел")
		msg.ReplyMarkup = NewWavesKeyboard()
		b.BotAPI.Send(msg)

	case "stage2":
		b.StationBindInfo(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Выбери раздел")
		msg.ReplyMarkup = NewBindKeyboard()
		b.BotAPI.Send(msg)

	case "backToAuto":
		msg := tgbotapi.NewMessage(chatID, "Выбери раздел")
		msg.ReplyMarkup = NewAutostartKeyboard()
		b.BotAPI.Send(msg)

	case "timeBefore":
		b.TimeDoInfo(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Выбери раздел")
		msg.ReplyMarkup = NewAutostartKeyboard()
		b.BotAPI.Send(msg)

	case "timeWith":
		b.TimeWithInfo(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Выбери раздел")
		msg.ReplyMarkup = NewAutostartKeyboard()
		b.BotAPI.Send(msg)

	case "handStart":
		b.HandLetGoInfo(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Выбери раздел")
		msg.ReplyMarkup = NewLaunchKeyboard()
		b.BotAPI.Send(msg)

	default:

	}
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
		msg := tgbotapi.NewMessage(chatID, "Командир, еще ячейки упаковки часто называют буфер консолидации, поэтому основная мысль - доехал ли заказ до буфера упаковки?")
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

	case "Unknown":
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

	case "kpi":
		b.KPIInfo(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Все, что человеческий разум способен понять и во что он способен поверить, достижимо. — Наполеон Хилл.")
		msg.ReplyMarkup = NewPrimerKeyboard()
		b.BotAPI.Send(msg)

	case "sla":
		b.SlaInfo(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Сложнее всего начать действовать, все остальное зависит только от упорства. — Амелия Эрхарт.")
		msg.ReplyMarkup = NewPrimerKeyboard()
		b.BotAPI.Send(msg)

	case "lost":
		b.LOSTInfo(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Вдохновение приходит только во время работы. — Габриэль Гарсиа Маркес.")
		msg.ReplyMarkup = NewPrimerKeyboard()
		b.BotAPI.Send(msg)

	case "capasity":
		b.CapInfo(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Если проблему можно решить, не стоит о ней беспокоиться. — Далай-лама XIV.")
		msg.ReplyMarkup = NewPrimerKeyboard()
		b.BotAPI.Send(msg)

	case "dayoff":
		b.DayOffInfo(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Постановка целей — это первый шаг к превращению невидимого в видимое. — Тони Роббинс.")
		msg.ReplyMarkup = NewPrimerKeyboard()
		b.BotAPI.Send(msg)

	case "dashboards":
		b.DashInfo(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Для того, чтобы преуспеть, мы первым делом должны верить, что мы можем. — Никос Казантзакис.")
		msg.ReplyMarkup = NewTutorialKeyboard()
		b.BotAPI.Send(msg)

	case "alarm":
		b.AlarmInfo(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Если что-то вообще стоит пробовать, это стоит попробовать не меньше 10 раз. — Артур Гордон Линклеттер.")
		msg.ReplyMarkup = NewPrimerKeyboard()
		b.BotAPI.Send(msg)

	case "NpoInfo":
		b.NPOInfo(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Дойдя до конца, люди смеются над страхами, мучившими их вначале. — Пауло Коэльо.")
		msg.ReplyMarkup = NewPrimerKeyboard()
		b.BotAPI.Send(msg)

	case "smoothed":
		b.SmoothedInfo(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Постарайтесь получить то, что любите, иначе придется полюбить то, что получили. — Джордж Бернард Шоу.")
		msg.ReplyMarkup = NewPrimerKeyboard()
		b.BotAPI.Send(msg)

	case "shipmentPro":
		b.ShipmentProInfo(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Живите в соответствии со своим воображением, а не своим прошлым. — Стивен Кови.")
		msg.ReplyMarkup = NewPrimerKeyboard()
		b.BotAPI.Send(msg)

	case "OpenPro":
		b.OpenProInfo(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Лидеры не рождаются и не делаются кем-либо — они делают себя сами. — Стивен Кови.")
		msg.ReplyMarkup = NewPrimerKeyboard()
		b.BotAPI.Send(msg)

	case "fte":
		b.FTEInfo(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Единственным пределом наших завтрашних свершений станут наши сегодняшние сомнения. — Франклин Д. Рузвельт.")
		msg.ReplyMarkup = NewPrimerKeyboard()
		b.BotAPI.Send(msg)

	case "germes":
		b.GermesInfo(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Храбрость не всегда кричит. Иногда храбрость говорит тихим голосом в конце дня: Завтра я попытаюсь еще раз. — Мэри Энн Редмачер.")
		msg.ReplyMarkup = NewPrimerKeyboard()
		b.BotAPI.Send(msg)

	case "ticket":
		b.TicketInfo(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Каждый день не может быть хорошим, но есть что-то хорошее в каждом дне. — Элис Морз Эрл")
		msg.ReplyMarkup = NewPrimerKeyboard()
		b.BotAPI.Send(msg)

	case "withdrawal":
		b.WithdrawalInfo(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Постановка целей является первым шагом на пути превращения мечты в реальность. — Энтони Роббинс.")
		msg.ReplyMarkup = NewPrimerKeyboard()
		b.BotAPI.Send(msg)

	case "infoscan":
		b.ScanInfo(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Мы должны использовать время мудро и помнить: правое дело можно начать в любую минуту. — Нельсон Мандела.")
		msg.ReplyMarkup = NewPrimerKeyboard()
		b.BotAPI.Send(msg)

	case "catoff":
		b.CutOffInfo(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Если Вы работаете над поставленными целями, то эти цели будут работать на вас. — Джим Рон.")
		msg.ReplyMarkup = NewPrimerKeyboard()
		b.BotAPI.Send(msg)

	case "spichRichtrack":
		msg := tgbotapi.NewMessage(chatID, "ФОРМА ИНСТРУКТАЖА")
		msg.ReplyMarkup = NewVashKeyboard()
		b.BotAPI.Send(msg)

	case "shortBriefing":
		b.ShortBriefing(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Величайшая польза, которую можно извлечь из жизни, — потратить жизнь на дело, которое переживёт нас. — Уильям Джеймс.")
		msg.ReplyMarkup = NewVashKeyboard()
		b.BotAPI.Send(msg)

	case "longBriefing":
		b.LongBriefing(chatID, "")
		b.LongBriefing2(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Только я могу изменить свою жизнь. Никто не может сделать это за меня. — Кэрол Бернетт.")
		msg.ReplyMarkup = NewVashKeyboard()
		b.BotAPI.Send(msg)

	case "formsOT":
		b.FormsOT(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Успех — это сумма мелких усилий, повторяющихся день за днем. — Роберт Кольер.")
		msg.ReplyMarkup = NewTutorialKeyboard()
		b.BotAPI.Send(msg)

	case "ekb":
		b.EkbShipInfo(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Когда необходимо сделать выбор, а вы его не делаете, — это тоже выбор. — Уильям Джеймс")
		msg.ReplyMarkup = NewFactInsKeyboard()
		b.BotAPI.Send(msg)

	case "transport":
		b.TCShipInfo(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Стремитесь не к успеху, а к ценностям, которые он дает. — Альберт Эйнштейн.")
		msg.ReplyMarkup = NewFactInsKeyboard()
		b.BotAPI.Send(msg)

	case "tranzit":
		b.TransitShipInfo(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Неважно, как медленно вы идете, до тех пор, пока вы не остановитесь. — Конфуций.")
		msg.ReplyMarkup = NewFactInsKeyboard()
		b.BotAPI.Send(msg)

	case "system":
		b.SystemShipInfo(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Верный способ начать что-то — бросить говорить и делать. — Уолт Дисней.")
		msg.ReplyMarkup = NewShippedInsKeyboard()
		b.BotAPI.Send(msg)

	case "skillsPro":
		b.SkillsProInfo(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Будьте собой. Все остальные роли уже заняты. — Оскар Уайльд.")
		msg.ReplyMarkup = NewTutorialKeyboard()
		b.BotAPI.Send(msg)

	case "lockingINFO":
		b.LockingINFO(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Секрет перемен состоит в том, чтобы сосредоточиться на создании нового, а не на борьбе со старым. — Сократ.")
		msg.ReplyMarkup = NewTutorialKeyboard()
		b.BotAPI.Send(msg)

	case "YP":
		b.TransactionYP(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Существует только один путь к счастью — перестать беспокоиться о вещах, которые не подвластны нашей воле. — Эпиктет.")
		msg.ReplyMarkup = NewTransactionInsMenu()
		b.BotAPI.Send(msg)

	case "UIT":
		b.TransactionUitInfo(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Действие даже самого крохотного существа приводит к изменениям во всей вселенной. — Никола Тесла.")
		msg.ReplyMarkup = NewTransactionInsMenu()
		b.BotAPI.Send(msg)

	case "ROV":
		b.TransactionROV(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Выживает не самый сильный из видов и не самый умный, а тот, кто лучше других реагирует на изменения. — Леон Меггинсон.")
		msg.ReplyMarkup = NewTransactionInsMenu()
		b.BotAPI.Send(msg)

	case "Tara":
		b.TransactionNZN(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Не будучи в силах свершить великое, свершай малое великим способом. — Наполеон Хилл.")
		msg.ReplyMarkup = NewTransactionInsMenu()
		b.BotAPI.Send(msg)

	case "nuance":
		b.TransactionNuance(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Грамм собственного опыта стоит дороже тонны чужих наставлений. — Махатма Ганди.")
		msg.ReplyMarkup = NewTransactionInsMenu()
		b.BotAPI.Send(msg)

	case "DRP":
		b.SealDRP(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Вы не сможете пересечь море, просто стоя и вглядываясь в воду. Не тратьте время на напрасные желания. — Рабиндранат Тагор.")
		msg.ReplyMarkup = NewSealKeyboard()
		b.BotAPI.Send(msg)

	case "TRP":
		b.SealTRP(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Не стоит быть тенью в чужом фонаре, когда можно сиять своим собственным светом. — Алексей Христинин.")
		msg.ReplyMarkup = NewSealKeyboard()
		b.BotAPI.Send(msg)

	case "uitSeal":
		b.SealUIT(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Мы должны использовать время, как инструмент, а не как диван. — Джон Фицджеральд Кеннеди.")
		msg.ReplyMarkup = NewSealKeyboard()
		b.BotAPI.Send(msg)

	case "nzn":
		b.SealNZN(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Большую работу поручают тому, кто показывает способность перерасти маленькую. — Теодор Рузвельт.")
		msg.ReplyMarkup = NewSealKeyboard()
		b.BotAPI.Send(msg)

	case "label":
		b.SealP(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Не отступай перед трудностями. Смотри им прямо в лицо. Смотри, пока не одолеешь их. — Чарльз Диккенс.")
		msg.ReplyMarkup = NewSealKeyboard()
		b.BotAPI.Send(msg)

	case "clientOrder":
		b.ReplenishClientOrder(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Чем больше трудностей мы преодолеваем, тем легче принимаем свое будущее. — Майкл Хайятт.")
		msg.ReplyMarkup = NewReplenishKeyboard()
		b.BotAPI.Send(msg)

	case "withdrawalOrder":
		b.ReplenishWithdrawal(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Вы не должны быть великими, чтобы начать, но вы должны начать, чтобы быть великими. — Зиг Зиглар.")
		msg.ReplyMarkup = NewReplenishKeyboard()
		b.BotAPI.Send(msg)

	}
}

func (b *Bot) HandleCallbackKPI(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID

	switch callback.Data {

	case "KPI":
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

	case "backToKPI":
		msg := tgbotapi.NewMessage(chatID, "Выбери должность")
		msg.ReplyMarkup = NewKPIKeyboard()
		b.BotAPI.Send(msg)

	case "dispatchZRU":
		b.DispatchInfoZRU(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Выбери должность")
		msg.ReplyMarkup = NewKPIKeyboard()
		b.BotAPI.Send(msg)

	case "dispatchRU":
		b.DispatchInfoRU(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Выбери должность")
		msg.ReplyMarkup = NewKPIKeyboard()
		b.BotAPI.Send(msg)

	case "refresh1":
		msg := tgbotapi.NewMessage(chatID, "Выбери показатель для обновления")
		msg.ReplyMarkup = NewRefreshKpiZRUKeyboard()
		b.BotAPI.Send(msg)

	case "refresh":
		msg := tgbotapi.NewMessage(chatID, "Выбери показатель для обновления")
		msg.ReplyMarkup = NewRefreshKpiRUKeyboard()
		b.BotAPI.Send(msg)

	case "dispatchUpZRU":
		b.HandleUpdateButton(chatID, "% Засылов(ЗРУ)")

	case "dispatchUpRU":
		b.HandleUpdateButton(chatID, "% Засылов(РУ)")

	case "vpKPI":
		b.VpKpi(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Выбери должность")
		msg.ReplyMarkup = NewKPIKeyboard()
		b.BotAPI.Send(msg)

	case "vpKPIUP":
		b.HandleUpdateButton(chatID, "Производ возвратов")

	case "cancelkpiUP":
		b.HandleUpdateButton(chatID, "% Отмен")

	case "lostKPIUp":
		b.HandleUpdateButton(chatID, "% Потерь склад")

	case "timelinesskpiUp":
		b.HandleUpdateButton(chatID, "Своевременность отгрузки")

	case "potokKPIUp":
		b.HandleUpdateButton(chatID, "Производ исхода")

	case "smoothedkpiUpZRU":
		b.HandleUpdateButton(chatID, "Сглаженный производ(ЗРУ)")

	case "smoothedkpiUpRU":
		b.HandleUpdateButton(chatID, "Сглаженный производ(РУ)")

	case "npoKPIUp":
		b.HandleUpdateButton(chatID, "Доля НПО")

	case "lostKPI":
		b.LostKpi(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Выбери должность")
		msg.ReplyMarkup = NewKPIKeyboard()
		b.BotAPI.Send(msg)

	case "timelinesskpi":
		b.TimelinessKpi(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Выбери должность")
		msg.ReplyMarkup = NewKPIKeyboard()
		b.BotAPI.Send(msg)

	case "potokKPI":
		b.OutPotokKPI(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Выбери должность")
		msg.ReplyMarkup = NewKPIKeyboard()
		b.BotAPI.Send(msg)

	case "smoothedkpiZRU":
		b.SmoothedKpiZRU(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Выбери должность")
		msg.ReplyMarkup = NewKPIKeyboard()
		b.BotAPI.Send(msg)

	case "smoothedkpiRU":
		b.SmoothedKpiRU(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Выбери должность")
		msg.ReplyMarkup = NewKPIKeyboard()
		b.BotAPI.Send(msg)

	case "npoKPI":
		b.NpoKPI(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Выбери должность")
		msg.ReplyMarkup = NewKPIKeyboard()
		b.BotAPI.Send(msg)

	case "cancelkpi":
		b.CancelKpi(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Выбери должность")
		msg.ReplyMarkup = NewKPIKeyboard()
		b.BotAPI.Send(msg)

	}
}

func (b *Bot) HandleCallbackStandard(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID

	switch callback.Data {
	case "standart":
		msg := tgbotapi.NewMessage(chatID, "Текущие нормативы")
		msg.ReplyMarkup = NewStandardKeyboard()
		b.BotAPI.Send(msg)

	case "mez1":
		b.Mez1Info(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Нормативы")
		msg.ReplyMarkup = NewStandardKeyboard()
		b.BotAPI.Send(msg)

	case "mez2":
		b.Mez2Info(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Нормативы")
		msg.ReplyMarkup = NewStandardKeyboard()
		b.BotAPI.Send(msg)

	case "mez3":
		b.Mez3Info(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Нормативы")
		msg.ReplyMarkup = NewStandardKeyboard()
		b.BotAPI.Send(msg)

	case "mez4":
		b.Mez4Info(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Нормативы")
		msg.ReplyMarkup = NewStandardKeyboard()
		b.BotAPI.Send(msg)

	case "mez5":
		b.Mez5Info(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Нормативы")
		msg.ReplyMarkup = NewStandardKeyboard()
		b.BotAPI.Send(msg)

	case "specHran":
		b.SpecInfo(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Нормативы")
		msg.ReplyMarkup = NewStandardKeyboard()
		b.BotAPI.Send(msg)

	case "selectionKGT":
		b.SelectionKGTInfo(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Нормативы")
		msg.ReplyMarkup = NewStandardKeyboard()
		b.BotAPI.Send(msg)

	case "selectionSGT":
		b.SelectionSGTInfo(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Нормативы")
		msg.ReplyMarkup = NewStandardKeyboard()
		b.BotAPI.Send(msg)

	case "selectionIZ":
		b.SelectionIZInfo(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Нормативы")
		msg.ReplyMarkup = NewStandardKeyboard()
		b.BotAPI.Send(msg)

	case "sortKGT":
		b.SortKGTInfo(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Нормативы")
		msg.ReplyMarkup = NewStandardKeyboard()
		b.BotAPI.Send(msg)

	case "packPot":
		b.PackPotInfo(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Нормативы")
		msg.ReplyMarkup = NewStandardKeyboard()
		b.BotAPI.Send(msg)

	case "sortSGT":
		b.SortSGTInfo(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Нормативы")
		msg.ReplyMarkup = NewStandardKeyboard()
		b.BotAPI.Send(msg)

	case "sortLOT":
		b.SortLotInfo(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Нормативы")
		msg.ReplyMarkup = NewStandardKeyboard()
		b.BotAPI.Send(msg)

	case "bufNorm":
		b.BufNormaInfo(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Нормативы")
		msg.ReplyMarkup = NewStandardKeyboard()
		b.BotAPI.Send(msg)

	case "shipNorm":
		b.ShipNormaInfo(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Нормативы")
		msg.ReplyMarkup = NewStandardKeyboard()
		b.BotAPI.Send(msg)

	case "refresh3":
		msg := tgbotapi.NewMessage(chatID, "Выбери норматив для обновления")
		msg.ReplyMarkup = NewStandardRefreshKeyboard()
		b.BotAPI.Send(msg)

	case "mez1UP":
		b.HandleUpdateButtonStan(chatID, "Отбор MEZ-1")
	case "mez2UP":
		b.HandleUpdateButtonStan(chatID, "Отбор MEZ-2")
	case "mez3UP":
		b.HandleUpdateButtonStan(chatID, "Отбор MEZ-3")
	case "mez4UP":
		b.HandleUpdateButtonStan(chatID, "Отбор MEZ-4")
	case "mez5UP":
		b.HandleUpdateButtonStan(chatID, "Отбор MEZ-5А")
	case "specHranUP":
		b.HandleUpdateButtonStan(chatID, "Спецхран")
	case "selectionKGTUP":
		b.HandleUpdateButtonStan(chatID, "Отбор КГТ RACK ")
	case "selectionSGTUP":
		b.HandleUpdateButtonStan(chatID, "Отбор СГТ RACK")
	case "selectionIZUP":
		b.HandleUpdateButtonStan(chatID, "Отбор Изъятия")
	case "sortKGTUP":
		b.HandleUpdateButtonStan(chatID, "Сортировка КГТ")
	case "packPotUP":
		b.HandleUpdateButtonStan(chatID, "Потоварная сортировка(упак)")
	case "sortSGTUP":
		b.HandleUpdateButtonStan(chatID, "Сортировка СГТ")
	case "sortLOTUP":
		b.HandleUpdateButtonStan(chatID, "Сортировка в LOT")
	case "bufNormUP":
		b.HandleUpdateButtonStan(chatID, "Размещение для отгрузки")
	case "shipNormUP":
		b.HandleUpdateButtonStan(chatID, "Отгрузка")

	case "selectionBalkon":
		b.SelectionBalkon(chatID, "")
		msg := tgbotapi.NewMessage(chatID, "Нормативы")
		msg.ReplyMarkup = NewStandardKeyboard()
		b.BotAPI.Send(msg)

	case "selectionBalkonUP":
		b.HandleUpdateButtonStan(chatID, "Отбор КГТ Балкон")
	}
}
