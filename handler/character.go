package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Yamaguchi-Katsuya/ca-tech-dojo-go/model"
	"github.com/Yamaguchi-Katsuya/ca-tech-dojo-go/service"
)

type CharacterHandler struct {
	svc *service.CharacterService
}

func NewCharacterHandler(svc *service.CharacterService) *CharacterHandler {
	return &CharacterHandler{
		svc: svc,
	}
}

func (c *CharacterHandler) List(ctx context.Context, req *model.GetUserCharactersRequest) (*model.GetUserCharactersResponse, error) {
	characters, err := c.svc.ListCharacters(ctx, req.Token)
	if err != nil {
		return nil, err
	}
	return &model.GetUserCharactersResponse{Characters: characters}, nil
}

func (c *CharacterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && r.URL.Path == "/character/list":
		var req model.GetUserCharactersRequest

		token := r.Header.Get("x-token")
		if token == "" {
			http.Error(w, "x-token is required", http.StatusBadRequest)
			return
		}
		req.Token = token
		resp, err := c.List(r.Context(), &req)
		if err != nil {
			if errors.Is(err, &model.UnauthorizedError{}) {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
