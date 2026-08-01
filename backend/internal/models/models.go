package models

import (
	"database/sql"
	"errors"
	"time"

	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

var ErrUsernameExists = errors.New("username already exists")

type User struct {
	ID           int
	Username     string
	PasswordHash string
}

func CreateUser(db *sql.DB, username, passwordHash string) error {
	_, err := db.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", username, passwordHash)
	if err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
				return ErrUsernameExists
			}

		}
		return err
	}
	return nil
}

func GetUserByUsername(db *sql.DB, username string) (*User, error) {
	var u User
	row := db.QueryRow("SELECT id, username, password_hash FROM users WHERE username = ?", username)
	err := row.Scan(&u.ID, &u.Username, &u.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func GetSessionByTokenHash(db *sql.DB, tokenHash string) (int, error) {
	var userID int
	row := db.QueryRow(
		"SELECT user_id FROM sessions WHERE token_hash = ? AND expires_at > CURRENT_TIMESTAMP",
		tokenHash,
	)
	err := row.Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func CreateSession(db *sql.DB, userID int, tokenHash string, expiresAt time.Time) error {
	_, err := db.Exec(
		"INSERT INTO sessions (user_id, token_hash, expires_at) VALUES (?, ?, ?)",
		userID, tokenHash, expiresAt.UTC().Format("2006-01-02 15:04:05"),
	)
	return err
}
