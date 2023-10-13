package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Note struct {
	ID        int
	UserID    int
	Title     string
	Content   string
	CreatedAt string
}

//Functions

// noteRouter returns a router with the handlers for the "/notes" path
func (app *App) noteRouter() http.Handler {
	router := chi.NewRouter()

	//Routes
	router.Get("/", app.handleGetAllNotes)
	router.Get("/{id}", app.handleGetNoteByID)
	router.Post("/", app.handlePostNote)

	return router
}

// getAllNotes gets all notes from the db connection and returns them as a list of notes
func (app *App) getAllNotes(userID int) []Note {
	var notes []Note
	rows, err := app.db.Query("SELECT * FROM notes WHERE user_id = $1", userID)
	if err != nil {
		fmt.Println(err.Error())
	}
	for rows.Next() {
		var note Note
		rows.Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &note.CreatedAt)
		notes = append(notes, note)
	}

	return notes
}

// getNoteByID takes an id as an argument and queries the db connection for a Note matching that id
// and then returns a Note object
func (app *App) getNoteByID(userID, id int) Note {
	var note Note
	row, err := app.db.Query("SELECT * FROM notes WHERE id = $1 AND user_id = $2", id, userID)
	if err != nil {
		fmt.Println(err.Error())
	}
	row.Next()
	err = row.Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &note.CreatedAt)
	if err != nil {
		fmt.Println(err)
	}

	return note
}

// postNote takes a user id, title and content as arguments and inserts a new note into the database
// and then returns a Note object
func (app *App) postNote(userID int, title, content string) Note {
	var note Note

	row, err := app.db.Query("INSERT INTO notes(user_id, title, content) VALUES($1, $2, $3) RETURNING *", userID, title, content)
	if err != nil {
		fmt.Println(err.Error())
	}
	if row.Next() {
		err = row.Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &note.CreatedAt)
		if err != nil {
			fmt.Println(err)
		}
	}

	return note
}

//Handlers

// handleGetAllNotes calls the queryAllNotes function and renders the returned notes to the ResponseWriter
func (app *App) handleGetAllNotes(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	notes := app.getAllNotes(userID)
	app.templates.ExecuteTemplate(w, "notes", notes)
}

// handleGetNoteByID calls the queryNoteByID function and renders the returned note to the ResponseWriter
func (app *App) handleGetNoteByID(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	requestedId := chi.URLParam(r, "id")
	id, err := strconv.Atoi(requestedId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Wrapped in note slice so that template is compatible with any number of notes
	note := app.getNoteByID(id, userID)

	if note.UserID != userID {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	app.templates.ExecuteTemplate(w, "notes", note)
}

// handlePostNote collects title and content from form and calls the postNote function with them as arguments.
// It then renders the returned note to the ResponseWriter
func (app *App) handlePostNote(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	title := r.FormValue("title")
	content := r.FormValue("content")

	note := []Note{app.postNote(userID, title, content)}
	app.templates.ExecuteTemplate(w, "notes", note)
}
