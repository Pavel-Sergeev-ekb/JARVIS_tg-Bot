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
			tgbotapi.NewInlineKeyboardButtonData("KPI", "kpi"),
			tgbotapi.NewInlineKeyboardButtonData("Нормативы", "standart"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Часовой отчет", "hourlyReport"),
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
			tgbotapi.NewInlineKeyboardButtonData("Потоварка", "potCons"),
			tgbotapi.NewInlineKeyboardButtonData("КГТ", "kgt"),
			tgbotapi.NewInlineKeyboardButtonData("Сингл", "single"),
			tgbotapi.NewInlineKeyboardButtonData("Станция", "stage"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "back"),
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
			tgbotapi.NewInlineKeyboardButtonData("Назад", "back"),
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
			tgbotapi.NewInlineKeyboardButtonData("Запирание", "lockingIns"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Навыки и скиллы сотрудников", "skillsPro"),
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
			tgbotapi.NewInlineKeyboardButtonData("Дейофф", "lockingIns"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Сглаженный производ", "smoothed"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Производ отгрузки", "shipmentPro"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Тикет", "ticket"),
			tgbotapi.NewInlineKeyboardButtonData("Гермес", "germes"),
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
			tgbotapi.NewInlineKeyboardButtonData("% Засылов", "dispatch"),
			tgbotapi.NewInlineKeyboardButtonData("Доля НПО", "npoKPI"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Сглаженный Производ", "smothedkpi"),
			tgbotapi.NewInlineKeyboardButtonData("% Отмен", "cancelkpi"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Своевременность Отгрузки", "timelinesskpi"),
			tgbotapi.NewInlineKeyboardButtonData("Производ Возвратов", "vpKPI"),
		),

		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Обновить", "refresh1"),
			tgbotapi.NewInlineKeyboardButtonData("Назад", "back"),
		),
	)
}
func NewBrigadirKpiKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("% Засылов", "dispatch"),
			tgbotapi.NewInlineKeyboardButtonData("Доля НПО", "npoKPI"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Сглаженный Производ", "smothedkpi"),
			tgbotapi.NewInlineKeyboardButtonData("Производ Исхода", "potokKPI"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Своевременность Отгрузки", "timelinesskpi"),
			tgbotapi.NewInlineKeyboardButtonData("% Потерь (склад)", "vpKPI"),
		),

		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Обновить", "refresh1"),
			tgbotapi.NewInlineKeyboardButtonData("Назад", "back"),
		),
	)
}
