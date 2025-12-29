package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Pavel-Sergeev-ekb/JARVIS_tg-Bot/internal/config"
	"github.com/jackc/pgx/v4"
)

type AccessRequest struct {
	ID          int64     `db:"id"`
	UserID      string    `db:"user_id"`
	Status      string    `db:"status"` // "pending", "approved", "rejected"
	RequestedAt time.Time `db:"requested_at"`
	ApprovedBy  string    `db:"approved_by"`
	UpdatedAt   time.Time `db:"updated_at"`
	UserName    string    `db:"user_name"`
}

func SaveUser(chatID int64, userName string, isApproved bool) error {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Ошибка загрузки конфигурации: %v", err)
		return err
	}

	db, err := ConnectDB(cfg)
	if err != nil {
		log.Printf("Ошибка подключения к БД: %v", err)
		return err
	}
	defer db.Close(context.Background())

	query := `
        INSERT INTO users (chat_id, user_name, is_approved)
		VALUES ($1, $2, $3)
		ON CONFLICT (chat_id) DO UPDATE
		SET user_name = EXCLUDED.user_name,
		    is_approved = EXCLUDED.is_approved,
		    updated_at = NOW()
    `

	_, err = db.Exec(context.Background(), query, chatID, userName, isApproved)
	if err != nil {
		log.Printf("Ошибка при сохранении пользователя chat_id=%d, user_name=%s: %v", chatID, userName, err)
		return err
	}

	log.Printf("Пользователь успешно сохранён: chat_id=%d, user_name=%s, is_approved=%v", chatID, userName, isApproved)
	return nil
}

func SaveAccessRequest(conn *pgx.Conn, chatID int64, username string) error {
	ctx := context.Background()

	// SQL-запрос с ON CONFLICT для обновления существующей заявки
	query := `
		INSERT INTO access_requests (user_id, username, status, requested_at, updated_at)
		VALUES ($1, $2, 'pending', NOW(), NOW())
		ON CONFLICT (user_id)
		DO UPDATE SET
			status = EXCLUDED.status,
			updated_at = NOW()
		WHERE access_requests.status != 'pending'
	`

	_, err := conn.Exec(ctx, query, chatID, username)
	if err != nil {
		return fmt.Errorf("ошибка выполнения SQL-запроса: %w", err)
	}

	return nil
}
