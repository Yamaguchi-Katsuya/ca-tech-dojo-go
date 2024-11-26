package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Yamaguchi-Katsuya/golang-gacha-app-api/model"
)

type CharacterService struct {
	db *sql.DB
}

func NewCharacterService(db *sql.DB) *CharacterService {
	return &CharacterService{db: db}
}

func (c *CharacterService) ListCharacters(ctx context.Context, token string) ([]*model.UserCharacter, error) {
	const (
		querySelectUser           = "SELECT id FROM users WHERE token = ?"
		querySelectUserCharacters = `
		SELECT
			uc.id,
			uc.character_id,
			c.name
		FROM
			user_characters AS uc
		JOIN
			characters AS c
		ON
			uc.character_id = c.id
		WHERE
			uc.user_id = ?
		ORDER BY
			uc.created_at DESC
	`
	)

	u := &model.User{}
	err := c.db.QueryRowContext(ctx, querySelectUser, token).Scan(&u.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &model.UnauthorizedError{}
		}
		return nil, err
	}

	ucs := []*model.UserCharacter{}
	rows, err := c.db.QueryContext(ctx, querySelectUserCharacters, u.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var uc model.UserCharacter
		if err := rows.Scan(&uc.ID, &uc.CharacterID, &uc.Name); err != nil {
			return nil, err
		}
		ucs = append(ucs, &uc)
	}

	return ucs, nil
}
