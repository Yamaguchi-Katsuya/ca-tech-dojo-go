package service

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"math/big"

	"github.com/Yamaguchi-Katsuya/golang-gacha-app-api/model"
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

	// ユーザーを取得
	user := &model.User{}
	err := g.db.QueryRowContext(ctx, querySelectUser, token).Scan(&user.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &model.UnauthorizedError{}
		}
		return nil, err
	}

	// ガチャ確率データを取得
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

	// 並列処理用のチャンネルとエラー収集
	resultCh := make(chan *model.UserCharacter, times)
	errCh := make(chan error, times)
	defer close(resultCh)
	defer close(errCh)

	// 並列でガチャを引く
	for i := 0; i < times; i++ {
		go func() {
			character, err := drawGacha(gps)
			if err != nil {
				errCh <- err
				return
			}
			resultCh <- &model.UserCharacter{
				CharacterID: character.ID,
				Name:        character.Name,
			}
		}()
	}

	// 結果を収集
	userCharacters := []*model.UserCharacter{}
	for i := 0; i < times; i++ {
		select {
		case uc := <-resultCh:
			userCharacters = append(userCharacters, uc)
		case err := <-errCh:
			return nil, err
		}
	}

	// データベースにまとめて挿入
	tx, err := g.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	for _, uc := range userCharacters {
		res, err := tx.ExecContext(ctx, insertUserCharacter, user.ID, uc.CharacterID)
		if err != nil {
			return nil, err
		}
		uc.ID, err = res.LastInsertId()
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
	randomValue, err := secureRandomFloat()
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

func secureRandomFloat() (float64, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return 0, err
	}
	return float64(n.Int64()) / 1000000.0, nil
}
