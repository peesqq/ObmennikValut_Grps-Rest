package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB(conn *pgxpool.Pool) {
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		password_hash TEXT NOT NULL
	);
	`

	createWalletsTable := `
	CREATE TABLE IF NOT EXISTS wallets (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL REFERENCES users(id),
		balance_usd NUMERIC(12, 2) DEFAULT 0,
		balance_rub NUMERIC(12, 2) DEFAULT 0,
		balance_eur NUMERIC(12, 2) DEFAULT 0
	);
	`

	// Выполняем запросы для создания таблиц
	queries := []string{createUsersTable, createWalletsTable}
	for _, query := range queries {
		_, err := conn.Exec(context.Background(), query)
		if err != nil {
			log.Fatalf("Failed to execute query: %v", err)
		}
	}

	log.Println("Database initialized successfully!")
}
