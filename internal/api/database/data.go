package database

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/Pavel-Sergeev-ekb/JARVIS_tg-Bot/internal/config"
	"github.com/jackc/pgx/v4"
)

var db *pgx.Conn

func ConnectDB(cfg config.Config) (*pgx.Conn, error) {

	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		convertPort(cfg.DBPort),
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName)

	var err error

	db, err = pgx.Connect(context.Background(), connString)

	if err != nil {
		log.Fatal("Не удалось подключиться к бд", err)
		return nil, fmt.Errorf("ошибка при подключении к БД: %w", err)
	}
	log.Println("Успешное подключение к базе данных!")
	return db, nil
}
func convertPort(port string) int {
	p, err := strconv.Atoi(port)
	if err != nil {
		return 5432 // порт по умолчанию
	}
	return p
}

func CloseDB(db *pgx.Conn) {
	if db != nil {
		db.Close(context.Background())
		log.Println("соединение с бд закрыто")
	}
}
