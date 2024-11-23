package model

type (
	Character struct {
		ID   int64  `json:"characterID"`
		Name string `json:"name"`
	}

	UserCharacter struct {
		ID          int64  `json:"userCharacterID"`
		UserID      int64  `json:"userID"`
		CharacterID int64  `json:"characterID"`
		Name        string `json:"name"`
	}
)
