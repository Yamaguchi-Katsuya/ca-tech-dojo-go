package service

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"io"
)

type UserService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{db: db}
}

func (u *UserService) CreateUser(ctx context.Context, name string) (string, error) {
	token, err := u.genToken()
	if err != nil {
		return "", err
	}

	const query = "INSERT INTO users (name, token) VALUES (?, ?)"
	_, err = u.db.ExecContext(ctx, query, name, token)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *UserService) GetUser(ctx context.Context, token string) (string, error) {
	const query = "SELECT name FROM users WHERE token = ?"

	var name string
	err := u.db.QueryRowContext(ctx, query, token).Scan(&name)
	if err != nil {
		return "", err
	}

	return name, nil
}

func (u *UserService) UpdateUser(ctx context.Context, token, name string) error {
	const query = "UPDATE users SET name = ? WHERE token = ?"
	updated, err := u.db.ExecContext(ctx, query, name, token)
	if err != nil {
		return err
	}

	rowsAffected, err := updated.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (u *UserService) genToken() (string, error) {
	randBytes := make([]byte, 64)
	_, err := io.ReadFull(rand.Reader, randBytes)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.WithPadding(base64.NoPadding).EncodeToString(randBytes), nil
}
