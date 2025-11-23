package bot

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Pavel-Sergeev-ekb/JARVIS_tg-Bot/internal/api/database"
	"github.com/Pavel-Sergeev-ekb/JARVIS_tg-Bot/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) updateKPI(chatID int64, name string, newIndicator float64, newWeight int) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Ошибка загрузки конфигурации: %v", err)
		b.SendMessage(chatID, "Произошла ошибка при загрузке конфигурации", tgbotapi.ModeHTML)
		return err
	}

	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Printf("Ошибка подключения к БД: %v", err)
		b.SendMessage(chatID, "Произошла ошибка при подключении к базе данных", tgbotapi.ModeHTML)
		return err
	}
	defer db.Close(context.Background())

	// SQL-запрос на обновление
	query := `
        UPDATE kpi
        SET indicator = $2, 
				weight = $3,
				user_chatid = $4,  
    		data = $5 
        WHERE name = $1
    `

	// Выполняем UPDATE
	result, err := db.Exec(context.Background(),
		query,
		name,
		newIndicator,
		newWeight,
		chatID,
		time.Now(),
	)
	if err != nil {
		log.Printf("Ошибка при выполнении UPDATE-запроса: %v", err)
		b.SendMessage(chatID, "Произошла ошибка при обновлении данных", tgbotapi.ModeHTML)
		return err
	}

	// Проверяем, сколько строк было обновлено
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		msg := fmt.Sprintf("Показатель с ID %s не найден", name)
		b.SendMessage(chatID, msg, tgbotapi.ModeHTML)
		return fmt.Errorf("показатель не найден в базе данных")
	}

	// Формируем сообщение об успешном обновлении
	message := fmt.Sprintf(
		"Показатель <b>%s</b> успешно обновлён:\n"+
			"• Новый показатель: <b>%v</b>\n"+
			"• Новый вес: <b>%v%%</b>",
		name, newIndicator, newWeight,
	)
	b.SendMessage(chatID, message, tgbotapi.ModeHTML)
	msg := tgbotapi.NewMessage(chatID, "Выбери должность")
	msg.ReplyMarkup = NewKPIKeyboard()
	b.BotAPI.Send(msg)

	return nil
}

func (b *Bot) handleUpdateInput(chatID int64, message string, kpiID string) error {
	// Разбиваем входное сообщение по запятой
	parts := strings.Split(message, ",")
	if len(parts) != 2 {
		b.SendMessage(chatID, "Неверный формат ввода. Напишите числа через запятую: показатель, вес (например: 98.5, 30)", tgbotapi.ModeHTML)
		return errors.New("неверный формат входных данных")
	}

	// Очищаем строки от пробелов и пытаемся преобразовать в числа
	indicatorStr := strings.TrimSpace(parts[0])
	weightStr := strings.TrimSpace(parts[1])

	newIndicator, err := strconv.ParseFloat(indicatorStr, 64)
	if err != nil {
		b.SendMessage(chatID, "Ошибка: показатель должен быть числом.", tgbotapi.ModeHTML)
		return err
	}

	newWeight, err := strconv.Atoi(weightStr)
	if err != nil || newWeight < 0 || newWeight > 100 {
		b.SendMessage(chatID, "Ошибка: вес должен быть целым числом от 0 до 100.", tgbotapi.ModeHTML)
		return err
	}

	// Вызываем функцию обновления
	return b.updateKPI(chatID, kpiID, newIndicator, newWeight)
}
func (b *Bot) HandleUpdateButton(chatID int64, kpiID string) {
	// Отправляем запрос на ввод данных
	msg := fmt.Sprintf(
		"Введите новые значения для показателя %s в формате:\n"+
			"Показатель, Вес (например: 98.5, 30)\n\n"+
			"Показатель — десятичное или целое число, вес — целое число (0–100).",
		kpiID,
	)
	b.SendMessage(chatID, msg, tgbotapi.ModeHTML)

	// Сохраняем контекст ожидания в карту
	b.waitingForKPIInput[chatID] = kpiID

}
