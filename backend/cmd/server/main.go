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
	mux.Handle("GET /notes", middleware.RequireAuth(database, http.HandlerFunc(handlers.GetNotesHandler(database))))
	mux.Handle("POST /notes", middleware.RequireAuth(database, http.HandlerFunc(handlers.CreateNoteHandler(database))))
	mux.Handle("PUT /notes/{id}", middleware.RequireAuth(database, http.HandlerFunc(handlers.UpdateNoteHandler(database))))
	mux.Handle("DELETE /notes/{id}", middleware.RequireAuth(database, http.HandlerFunc(handlers.DeleteNoteHandler(database))))

	log.Fatal(http.ListenAndServe(":8080", mux))
}
