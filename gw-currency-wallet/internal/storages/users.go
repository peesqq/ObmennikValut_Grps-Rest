package storages

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserStorage struct {
	db *pgxpool.Pool
}

func NewUserStorage(db *pgxpool.Pool) *UserStorage {
	return &UserStorage{db: db}
}

func (s *UserStorage) CreateUser(ctx context.Context, username, email, password string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)`
	_, err = s.db.Exec(ctx, query, username, email, string(passwordHash))
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}
	return nil
}

func (s *UserStorage) AuthenticateUser(ctx context.Context, username, password string) (bool, error) {
	var passwordHash string
	query := `SELECT password_hash FROM users WHERE username = $1`
	err := s.db.QueryRow(ctx, query, username).Scan(&passwordHash)
	if err != nil {
		return false, fmt.Errorf("failed to find user: %v", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return false, nil
	}

	return true, nil
}
