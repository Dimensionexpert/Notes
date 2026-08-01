package middleware

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/Dimensionexpert/notes-app/internal/auth"
	"github.com/Dimensionexpert/notes-app/internal/models"
)

type contextKey string

const userIDKey contextKey = "userID"

func RequireAuth(db *sql.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		tokenHash := auth.HashToken(cookie.Value)
		userID, err := models.GetSessionByTokenHash(db, tokenHash)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)

	})
}

func GetUserID(r *http.Request) (int, bool) {
	userID, ok := r.Context().Value(userIDKey).(int)
	return userID, ok
}
