package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type App struct {
	templates *template.Template
	db  *sql.DB
	log *log.Logger
}

type contextKey string
const userIDKey contextKey = "userID"

// initializeApp initializes the app database and routes and starts the HTTP server on the given port
func startApp() {
	mux := chi.NewMux()

	//Recommended default middleware stack
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	//Custom middleware
	mux.Use(checkAuthentication)

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

	//Parse all templates
	templates := template.Must(template.ParseGlob("templates/*/*.html"))

	//Create new app struct to pass db connection
	app := &App{
		templates: templates,
		db:  db,
		log: log.Default(),
	}

	//Mount routers and utility handlers
	mux.Mount("/", app.frontendRouter())
	mux.Mount("/api", app.apiRouter())
	mux.Get("/freshtoast", app.handleEmptyToast)

	//Add static file server, pattern from Alex Edwards
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static/", fs))

	fmt.Println("Listening on port " + port)
	http.ListenAndServe(":"+port, mux)
}

// sendToast takes a ResponseWriter and message string and sends back a toast
// notification to the client front end using the toast template
func sendToast(w http.ResponseWriter, message string) {
	toastTemplate, err := template.ParseFiles("templates/components/toast.html")
	if err != nil {
		fmt.Println("Error parsing toast.html")
	}

	toastTemplate.Execute(w, message)
}

// sendErrorToast takes a ResponseWriter and message string and sends back am error toast
// notification to the client front end using the toast template
func (app *App) sendErrorToast(w http.ResponseWriter, errorMessage string) {
	app.templates.ExecuteTemplate(w, "error_toast", errorMessage)
}

// handleEmptyToast is called after a toast message has timed out on the front end
// it responds with an empty toast message to serve as a placeholder for the next
// toast message.
func (app *App) handleEmptyToast(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<div id=\"toast\"></div>"))
}

func (app *App) frontendRouter() *chi.Mux{
	router := chi.NewRouter()

	//Pages
	router.Get("/", app.handleIndex)
	router.Get("/login", app.handleLoginPage)
	router.Get("/register", app.handleRegisterPage)
	router.Get("/notes", app.handleNotesPage)
	router.Get("/notes/{id}", app.handleIndividualNotePage)

	return router
}

func (app *App) apiRouter() *chi.Mux{
	router := chi.NewRouter()

	router.Mount("/notes", app.noteRouter())
	router.Mount("/auth", app.authRouter())

	return router
}