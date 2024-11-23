package model

import "time"

type (
	User struct {
		ID        int64
		Name      string
		Token     string
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	CreateUserRequest struct {
		Name string `json:"name"`
	}

	CreateUserResponse struct {
		Token string `json:"token"`
	}

	GetUserRequest struct {
		Token string `json:"token"`
	}

	GetUserResponse struct {
		Name string `json:"name"`
	}

	UpdateUserRequest struct {
		Token string `json:"token"`
		Name  string `json:"name"`
	}

	UpdateUserResponse struct{}
)
