package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Dimensionexpert/notes-app/internal/auth"
	"github.com/Dimensionexpert/notes-app/internal/middleware"
	"github.com/Dimensionexpert/notes-app/internal/models"
)

type signupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func SignupHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req signupRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		hashedPassword, err := auth.HashPassword(req.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = models.CreateUser(db, req.Username, hashedPassword)
		if err != nil {
			if errors.Is(err, models.ErrUsernameExists) {
				http.Error(w, "username already exists", http.StatusConflict)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User created successfully"))
	}
}

func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user, err := models.GetUserByUsername(db, req.Username)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "invalid username or password", http.StatusUnauthorized)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := auth.VerifyPassword(user.PasswordHash, req.Password); err != nil {
			http.Error(w, "invalid username or password", http.StatusUnauthorized)
			return
		}
		// next: generate session token, insert into sessions table, set cookie
		token, err := auth.GenToken()
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		tokenhash := auth.HashToken(token)
		expiresAt := time.Now().Add(24 * time.Hour)
		if err := models.CreateSession(db, user.ID, tokenhash, expiresAt); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    token,
			Expires:  expiresAt,
			HttpOnly: true,
			Path:     "/",
		})
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Session created successfuly"))
	}
}

func WhoAmIHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "no user id in context", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "You are user #%d", userID)
}
