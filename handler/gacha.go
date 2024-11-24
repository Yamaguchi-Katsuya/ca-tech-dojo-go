package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Yamaguchi-Katsuya/ca-tech-dojo-go/model"
	"github.com/Yamaguchi-Katsuya/ca-tech-dojo-go/service"
)

type GachaHandler struct {
	svc *service.GachaService
}

func NewGachaHandler(svc *service.GachaService) *GachaHandler {
	return &GachaHandler{
		svc: svc,
	}
}

func (g *GachaHandler) Draw(ctx context.Context, req *model.DrawGachaRequest) (*model.DrawGachaResponse, error) {
	characters, err := g.svc.DrawGacha(ctx, req.Token, req.Times)
	if err != nil {
		return nil, err
	}
	return &model.DrawGachaResponse{Characters: characters}, nil
}

func (g *GachaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && r.URL.Path == "/gacha/draw":
		var req model.DrawGachaRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON body", http.StatusBadRequest)
			return
		}
		if req.Times <= 0 {
			http.Error(w, "times is required", http.StatusBadRequest)
			return
		}

		token := r.Header.Get("x-token")
		if token == "" {
			http.Error(w, "x-token is required", http.StatusBadRequest)
			return
		}
		req.Token = token
		resp, err := g.Draw(r.Context(), &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	default:
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
