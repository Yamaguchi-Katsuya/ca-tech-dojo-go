package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Yamaguchi-Katsuya/ca-tech-dojo-go/model"
	"github.com/Yamaguchi-Katsuya/ca-tech-dojo-go/service"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		svc: svc,
	}
}

func (u *UserHandler) Create(ctx context.Context, req *model.CreateUserRequest) (*model.CreateUserResponse, error) {
	token, err := u.svc.CreateUser(ctx, req.Name)
	if err != nil {
		return nil, err
	}

	return &model.CreateUserResponse{Token: token}, nil
}

func (u *UserHandler) Get(ctx context.Context, req *model.GetUserRequest) (*model.GetUserResponse, error) {
	name, err := u.svc.GetUser(ctx, req.Token)
	if err != nil {
		return nil, err
	}
	return &model.GetUserResponse{Name: name}, nil
}

func (u *UserHandler) Update(ctx context.Context, req *model.UpdateUserRequest) (*model.UpdateUserResponse, error) {
	err := u.svc.UpdateUser(ctx, req.Token, req.Name)
	if err != nil {
		return nil, err
	}

	return &model.UpdateUserResponse{}, nil
}

func (u *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && r.URL.Path == "/user/create":
		var req model.CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON body", http.StatusBadRequest)
			return
		}
		if req.Name == "" {
			http.Error(w, "name is required", http.StatusBadRequest)
			return
		}
		resp, err := u.Create(r.Context(), &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(resp)
	case r.Method == http.MethodGet && r.URL.Path == "/user/get":
		var req model.GetUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON body", http.StatusBadRequest)
			return
		}
		token := r.Header.Get("x-token")
		if token == "" {
			http.Error(w, "x-token is required", http.StatusBadRequest)
			return
		}
		req.Token = token
		resp, err := u.Get(r.Context(), &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(resp)
	case r.Method == http.MethodPut && r.URL.Path == "/user/update":
		var req model.UpdateUserRequest
		token := r.Header.Get("x-token")
		if token == "" {
			http.Error(w, "x-token is required", http.StatusBadRequest)
			return
		}
		req.Token = token

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON body", http.StatusBadRequest)
			return
		}
		if req.Name == "" {
			http.Error(w, "name is required", http.StatusBadRequest)
			return
		}

		resp, err := u.Update(r.Context(), &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(resp)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
