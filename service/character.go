package service

import (
	"context"
	"database/sql"

	"github.com/Yamaguchi-Katsuya/ca-tech-dojo-go/model"
)

type CharacterService struct {
	db *sql.DB
}

func NewCharacterService(db *sql.DB) *CharacterService {
	return &CharacterService{db: db}
}

func (c *CharacterService) ListCharacters(ctx context.Context, token string) ([]*model.UserCharacter, error) {
	const query = `
		SELECT
			uc.id,
			uc.character_id,
			c.name
		FROM
			users AS u
		JOIN
			user_characters AS uc
		ON
			uc.user_id = u.id
		JOIN
			characters AS c
		ON
			uc.character_id = c.id
		WHERE
			u.token = ?
		ORDER BY
			uc.created_at DESC
	`

	ucs := []*model.UserCharacter{}
	rows, err := c.db.QueryContext(ctx, query, token)
	if err != nil {
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
