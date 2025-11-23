package bot

import (
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type NeedReport struct {
	TotalDrop   int
	DayShift    ShiftData
	NightShift  ShiftData // будет рассчитан автоматически
	ActualStaff int
}

type ShiftData struct {
	Pick  int     // Отбор
	Pack  int     // Упаковка
	Sort  int     // Сортировка
	Staff float64 // Рассчитанная потребность
}

func (b *Bot) startNeedReport(chatID int64) {
	if _, inProgress := b.waitingForNeedInput[chatID]; inProgress {
		b.BotAPI.Send(tgbotapi.NewMessage(chatID, "Завершите текущий отчёт!"))
		return
	}
	b.needData[chatID] = &NeedReport{}
	b.waitingForNeedInput[chatID] = "total_drop"

	msg := tgbotapi.NewMessage(chatID, "Введите общее падение за день (шт.):")
	b.BotAPI.Send(msg)
}

func (b *Bot) processNeedInput(chatID int64, stage, text string) {
	data := b.needData[chatID]

	switch stage {
	case "total_drop":
		if num, err := strconv.Atoi(text); err == nil && num > 0 {
			data.TotalDrop = num
			b.sendDayShiftPrompt(chatID)
		} else {
			b.SendMessage(chatID, "Ошибка: введите положительное число!", tgbotapi.ModeHTML)
		}

	case "day_pick":
		if num, err := strconv.Atoi(text); err == nil && num >= 0 {
			data.DayShift.Pick = num
			b.sendNextDayShiftPrompt(chatID, "pack")
		} else {
			b.SendMessage(chatID, "Ошибка: введите число ≥ 0!", tgbotapi.ModeHTML)
		}

	case "day_pack":
		if num, err := strconv.Atoi(text); err == nil && num >= 0 {
			data.DayShift.Pack = num
			b.sendNextDayShiftPrompt(chatID, "sort")
		} else {
			b.SendMessage(chatID, "Ошибка: введите число ≥ 0!", tgbotapi.ModeHTML)
		}

	case "day_sort":
		if num, err := strconv.Atoi(text); err == nil && num >= 0 {
			data.DayShift.Sort = num
			b.askActualStaff(chatID) // переходим к запросу фактического штата
		} else {
			b.SendMessage(chatID, "Ошибка: введите число ≥ 0!", tgbotapi.ModeHTML)
		}

	case "actual_staff":
		if num, err := strconv.Atoi(text); err == nil && num >= 0 {
			data.ActualStaff = num
			b.calculateNightShift(data) // рассчитываем ночную смену
			b.generateAndSendNeedReport(chatID)
		} else {
			b.SendMessage(chatID, "Ошибка: введите число ≥ 0!", tgbotapi.ModeHTML)
		}
	}
}
func (b *Bot) calculateNightShift(data *NeedReport) {
	// Отбор: чтобы осталось 0
	nightPick := data.TotalDrop - data.DayShift.Pick
	if nightPick < 0 {
		nightPick = 0
	}

	// Упаковка: чтобы осталось 1000
	nightPack := data.TotalDrop - data.DayShift.Pack - 1000
	if nightPack < 0 {
		nightPack = 0
	}

	// Сортировка: чтобы осталось 3000
	nightSort := data.TotalDrop - data.DayShift.Sort - 3000
	if nightSort < 0 {
		nightSort = 0
	}

	data.NightShift.Pick = nightPick
	data.NightShift.Pack = nightPack
	data.NightShift.Sort = nightSort

	// Расчёт потребности в персонале (как раньше)
	data.NightShift.Staff = float64(nightPick)/11/165 +
		float64(nightPack)/11/360 +
		float64(nightSort)/11/550
}

func (b *Bot) generateAndSendNeedReport(chatID int64) {
	data := b.needData[chatID]
	totalStaff := data.NightShift.Staff
	shortage := totalStaff - float64(data.ActualStaff)

	report := fmt.Sprintf(`<b>ОТЧЁТ ПО ПОТРЕБНОСТИ В ПЕРСОНАЛЕ</b>
	───────────────────────────────────────────────

		<b>Общее падение:</b> %dшт.

		<b>Дневная смена:</b>

		 - Отбор: <b>%dшт.</b>
		 - Упаковка: <b>%dшт.</b>
		- Сортировка: <b>%dшт.</b>


		<b>Ночная смена (план для достижения остатков):</b>

		- Отбор: <b>%dшт.</b> (чтобы осталось <b>0</b>)
		- Упаковка: <b>%dшт.</b> (чтобы осталось <b>1000</b>)
		- Сортировка: <b>%dшт.</b> (чтобы осталось <b>3000</b>)


		<b>Итого потребность (ночная смена):</b> %.1fчел.

		<b>Фактически заявлено:</b> %dчел.

		<b>Необходимо добавить:</b> <b>+%.1fчел.</b>`,
		data.TotalDrop,
		data.DayShift.Pick, data.DayShift.Pack, data.DayShift.Sort,
		data.NightShift.Pick, data.NightShift.Pack, data.NightShift.Sort,
		totalStaff, data.ActualStaff, shortage)

	b.SendMessage(chatID, report, tgbotapi.ModeHTML)

	msg := tgbotapi.NewMessage(chatID, "Выбери раздел:")
	msg.ReplyMarkup = NewMainKeyboard()
	b.BotAPI.Send(msg)

	delete(b.waitingForNeedInput, chatID)
	delete(b.needData, chatID)
}

func (b *Bot) sendDayShiftPrompt(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "План для дневной смены:\n- Отбор (шт.):")
	b.BotAPI.Send(msg)
	b.waitingForNeedInput[chatID] = "day_pick"
}

func (b *Bot) sendNextDayShiftPrompt(chatID int64, next string) {
	switch next {
	case "pack":
		msg := tgbotapi.NewMessage(chatID, "- Упаковка (шт.):")
		b.BotAPI.Send(msg)
		b.waitingForNeedInput[chatID] = "day_pack"
	case "sort":
		msg := tgbotapi.NewMessage(chatID, "- Сортировка (шт.):")
		b.BotAPI.Send(msg)
		b.waitingForNeedInput[chatID] = "day_sort"
	}
}

func (b *Bot) sendNightShiftPrompt(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "План для ночной смены:\n- Отбор (шт.):")
	b.BotAPI.Send(msg)
	b.waitingForNeedInput[chatID] = "night_pick"
}

func (b *Bot) sendNextNightShiftPrompt(chatID int64, next string) {
	switch next {
	case "pack":
		msg := tgbotapi.NewMessage(chatID, "- Упаковка (шт.):")
		b.BotAPI.Send(msg)
		b.waitingForNeedInput[chatID] = "night_pack"
	case "sort":
		msg := tgbotapi.NewMessage(chatID, "- Сортировка (шт.):")
		b.BotAPI.Send(msg)
		b.waitingForNeedInput[chatID] = "night_sort"
	}
}

func (b *Bot) askActualStaff(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Фактическое количество заявленных людей:")
	b.BotAPI.Send(msg)
	b.waitingForNeedInput[chatID] = "actual_staff"
}
