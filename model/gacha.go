package model

type (
	GachaProbability struct {
		ID          int64   `json:"gachaProbabilityID"`
		CID         int64   `json:"characterID"`
		Probability float64 `json:"probability"`
		CName       string  `json:"characterName"`
	}

	DrawGachaRequest struct {
		Token string `json:"token"`
		Times int    `json:"times"`
	}

	DrawGachaResponse struct {
		Characters []*UserCharacter `json:"characters"`
	}
)
