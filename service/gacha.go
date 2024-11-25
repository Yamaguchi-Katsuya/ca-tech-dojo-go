package service

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"math/big"

	"github.com/Yamaguchi-Katsuya/ca-tech-dojo-go/model"
)

type GachaService struct {
	db *sql.DB
}

func NewGachaService(db *sql.DB) *GachaService {
	return &GachaService{db: db}
}

func (g *GachaService) DrawGacha(ctx context.Context, token string, times int) ([]*model.UserCharacter, error) {
	const (
		querySelectUser               = "SELECT id FROM users WHERE token = ?"
		querySelectGachaProbabilities = `
			SELECT
				gp.id,
				gp.character_id,
				gp.probability,
				c.name
			FROM
				gacha_probabilities AS gp
			JOIN
				characters AS c
			ON
				gp.character_id = c.id`
		insertUserCharacter = "INSERT INTO user_characters (user_id, character_id) VALUES (?, ?)"
	)

	user := &model.User{}
	err := g.db.QueryRowContext(ctx, querySelectUser, token).Scan(&user.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &model.UnauthorizedError{}
		}
		return nil, err
	}

	gps := []*model.GachaProbability{}
	rows, err := g.db.QueryContext(ctx, querySelectGachaProbabilities)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		gp := &model.GachaProbability{}
		err := rows.Scan(&gp.ID, &gp.CID, &gp.Probability, &gp.CName)
		if err != nil {
			return nil, err
		}
		gps = append(gps, gp)
	}

	tx, err := g.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	userCharacters := []*model.UserCharacter{}
	for i := 0; i < times; i++ {
		character, err := drawGacha(gps)
		if err != nil {
			return nil, err
		}

		userCharacter := &model.UserCharacter{
			CharacterID: character.ID,
			Name:        character.Name,
		}
		userCharacters = append(userCharacters, userCharacter)

		res, err := tx.ExecContext(ctx, insertUserCharacter, user.ID, character.ID)
		if err != nil {
			return nil, err
		}
		userCharacter.ID, err = res.LastInsertId()
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return userCharacters, nil
}

func drawGacha(gps []*model.GachaProbability) (*model.Character, error) {
	randomValue, err := SecureRandomFloat()
	if err != nil {
		return nil, err
	}

	character := &model.Character{}
	accumulatedProbability := 0.0
	for _, gp := range gps {
		accumulatedProbability += gp.Probability
		if randomValue < accumulatedProbability {
			character = &model.Character{
				ID:   gp.CID,
				Name: gp.CName,
			}
			return character, nil
		}
	}
	return nil, errors.New("failed to draw gacha")
}

func SecureRandomFloat() (float64, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000)) // 0 ~ 999999
	if err != nil {
		return 0, err
	}
	return float64(n.Int64()) / 1000000.0, nil
}
