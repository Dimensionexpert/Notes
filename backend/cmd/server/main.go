package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Dimensionexpert/notes-app/internal/db"
	"github.com/Dimensionexpert/notes-app/internal/handlers"
	"github.com/Dimensionexpert/notes-app/internal/middleware"
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

	mux := http.NewServeMux()
	mux.HandleFunc("/signup", handlers.SignupHandler(database))
	mux.HandleFunc("/login", handlers.LoginHandler(database))
	mux.Handle("/whoami", middleware.RequireAuth(database, http.HandlerFunc(handlers.WhoAmIHandler)))

	log.Fatal(http.ListenAndServe(":8080", mux))
}
