package main

import (
	"fmt"
	"log"

	"github.com/Dimensionexpert/notes-app/internal/db"
)

func main() {
	database, err := db.OpenDB("notes.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	log.Println("Connected to database :)")

	if err := db.Migrate(database); err != nil {
		log.Fatal(err)
	}

	fmt.Println("tables created.")
}
