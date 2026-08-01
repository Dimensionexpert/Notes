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

type Note struct {
	ID        int
	UserID    int
	Title     string
	Content   string
	CreatedAt string
	UpdatedAt string
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

func CreateNote(db *sql.DB, userID int, title, content string) error {
	_, err := db.Exec(
		"INSERT INTO notes (user_id, title, content) VALUES (?, ?, ?)",
		userID, title, content,
	)
	return err
}

func GetNotesForUser(db *sql.DB, userID int) ([]Note, error) {
	rows, err := db.Query("SELECT id, user_id, title, content, created_at, updated_at FROM notes WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var n Note
		if err := rows.Scan(&n.ID, &n.UserID, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}
	return notes, nil
}
func UpdateNote(db *sql.DB, noteID, userID int, title, content string) error {
	result, err := db.Exec(
		"UPDATE notes SET title = ?, content = ?, updated_at = ? WHERE id = ? AND user_id = ?",
		title, content, time.Now().UTC().Format("2006-01-02 15:04:05"), noteID, userID,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows // note didn't exist, or didn't belong to this user
	}
	return nil
}

func DeleteNote(db *sql.DB, noteID, userID int) error {
	result, err := db.Exec("DELETE FROM notes WHERE id = ? AND user_id = ?", noteID, userID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
