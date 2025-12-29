package bot

import (
	"context"
	"fmt"
	"log"

	"github.com/Pavel-Sergeev-ekb/JARVIS_tg-Bot/internal/api/database"
	"github.com/Pavel-Sergeev-ekb/JARVIS_tg-Bot/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v4"
)

func (b *Bot) StartCheck(chatID int64, msg *tgbotapi.Message) (bool, error) {
	// Используем существующее соединение Bot.Conn вместо нового подключения
	if b.Conn == nil {
		log.Printf("Ошибка: соединение с БД не инициализировано")
		b.SendMessage(chatID, "Произошла ошибка при подключении к базе данных", tgbotapi.ModeHTML)
		return false, fmt.Errorf("БД не подключена")
	}

	var exists bool
	err := b.Conn.QueryRow(context.Background(),
		"SELECT TRUE FROM users WHERE chat_id = $1 AND is_approved = TRUE", chatID).Scan(&exists)

	switch {
	case err == pgx.ErrNoRows:
		// Пользователь не найден или не одобрен → проверяем заявки
		var hasRequest bool
		errReq := b.Conn.QueryRow(context.Background(),
			"SELECT TRUE FROM access_requests WHERE user_id = $1", chatID).Scan(&hasRequest)

		if errReq == nil && hasRequest {
			// Заявка уже есть → не отправляем кнопку
			b.SendMessage(chatID, "Ваша заявка на доступ уже отправлена. Ожидайте решения.", tgbotapi.ModeHTML)
			return false, nil
		}

		// Нет заявки → регистрируем и показываем кнопку
		_, errReg := b.Conn.Exec(context.Background(),
			"INSERT INTO users (chat_id, user_name, is_approved) VALUES ($1, $2, FALSE) ON CONFLICT (chat_id) DO NOTHING",
			chatID, msg.From.UserName)
		if errReg != nil {
			log.Printf("Ошибка регистрации пользователя %d: %v", chatID, errReg)
			b.SendMessage(chatID, "Произошла ошибка при регистрации", tgbotapi.ModeHTML)
			return false, errReg
		}

		keyboard := NewDoorKeyboard()
		msg := tgbotapi.NewMessage(chatID, "Доступ ограничен. Нажмите кнопку ниже, чтобы запросить доступ.")
		msg.ReplyMarkup = keyboard
		if _, errSend := b.BotAPI.Send(msg); errSend != nil {
			log.Printf("Ошибка отправки сообщения с клавиатурой: %v", errSend)
		}
		return false, nil

	case err != nil:
		log.Printf("Ошибка проверки пользователя %d: %v", chatID, err)
		b.SendMessage(chatID, "Произошла ошибка при проверке доступа", tgbotapi.ModeHTML)
		return false, err

	default:
		// Пользователь найден и одобрен → разрешаем доступ
		log.Printf("Пользователь %d разрешён", chatID)
		if msg != nil && msg.Text != "" {
			b.SendMessage(chatID, msg.Text, tgbotapi.ModeHTML)
		}
		return true, nil
	}
}

func (b *Bot) ApproveUser(conn *pgx.Conn, userID int64, approverName string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("ошибка загрузки конфигурации: %v", err)
		return nil
	}
	// Обновляем статус в БД
	err = database.SaveUser(userID, "", true) // true — доступ разрешён
	if err != nil {
		log.Printf("Ошибка при одобрении пользователя %d: %v", userID, err)
		b.SendMessage(cfg.MyChatID,
			fmt.Sprintf("Ошибка при одобрении пользователя %d", userID),
			tgbotapi.ModeHTML)
		return nil
	}

	// Уведомляем пользователя
	b.SendMessage(userID,
		"Доступ разрешён! Теперь вы можете пользоваться ботом.",
		tgbotapi.ModeHTML)

	// Уведомляем владельца
	b.SendMessage(cfg.MyChatID,
		fmt.Sprintf("Пользователь %d одобрен. Инициатор: @%s", userID, approverName),
		tgbotapi.ModeHTML)
	return nil
}

func (b *Bot) RejectUser(conn *pgx.Conn, userID int64, rejectorName string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("ошибка загрузки конфигурации: %v", err)
		return nil
	}

	err = database.SaveUser(userID, "", false) // false — доступ запрещён
	if err != nil {
		log.Printf("Ошибка при отклонении пользователя %d: %v", userID, err)
		b.SendMessage(cfg.MyChatID,
			fmt.Sprintf("Ошибка при отклонении пользователя %d", userID),
			tgbotapi.ModeHTML)
		return nil
	}

	// Уведомляем пользователя
	b.SendMessage(userID,
		"На данный момент доступ отклонён. Обратитесь напрямую к владельцу: @Pavel_Sergeev_1",
		tgbotapi.ModeHTML)

	// Уведомляем владельца
	b.SendMessage(cfg.MyChatID,
		fmt.Sprintf("Запрос пользователя %d отклонён. Инициатор: @%s", userID, rejectorName),
		tgbotapi.ModeHTML)
	return nil
}

func (b *Bot) CheckAccess(chatID int64) (bool, error) {
	if b.Conn == nil {
		log.Printf("Ошибка: соединение с БД не инициализировано")
		return false, fmt.Errorf("БД не подключена")
	}

	var isApproved bool
	err := b.Conn.QueryRow(context.Background(),
		"SELECT is_approved FROM users WHERE chat_id = $1", chatID).Scan(&isApproved)

	switch {
	case err == pgx.ErrNoRows:
		// Пользователь не найден в БД → доступ запрещён
		return false, nil
	case err != nil:
		log.Printf("Ошибка проверки доступа для chatID=%d: %v", chatID, err)
		return false, err
	default:
		// Пользователь найден → возвращаем статус доступа
		return isApproved, nil
	}
}

func (b *Bot) WithAccessCheck(next func(tgbotapi.Update)) func(tgbotapi.Update) {
	return func(update tgbotapi.Update) {
		var chatID int64

		// Извлекаем chatID
		if update.Message != nil {
			chatID = update.Message.Chat.ID
		} else if update.CallbackQuery != nil && update.CallbackQuery.Message != nil {
			chatID = update.CallbackQuery.Message.Chat.ID
		} else {
			return
		}

		// Проверяем доступ
		hasAccess, err := b.CheckAccess(chatID)
		if err != nil {
			log.Printf("Ошибка проверки доступа для chatID=%d: %v", chatID, err)
			b.SendMessage(chatID, "Произошла ошибка при проверке доступа. Попробуйте позже.", tgbotapi.ModeHTML)
			return
		}
		if !hasAccess {
			b.SendMessage(chatID, "Доступ ограничен. Отправьте /start для проверки статуса.", tgbotapi.ModeHTML)
			return
		}

		// Вызываем оригинальный хендлер
		next(update)
	}
}

// WrapHandler — удобная обёртка для хендлеров
func (b *Bot) WrapHandler(handler func(tgbotapi.Update)) func(tgbotapi.Update) {
	return b.WithAccessCheck(handler)
}
