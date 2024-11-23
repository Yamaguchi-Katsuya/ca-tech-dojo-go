package model

type (
	Character struct {
		ID   int64  `json:"characterID"`
		Name string `json:"name"`
	}

	UserCharacter struct {
		ID          int64  `json:"userCharacterID"`
		CharacterID int64  `json:"characterID"`
		Name        string `json:"name"`
	}

	GetUserCharactersRequest struct {
		Token string `json:"token"`
	}

	GetUserCharactersResponse struct {
		Characters []*UserCharacter `json:"characters"`
	}
)
