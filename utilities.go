package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/go-chi/chi/v5"
)

type App struct {
	db  *sql.DB
	log *log.Logger
}

// initializeApp initializes the app database and routes and starts the HTTP server on the given port
func initializeApp() (*chi.Mux, string){
	mux := chi.NewMux()

	//Select port from environment variable or default to :3000
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	//Find URI for Postgres connection
	postgresUri := os.Getenv("GONOTE_POSTGRES_URI")
	if postgresUri == "" {
		log.Fatalln("Could not find Postgres connection URI in environment variable")
	}

	//Open DB using postgres driver
	db, err := sql.Open("postgres", postgresUri)
	if err != nil {
		log.Fatalln("Could not connect to Postgres database")
	}
	defer db.Close()

	//Create new app struct to pass db connection
	app := App{
		db:  db,
		log: log.Default(),
	}

	//Mount index handler, and router for notes
	mux.Get("/", app.handleIndex)
	mux.Mount("/notes", app.noteRouter())

	return mux, port
}

// handleIndex renders the index.html file to the ResponseWriter
func (app App) handleIndex(w http.ResponseWriter, r *http.Request) {

	indexTemplate, err := template.ParseFiles("templates/pages/index.html")
	if err != nil {
		fmt.Println("Error processing index template")
		w.Write([]byte("Server error"))
	}
	indexTemplate.Execute(w, nil)
}