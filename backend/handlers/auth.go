package handlers

import (
	"encoding/json"
	"net/http"

	"minibooking/db"
	"minibooking/middleware"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int    `json:"expires_in"` // seconds
	UserID       int64  `json:"user_id,omitempty"`
	Username     string `json:"username,omitempty"`
}

// Login authenticates a pre-seeded dummy user and returns JWT tokens.
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	var userID int64
	var password string
	err := db.DB.QueryRow(`SELECT id, password FROM users WHERE username = ?`, req.Username).Scan(&userID, &password)
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, "invalid username or password")
		return
	}
	if password != req.Password {
		writeJSONError(w, http.StatusUnauthorized, "invalid username or password")
		return
	}

	access, err := middleware.GenerateToken(userID, req.Username, "access", middleware.AccessTokenTTL)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}
	refresh, err := middleware.GenerateToken(userID, req.Username, "refresh", middleware.RefreshTokenTTL)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	writeJSON(w, http.StatusOK, tokenResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    int(middleware.AccessTokenTTL.Seconds()),
		UserID:       userID,
		Username:     req.Username,
	})
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// Refresh exchanges a valid refresh token for a new access token.
func Refresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	claims, err := middleware.ParseToken(req.RefreshToken)
	if err != nil || claims.Type != "refresh" {
		writeJSONError(w, http.StatusUnauthorized, "invalid or expired refresh token")
		return
	}
	access, err := middleware.GenerateToken(claims.UserID, claims.Username, "access", middleware.AccessTokenTTL)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}
	writeJSON(w, http.StatusOK, tokenResponse{
		AccessToken: access,
		ExpiresIn:   int(middleware.AccessTokenTTL.Seconds()),
	})
}
