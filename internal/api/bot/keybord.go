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
			tgbotapi.NewInlineKeyboardButtonData("Точки контроля WMS (1-1)", "op_1-1"),
			tgbotapi.NewInlineKeyboardButtonData("Отмененные заказы (2-7)", "op_2-7"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отбор 2.0 (3-3-8)", "op_3-3-8"),
			tgbotapi.NewInlineKeyboardButtonData("АвтоОтбор 2.0 (3-3-0)", "op_3-3-0"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Размещение для отгрузки (3-7-2)", "op_3-7-2"),
			tgbotapi.NewInlineKeyboardButtonData("Отгрузка ТС (3-7-3)", "op_3-7-3"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Единая упаковка (3-10)", "op_3-10"),
			tgbotapi.NewInlineKeyboardButtonData("Детали волн (4-1-2-1-2)", "op_4-1-2-1-2"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Диспетчер отгрузки (4-1-2-3-3)", "op_4-1-2-3-3"),
			tgbotapi.NewInlineKeyboardButtonData("Ручной запуск заказов (4-1-2-8-1-1)", "op_4-1-2-8-1-1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Транзакции по SKU (4-9-1)", "op_4-9-1"),
			tgbotapi.NewInlineKeyboardButtonData("Транзакции по УИТ (4-9-2)", "op_4-9-2"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Просмотр балансов (5-1)", "op_5-1"),
			tgbotapi.NewInlineKeyboardButtonData("Баланс УИТов (5-2-2)", "op_5-2-2"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Проверка посылок (5-3)", "op_5-3"),
			tgbotapi.NewInlineKeyboardButtonData("Точки контроля WMS (7-2-1)", "op_7-2-1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Запирание (9-5)", "op_9-5"),
			tgbotapi.NewInlineKeyboardButtonData("Автостарт (4-1-2-1-1)", "op_4-1-2-1-1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Подтверждение действия (9-7)", "op_9-7"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "back"),
		),
	)
}
func NewOrderStatusesKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Создано внешне (02)", "status:02"),
			tgbotapi.NewInlineKeyboardButtonData("Выпущено (29)", "status:29"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отсортирован (59)", "status:59"),
			tgbotapi.NewInlineKeyboardButtonData("Упакован (65)", "status:65"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отсортирован по СД (68)", "status:68"),
			tgbotapi.NewInlineKeyboardButtonData("Отгрузка завершена (95)", "status:95"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Проверка КИЗ (-7)", "status:-7"),
			tgbotapi.NewInlineKeyboardButtonData("Отбор завершен (55)", "status:55"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Неизвестно (-1)", "status:-1"),
			tgbotapi.NewInlineKeyboardButtonData("Назад", "back"),
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
