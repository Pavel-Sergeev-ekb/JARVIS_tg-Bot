package bot

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Pavel-Sergeev-ekb/JARVIS_tg-Bot/internal/api/database"
	"github.com/Pavel-Sergeev-ekb/JARVIS_tg-Bot/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) UpdateStand(chatID int64, name string, newIndicator float64) error {
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
        UPDATE standart
        SET indicator = $2, 
				user_chatid = $3,  
    		data_at = $4 
        WHERE name = $1
    `

	// Выполняем UPDATE
	result, err := db.Exec(context.Background(),
		query,
		name,
		newIndicator,
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
			"• Новый показатель: <b>%v</b>\n",
		name, newIndicator,
	)
	b.SendMessage(chatID, message, tgbotapi.ModeHTML)
	msg := tgbotapi.NewMessage(chatID, "Выбери норматив для обновления")
	msg.ReplyMarkup = NewStandardRefreshKeyboard()
	b.BotAPI.Send(msg)

	return nil
}

func (b *Bot) handleUpdateInputStand(chatID int64, message string, standartID string) error {
	// Очищаем строку от пробелов
	indicatorStr := strings.TrimSpace(message)

	// Пытаемся преобразовать в число
	newIndicator, err := strconv.ParseFloat(indicatorStr, 64)
	if err != nil {
		b.SendMessage(chatID, "Ошибка: показатель должен быть числом.", tgbotapi.ModeHTML)
		return err
	}

	// Вызываем функцию обновления
	return b.UpdateStand(chatID, standartID, newIndicator)
}

func (b *Bot) HandleUpdateButtonStan(chatID int64, standartID string) {
	// Отправляем запрос на ввод данных
	msg := fmt.Sprintf(
		"Введите новое значение для показателя %s в формате целого числа, без точек и запятых",
		standartID,
	)
	b.SendMessage(chatID, msg, tgbotapi.ModeHTML)

	// Сохраняем контекст ожидания в карту
	b.waitingForStandInput[chatID] = standartID

}
