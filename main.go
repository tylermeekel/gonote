package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/go-chi/chi/v5"
)

type App struct {
	db *sql.DB
}

func main() {
	mux := chi.NewMux()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	postgresUri := os.Getenv("GONOTE_POSTGRES_URI")
	if postgresUri == "" {
		log.Fatalln("Could not find Postgres connection URI in environment variable")
	}

	db, err := sql.Open("postgres", postgresUri)
	if err != nil {
		log.Fatalln("Could not connect to Postgres database")
	}
	defer db.Close()

	app := App{db: db}

	mux.Get("/", app.indexHandler)

	fmt.Println("Listening on port "+port)
	log.Fatalln(http.ListenAndServe(":"+port, mux))
}
