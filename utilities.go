package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/go-chi/chi/v5"
)

type App struct {
	db  *sql.DB
	log *log.Logger
}

// initializeApp initializes the app database and routes and starts the HTTP server on the given port
func startApp() {
	mux := chi.NewMux()

	//Select port from environment variable or default to :3000
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	//Load .env if one exists
	godotenv.Load()

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
	app := &App{
		db:  db,
		log: log.Default(),
	}

	//Mount index handler, and router for notes
	mux.Get("/", app.handleIndex)
	mux.Mount("/notes", app.noteRouter())
	mux.Mount("/users", app.userRouter())
	mux.Mount("/auth", app.authRouter())
	mux.Get("/freshtoast", app.handleEmptyToast)

	//Add static file server, pattern from Alex Edwards
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static/", fs))

	fmt.Println("Listening on port " + port)
	http.ListenAndServe(":"+port, mux)
}

// handleIndex renders the index.html file to the ResponseWriter
func (app *App) handleIndex(w http.ResponseWriter, r *http.Request) {

	indexTemplate, err := template.ParseFiles("templates/pages/index.html")
	if err != nil {
		fmt.Println("Error processing index template")
		w.Write([]byte("Server error"))
	}
	indexTemplate.Execute(w, nil)
}

func sendToast(w http.ResponseWriter, message string) {
	toastTemplate, err := template.ParseFiles("templates/components/toast.html")
	if err != nil {
		fmt.Println("Error parsing toast.html")
	}

	toastTemplate.Execute(w, message)
}

func sendErrorToast(w http.ResponseWriter, errorMessage string) {
	errorTemplate, err := template.ParseFiles("templates/components/errorToast.html")
	if err != nil {
		fmt.Println("error processing errorToast.html")
	}

	errorTemplate.Execute(w, errorMessage)
}

func (app *App) handleEmptyToast(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<div id=\"toast\"></div>"))
}

func signJWT(username string, expTime time.Time) (string, error){
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expTime),
		Subject:   username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func (app *App) handleTestJWT(w http.ResponseWriter, r *http.Request){

}