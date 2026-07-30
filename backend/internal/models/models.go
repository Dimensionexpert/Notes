package models

import "database/sql"

func CreateUser(db *sql.DB, username, passwordHash string) error {
	_, err := db.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", username, passwordHash)
	return err
}
