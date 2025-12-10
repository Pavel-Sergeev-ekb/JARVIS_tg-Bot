package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func NewMainKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("WMS", "wms"),
			tgbotapi.NewInlineKeyboardButtonData("Туториал", "tutorial"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("KPI", "KPI"),
			tgbotapi.NewInlineKeyboardButtonData("Нормативы", "standart"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отчеты", "reply"),
			tgbotapi.NewInlineKeyboardButtonData("Полезные ссылки", "links"),
		),
	)
}

func NewLinkMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Борды", "dashboards"),
			tgbotapi.NewInlineKeyboardButtonData("Тикеты", "tickets"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Основные формы по ОТ", "formsOT"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "back"),
		),
	)
}

func NewReplyMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Часовой отчет", "hourlyReport"),
			tgbotapi.NewInlineKeyboardButtonData("Своевременность", "timeliness"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Потребность", "needAK"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "back"),
		),
	)
}

func NewWMSMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Операции", "operations"),
			tgbotapi.NewInlineKeyboardButtonData("Статусы заказов", "orderStatuses"),
			tgbotapi.NewInlineKeyboardButtonData("Запуск", "launch"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "back"),
		),
	)
}

func NewWMSOperationsKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Точки контроля", "controlPoints"),
			tgbotapi.NewInlineKeyboardButtonData("Отмененные заказы", "canceled"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отбор 2.0", "selection"),
			tgbotapi.NewInlineKeyboardButtonData("АвтоОтбор 2.0", "autoSelection"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Размещение для отгрузки", "door"),
			tgbotapi.NewInlineKeyboardButtonData("Отгрузка ТС", "shipped"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Единая упаковка", "packing"),
			tgbotapi.NewInlineKeyboardButtonData("Детали волн", "waveDetails"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Диспетчер отгрузки", "dispatcher"),
			tgbotapi.NewInlineKeyboardButtonData("Ручной запуск заказов ", "handStartInfo"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Транзакции по SKU", "transactionSKU"),
			tgbotapi.NewInlineKeyboardButtonData("Транзакции по УИТ", "transactionUIT"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Просмотр балансов", "balance"),
			tgbotapi.NewInlineKeyboardButtonData("Баланс УИТов", "balanceUIT"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Проверка посылок", "checkOrders"),
			tgbotapi.NewInlineKeyboardButtonData("Точки контроля WMS", "controlPointsWMS"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Запирание", "locking"),
			tgbotapi.NewInlineKeyboardButtonData("Автостарт", "autoStart"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Подтверждение действия", "action"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Навыки", "skills"),
			tgbotapi.NewInlineKeyboardButtonData("Назначенная работа", "AssignedWork"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Принятие сотрудника", "hiring"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "backToWMS"),
		),
	)
}
func NewOrderStatusesKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Создано внешне", "createdExt"),
			tgbotapi.NewInlineKeyboardButtonData("Выпущено", "released"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отсортирован", "sorted"),
			tgbotapi.NewInlineKeyboardButtonData("Упакован", "packed"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отсортирован по СД", "sortedSD"),
			tgbotapi.NewInlineKeyboardButtonData("Отгрузка завершена", "shipment"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Проверка КИЗ", "KIZ"),
			tgbotapi.NewInlineKeyboardButtonData("Отбор завершен", "selectionCompleted"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Неизвестно ", "Unknown"),
			tgbotapi.NewInlineKeyboardButtonData("Назад", "backToWMS"),
		),
	)
}

func NewLaunchKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Автостарт", "autopilot"),
			tgbotapi.NewInlineKeyboardButtonData("Ручной запуск", "handStart"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "back"),
		),
	)
}

func NewAutostartKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Типы волн", "typeWaves"),
			tgbotapi.NewInlineKeyboardButtonData("Привязки", "binding"),
			tgbotapi.NewInlineKeyboardButtonData("Часы", "hours"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "back"),
		),
	)
}

func NewWavesKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Потоварка", "potCons1"),
			tgbotapi.NewInlineKeyboardButtonData("КГТ", "kgt1"),
			tgbotapi.NewInlineKeyboardButtonData("Сингл", "single1"),
			tgbotapi.NewInlineKeyboardButtonData("Станция", "stage1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "backToAuto"),
		),
	)
}

func NewBindKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Потоварка", "potCons2"),
			tgbotapi.NewInlineKeyboardButtonData("КГТ", "kgt2"),
			tgbotapi.NewInlineKeyboardButtonData("Сингл", "single2"),
			tgbotapi.NewInlineKeyboardButtonData("Станция", "stage2"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "backToAuto"),
		),
	)
}

func NewHoursKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Время ДО", "timeBefore"),
			tgbotapi.NewInlineKeyboardButtonData("Время С", "timeWith"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "backToAuto"),
		),
	)
}

func NewTutorialKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Букварь", "primer"),
			tgbotapi.NewInlineKeyboardButtonData("Отгрузка", "shippedIns"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Транзакции", "transactionIns"),
			tgbotapi.NewInlineKeyboardButtonData("Печать", "seal"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Пополнения", "Replenishments"),
			tgbotapi.NewInlineKeyboardButtonData("Запирание", "lockingINFO"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Навыки и скиллы сотрудников", "skillsPro"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Инструктаж ВЭШ", "spichRichtrack"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "back"),
		),
	)
}

func NewPrimerKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("KPI", "kpi"),
			tgbotapi.NewInlineKeyboardButtonData("SLA", "sla"),
			tgbotapi.NewInlineKeyboardButtonData("LOST", "lost"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Капасити", "capasity"),
			tgbotapi.NewInlineKeyboardButtonData("Факап", "alarm"),
			tgbotapi.NewInlineKeyboardButtonData("НПО", "NpoInfo"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Катофф", "catoff"),
			tgbotapi.NewInlineKeyboardButtonData("Дейофф", "dayoff"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Сглаженный производ", "smoothed"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Производ отгрузки", "shipmentPro"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Производ входящего", "OpenPro"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Тикет", "ticket"),
			tgbotapi.NewInlineKeyboardButtonData("Гермес", "germes"),
			tgbotapi.NewInlineKeyboardButtonData("FTE", "fte"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Изъятие", "withdrawal"),
			tgbotapi.NewInlineKeyboardButtonData("Инфоскан", "infoscan"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "backToTutorial"),
		),
	)

}

func NewShippedInsKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Физика", "fact"),
			tgbotapi.NewInlineKeyboardButtonData("Система", "system"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "backToTutorial"),
		),
	)
}

func NewFactInsKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Екат", "ekb"),
			tgbotapi.NewInlineKeyboardButtonData("ТС", "transport"),
			tgbotapi.NewInlineKeyboardButtonData("Транзит", "tranzit"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "backToShip"),
		),
	)
}

func NewTransactionInsMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("YP", "YP"),
			tgbotapi.NewInlineKeyboardButtonData("УИТ", "UIT"),
			tgbotapi.NewInlineKeyboardButtonData("ROV", "ROV"),
			tgbotapi.NewInlineKeyboardButtonData("Тара", "Tara"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Нюансы", "nuance"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "backToTutorial"),
		),
	)
}

func NewSealKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("DRP", "DRP"),
			tgbotapi.NewInlineKeyboardButtonData("TRP", "TRP"),
			tgbotapi.NewInlineKeyboardButtonData("УИТ", "uitSeal"),
			tgbotapi.NewInlineKeyboardButtonData("Контейнер", "nzn"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Посылочная этикетка(Pшка)", "label"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "backToTutorial"),
		),
	)
}

func NewReplenishKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Клиентский заказ", "clientOrder"),
			tgbotapi.NewInlineKeyboardButtonData("Изъятие", "withdrawalOrder"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "backToTutorial"),
		),
	)
}

func NewKPIKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Старший смены(ЗРУ)", "zru"),
			tgbotapi.NewInlineKeyboardButtonData("РУ и бригадир", "brigadir"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "back"),
		),
	)
}

func NewZruKpiKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("% Засылов", "dispatchZRU"),
			tgbotapi.NewInlineKeyboardButtonData("Доля НПО", "npoKPI"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Сглаженный Производ", "smoothedkpiZRU"),
			tgbotapi.NewInlineKeyboardButtonData("% Отмен", "cancelkpi"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Своевременность Отгрузки", "timelinesskpi"),
			tgbotapi.NewInlineKeyboardButtonData("Своевременность приемки и размещения", "vpKPI"),
		),

		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("ОБНОВИТЬ", "refresh1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "backToKPI"),
		),
	)
}
func NewBrigadirKpiKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("% Засылов", "dispatchRU"),
			tgbotapi.NewInlineKeyboardButtonData("Доля НПО", "npoKPI"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Сглаженный Производ", "smoothedkpiRU"),
			tgbotapi.NewInlineKeyboardButtonData("Производ Исхода", "potokKPI"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Своевременность Отгрузки", "timelinesskpi"),
			tgbotapi.NewInlineKeyboardButtonData("% Потерь (склад)", "lostKPI"),
		),

		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("ОБНОВИТЬ", "refresh"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "backToKPI"),
		),
	)
}
func NewRefreshKpiZRUKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("% Засылов", "dispatchUpZRU"),
			tgbotapi.NewInlineKeyboardButtonData("Доля НПО", "npoKPIUp"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Сглаженный Производ", "smoothedkpiUpZRU"),
			tgbotapi.NewInlineKeyboardButtonData("Производ Возвратов", "vpKPIUP"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Своевременность Отгрузки", "timelinesskpiUp"),
			tgbotapi.NewInlineKeyboardButtonData("% Отмен", "cancelkpiUP"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "backToKPI"),
		),
	)
}
func NewRefreshKpiRUKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("% Засылов", "dispatchUpRU"),
			tgbotapi.NewInlineKeyboardButtonData("Доля НПО", "npoKPIUp"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Сглаженный Производ", "smoothedkpiUpRU"),
			tgbotapi.NewInlineKeyboardButtonData("Производ Исхода", "potokKPIUp"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Своевременность Отгрузки", "timelinesskpiUp"),
			tgbotapi.NewInlineKeyboardButtonData("% Потерь (склад)", "lostKPIUp"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "backToKPI"),
		),
	)
}

func NewStandardKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отбор MEZ-1", "mez1"),
			tgbotapi.NewInlineKeyboardButtonData("Отбор MEZ-2", "mez2"),
			tgbotapi.NewInlineKeyboardButtonData("Отбор MEZ-3", "mez3"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отбор MEZ-4", "mez4"),
			tgbotapi.NewInlineKeyboardButtonData("Отбор MEZ-5А", "mez5"),
			tgbotapi.NewInlineKeyboardButtonData("Спецхран", "specHran"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отбор КГТ RACK", "selectionKGT"),
			tgbotapi.NewInlineKeyboardButtonData("Отбор CГТ RACK", "selectionSGT"),
			tgbotapi.NewInlineKeyboardButtonData("Отбор КГТ Балкон", "selectionBalkon"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отбор Изъятия", "selectionIZ"),
			tgbotapi.NewInlineKeyboardButtonData("Cортировка КГТ", "sortKGT"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Потоварная сортировка(упак)", "packPot"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Сортировка СГТ", "sortSGT"),
			tgbotapi.NewInlineKeyboardButtonData("Сортировка в LOT", "sortLOT"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Размещение для отгрузки", "bufNorm"),
			tgbotapi.NewInlineKeyboardButtonData("Отгрузка", "shipNorm"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("ОБНОВИТЬ", "refresh3"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "back"),
		),
	)
}

func NewVashKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Общий спич", "longBriefing"),
			tgbotapi.NewInlineKeyboardButtonData("Короткий спич 2.0", "shortBriefing"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "backToTutorial"),
		),
	)
}

func NewStandardRefreshKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отбор MEZ-1", "mez1UP"),
			tgbotapi.NewInlineKeyboardButtonData("Отбор MEZ-2", "mez2UP"),
			tgbotapi.NewInlineKeyboardButtonData("Отбор MEZ-3", "mez3UP"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отбор MEZ-4", "mez4UP"),
			tgbotapi.NewInlineKeyboardButtonData("Отбор MEZ-5А", "mez5UP"),
			tgbotapi.NewInlineKeyboardButtonData("Спецхран", "specHranUP"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отбор КГТ RACK", "selectionKGTUP"),
			tgbotapi.NewInlineKeyboardButtonData("Отбор CГТ RACK", "selectionSGTUP"),
			tgbotapi.NewInlineKeyboardButtonData("Отбор КГТ Балкон", "selectionBalkonUP"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отбор Изъятия", "selectionIZUP"),
			tgbotapi.NewInlineKeyboardButtonData("Cортировка КГТ", "sortKGTUP"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Потоварная сортировка(упак)", "packPotUP"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Сортировка СГТ", "sortSGTUP"),
			tgbotapi.NewInlineKeyboardButtonData("Сортировка в LOT", "sortLOTUP"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Размещение для отгрузки", "bufNormUP"),
			tgbotapi.NewInlineKeyboardButtonData("Отгрузка", "shipNormUP"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "back"),
		),
	)
}
