package controller

import (
	"encoding/json"
	"net/http"
	"rtdocs/model/domain"
	"rtdocs/model/web"
	"rtdocs/service"
)

type AuthController interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}

type authController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) AuthController {
	return &authController{authService: authService}
}

// Register creates a new user
func (c *authController) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdUser, err := c.authService.Register(ctx, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdUser)
}

// Login authenticates a user
func (c *authController) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req web.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	loginResponse, err := c.authService.Login(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loginResponse)
}

// Logout invalidates the access token
func (c *authController) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req web.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := c.authService.Logout(ctx, req.AccessToken); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
