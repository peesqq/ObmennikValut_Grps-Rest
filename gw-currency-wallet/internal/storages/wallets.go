package storages

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type WalletStorage struct {
	db *pgxpool.Pool
}

func NewWalletStorage(db *pgxpool.Pool) *WalletStorage {
	return &WalletStorage{db: db}
}

// Получение баланса пользователя
func (s *WalletStorage) GetBalance(ctx context.Context, userID int) (map[string]float64, error) {
	query := `SELECT balance_usd, balance_rub, balance_eur FROM wallets WHERE user_id = $1`
	var balanceUSD, balanceRUB, balanceEUR float64
	err := s.db.QueryRow(ctx, query, userID).Scan(&balanceUSD, &balanceRUB, &balanceEUR)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("wallet not found")
	} else if err != nil {
		return nil, fmt.Errorf("failed to get balance: %v", err)
	}

	return map[string]float64{
		"USD": balanceUSD,
		"RUB": balanceRUB,
		"EUR": balanceEUR,
	}, nil
}

// Пополнение баланса
func (s *WalletStorage) Deposit(ctx context.Context, userID int, currency string, amount float64) error {
	query := fmt.Sprintf(`UPDATE wallets SET balance_%s = balance_%s + $1 WHERE user_id = $2`, currency, currency)
	_, err := s.db.Exec(ctx, query, amount, userID)
	if err != nil {
		return fmt.Errorf("failed to deposit: %v", err)
	}
	return nil
}

// Снятие средств
func (s *WalletStorage) Withdraw(ctx context.Context, userID int, currency string, amount float64) error {
	query := fmt.Sprintf(`UPDATE wallets SET balance_%s = balance_%s - $1 WHERE user_id = $2 AND balance_%s >= $1`, currency, currency, currency)
	_, err := s.db.Exec(ctx, query, amount, userID)
	if err != nil {
		return fmt.Errorf("failed to withdraw: %v", err)
	}
	return nil
}

// Создание кошелька для пользователя
func (s *WalletStorage) CreateWallet(ctx context.Context, userID int) error {
	query := `INSERT INTO wallets (user_id) VALUES ($1)`
	_, err := s.db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to create wallet: %v", err)
	}
	return nil
}
