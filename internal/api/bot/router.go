package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v4"
)

type Bot struct {
	BotAPI                *tgbotapi.BotAPI
	waitingForStandInput  map[int64]string
	waitingForKPIInput    map[int64]string
	waitingForFirstFile   map[int64]bool
	waitingForSecondFile  map[int64]bool
	tempFilePaths         map[int64]map[int]string
	waitingForNeedInput   map[int64]string
	needData              map[int64]*NeedReport
	waitingForStaffCount  map[int64]bool
	waitingForPickedItems map[int64]bool
	staffCount            map[int64]int
	pickedItems           map[int64]int
	Conn                  *pgx.Conn
	commandToHandler      map[string]string
	handlerMap            map[string]func(tgbotapi.Update)
}

func NewOneBot(botAPI *tgbotapi.BotAPI, conn *pgx.Conn) *Bot {
	b := &Bot{
		BotAPI:                botAPI,
		Conn:                  conn,
		waitingForKPIInput:    make(map[int64]string),
		waitingForStandInput:  make(map[int64]string),
		waitingForFirstFile:   make(map[int64]bool),
		waitingForSecondFile:  make(map[int64]bool),
		tempFilePaths:         make(map[int64]map[int]string),
		waitingForNeedInput:   make(map[int64]string),
		needData:              make(map[int64]*NeedReport),
		waitingForStaffCount:  make(map[int64]bool),
		waitingForPickedItems: make(map[int64]bool),
		staffCount:            make(map[int64]int),
		pickedItems:           make(map[int64]int),
		commandToHandler:      make(map[string]string),
		handlerMap:            make(map[string]func(tgbotapi.Update)),
	}

	b.commandToHandler = map[string]string{

		"wms":               "HandleCallbackKeyboard",
		"launch":            "HandleCallbackKeyboard",
		"autopilot":         "HandleCallbackKeyboard",
		"typeWaves":         "HandleCallbackKeyboard",
		"needAK":            "HandleCallbackKeyboard",
		"binding":           "HandleCallbackKeyboard",
		"hours":             "HandleCallbackKeyboard",
		"back_to_autostart": "HandleCallbackKeyboard",
		"hourlyReport":      "HandleCallbackKeyboard",
		"back":              "HandleCallbackKeyboard",
		"potCons1":          "HandleCallbackKeyboard",
		"potCons2":          "HandleCallbackKeyboard",
		"kgt1":              "HandleCallbackKeyboard",
		"kgt2":              "HandleCallbackKeyboard",
		"single1":           "HandleCallbackKeyboard",
		"single2":           "HandleCallbackKeyboard",
		"stage1":            "HandleCallbackKeyboard",
		"stage2":            "HandleCallbackKeyboard",
		"backToAuto":        "HandleCallbackKeyboard",
		"timeBefore":        "HandleCallbackKeyboard",
		"timeWith":          "HandleCallbackKeyboard",
		"handStart":         "HandleCallbackKeyboard",
		"links":             "HandleCallbackKeyboard",
		"reply":             "HandleCallbackKeyboard",
		"tickets":           "HandleCallbackKeyboard",

		"operations":       "HandleCalbackOperation",
		"controlPoints":    "HandleCalbackOperation",
		"canceled":         "HandleCalbackOperation",
		"selection":        "HandleCalbackOperation",
		"autoSelection":    "HandleCalbackOperation",
		"door":             "HandleCalbackOperation",
		"shipped":          "HandleCalbackOperation",
		"packing":          "HandleCalbackOperation",
		"waveDetails":      "HandleCalbackOperation",
		"dispatcher":       "HandleCalbackOperation",
		"handStartInfo":    "HandleCalbackOperation",
		"transactionSKU":   "HandleCalbackOperation",
		"transactionUIT":   "HandleCalbackOperation",
		"balance":          "HandleCalbackOperation",
		"balanceUIT":       "HandleCalbackOperation",
		"checkOrders":      "HandleCalbackOperation",
		"controlPointsWMS": "HandleCalbackOperation",
		"locking":          "HandleCalbackOperation",
		"autoStart":        "HandleCalbackOperation",
		"action":           "HandleCalbackOperation",
		"skills":           "HandleCalbackOperation",
		"AssignedWork":     "HandleCalbackOperation",
		"hiring":           "HandleCalbackOperation",
		"backToWMS":        "HandleCalbackOperation",

		"createdExt":         "HandleCallbackOrders",
		"orderStatuses":      "HandleCallbackOrders",
		"released":           "HandleCallbackOrders",
		"sorted":             "HandleCallbackOrders",
		"packed":             "HandleCallbackOrders",
		"sortedSD":           "HandleCallbackOrders",
		"shipment":           "HandleCallbackOrders",
		"KIZ":                "HandleCallbackOrders",
		"selectionCompleted": "HandleCallbackOrders",
		"Unknown":            "HandleCallbackOrders",

		"tutorial":        "HandleCallbackTutorial",
		"primer":          "HandleCallbackTutorial",
		"backToTutorial":  "HandleCallbackTutorial",
		"shippedIns":      "HandleCallbackTutorial",
		"fact":            "HandleCallbackTutorial",
		"transactionIns":  "HandleCallbackTutorial",
		"seal":            "HandleCallbackTutorial",
		"backToShip":      "HandleCallbackTutorial",
		"Replenishments":  "HandleCallbackTutorial",
		"kpi":             "HandleCallbackTutorial",
		"sla":             "HandleCallbackTutorial",
		"lost":            "HandleCallbackTutorial",
		"capasity":        "HandleCallbackTutorial",
		"dayoff":          "HandleCallbackTutorial",
		"dashboards":      "HandleCallbackTutorial",
		"alarm":           "HandleCallbackTutorial",
		"NpoInfo":         "HandleCallbackTutorial",
		"smoothed":        "HandleCallbackTutorial",
		"shipmentPro":     "HandleCallbackTutorial",
		"OpenPro":         "HandleCallbackTutorial",
		"fte":             "HandleCallbackTutorial",
		"germes":          "HandleCallbackTutorial",
		"ticket":          "HandleCallbackTutorial",
		"withdrawal":      "HandleCallbackTutorial",
		"infoscan":        "HandleCallbackTutorial",
		"catoff":          "HandleCallbackTutorial",
		"spichRichtrack":  "HandleCallbackTutorial",
		"shortBriefing":   "HandleCallbackTutorial",
		"longBriefing":    "HandleCallbackTutorial",
		"formsOT":         "HandleCallbackTutorial",
		"ekb":             "HandleCallbackTutorial",
		"transport":       "HandleCallbackTutorial",
		"tranzit":         "HandleCallbackTutorial",
		"system":          "HandleCallbackTutorial",
		"skillsPro":       "HandleCallbackTutorial",
		"lockingINFO":     "HandleCallbackTutorial",
		"YP":              "HandleCallbackTutorial",
		"UIT":             "HandleCallbackTutorial",
		"ROV":             "HandleCallbackTutorial",
		"Tara":            "HandleCallbackTutorial",
		"nuance":          "HandleCallbackTutorial",
		"DRP":             "HandleCallbackTutorial",
		"TRP":             "HandleCallbackTutorial",
		"uitSeal":         "HandleCallbackTutorial",
		"nzn":             "HandleCallbackTutorial",
		"label":           "HandleCallbackTutorial",
		"clientOrder":     "HandleCallbackTutorial",
		"withdrawalOrder": "HandleCallbackTutorial",

		"KPI":              "HandleCallbackKPI",
		"zru":              "HandleCallbackKPI",
		"brigadir":         "HandleCallbackKPI",
		"backToKPI":        "HandleCallbackKPI",
		"dispatchZRU":      "HandleCallbackKPI",
		"dispatchRU":       "HandleCallbackKPI",
		"refresh1":         "HandleCallbackKPI",
		"dispatchUpZRU":    "HandleCallbackKPI",
		"dispatchUpRU":     "HandleCallbackKPI",
		"vpKPI":            "HandleCallbackKPI",
		"vpKPIUP":          "HandleCallbackKPI",
		"cancelkpiUP":      "HandleCallbackKPI",
		"lostKPIUp":        "HandleCallbackKPI",
		"timelinesskpiUp":  "HandleCallbackKPI",
		"potokKPIUp":       "HandleCallbackKPI",
		"smoothedkpiUpZRU": "HandleCallbackKPI",
		"smoothedkpiUpRU":  "HandleCallbackKPI",
		"npoKPIUp":         "HandleCallbackKPI",
		"lostKPI":          "HandleCallbackKPI",
		"timelinesskpi":    "HandleCallbackKPI",
		"potokKPI":         "HandleCallbackKPI",
		"smoothedkpiZRU":   "HandleCallbackKPI",
		"smoothedkpiRU":    "HandleCallbackKPI",
		"npoKPI":           "HandleCallbackKPI",
		"cancelkpi":        "HandleCallbackKPI",

		"standart":          "HandleCallbackStandart",
		"mez1":              "HandleCallbackStandart",
		"mez2":              "HandleCallbackStandart",
		"mez3":              "HandleCallbackStandart",
		"mez4":              "HandleCallbackStandart",
		"mez5":              "HandleCallbackStandart",
		"specHran":          "HandleCallbackStandart",
		"selectionKGT":      "HandleCallbackStandart",
		"selectionSGT":      "HandleCallbackStandart",
		"selectionIZ":       "HandleCallbackStandart",
		"sortKGT":           "HandleCallbackStandart",
		"packPot":           "HandleCallbackStandart",
		"sortSGT":           "HandleCallbackStandart",
		"sortLot":           "HandleCallbackStandart",
		"bufNorm":           "HandleCallbackStandart",
		"shipNorm":          "HandleCallbackStandart",
		"refresh3":          "HandleCallbackStandart",
		"mez1UP":            "HandleCallbackStandart",
		"mez2UP":            "HandleCallbackStandart",
		"mez3UP":            "HandleCallbackStandart",
		"mez4UP":            "HandleCallbackStandart",
		"mez5UP":            "HandleCallbackStandart",
		"specHranUP":        "HandleCallbackStandart",
		"selectionKGTUP":    "HandleCallbackStandart",
		"selectionSGTUP":    "HandleCallbackStandart",
		"selectionIZUP":     "HandleCallbackStandart",
		"sortKGTUP":         "HandleCallbackStandart",
		"packPotUP":         "HandleCallbackStandart",
		"sortSGTUP":         "HandleCallbackStandart",
		"sortLotUP":         "HandleCallbackStandart",
		"bufNormUP":         "HandleCallbackStandart",
		"shipNormUP":        "HandleCallbackStandart",
		"selectionBalkon":   "HandleCallbackStandart",
		"selectionBalkonUP": "HandleCallbackStandart",

		"approve": "HandleAccessCallback",
		"reject":  "HandleAccessCallback",
		"access":  "HandleAccessCallback",
	}

	b.handlerMap = map[string]func(tgbotapi.Update){
		"HandleCallbackKeyboard": b.HandleCallbackKeyboard,
		"HandleCalbackOperation": b.HandleCallbackOperation,
		"HandleCallbackOrders":   b.HandleCallbackOrders,
		"HandleCallbackTutorial": b.HandleCallbackTutorial,
		"HandleCallbackKPI":      b.HandleCallbackKPI,
		"HandleCallbackStandart": b.HandleCallbackStandard,
		"HandleAccessCallback":   b.HandleAccessCallback,
	}

	return b
}

func (b *Bot) RouteCallback(update tgbotapi.Update) {
	if update.CallbackQuery == nil {
		return
	}

	callbackData := update.CallbackQuery.Data

	// 1. Ищем имя хендлера по команде
	handlerName, ok := b.commandToHandler[callbackData]
	if !ok {
		b.SendMessage(
			update.CallbackQuery.Message.Chat.ID,
			"Неизвестная команда",
			tgbotapi.ModeHTML,
		)
		return
	}

	// 2. Ищем функцию‑обработчик по имени
	handler, ok := b.handlerMap[handlerName]
	if !ok {
		b.SendMessage(
			update.CallbackQuery.Message.Chat.ID,
			"Хендлер не найден",
			tgbotapi.ModeHTML,
		)
		return
	}

	// 3. Вызываем хендлер
	handler(update)
}
