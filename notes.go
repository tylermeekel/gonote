package main

import (
	"encoding/json"
	"fmt"
	"html/template"
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
func (app App) noteRouter() http.Handler {
	router := chi.NewRouter()
	router.Get("/", app.handleGetAllNotes)
	router.Get("/{id}", app.handleGetNoteByID)

	return router
}

// queryAllNotes gets all notes from the db connection and returns them as a list of notes
func (app App) queryAllNotes() []Note {
	var notes []Note
	rows, err := app.db.Query("SELECT title, content FROM notes;")
	if err != nil {
		fmt.Println(err.Error())
	}
	for rows.Next() {
		var note Note
		rows.Scan(&note.Title, &note.Content)
		notes = append(notes, note)
	}

	return notes
}

// queryOneNoteByID takes an id as an argument and queries the database connection for a Note matching that id, returning a Note object
func (app App) queryOneNoteByID(id int) Note {
	var note Note
	query := fmt.Sprintf("SELECT title, content FROM notes WHERE id = %d;", id)
	row, err := app.db.Query(query)
	if err != nil {
		fmt.Println("query error: " + query)
		fmt.Println(err.Error())
	}
	row.Next()
	row.Scan(&note.Title, &note.Content)

	return note
}

// handleGetAllNotes calls the queryAllNotes function and renders the returned notes to the ResponseWriter
func (app App) handleGetAllNotes(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/components/notes.html")
	if err != nil {
		w.Write([]byte("server error"))
		fmt.Println("error parsing file")
	}

	notes := app.queryAllNotes()
	t.Execute(w, notes)

}

// handleGetNoteByID calls the queryNoteByID function and renders the returned note to the ResponseWriter
func (app App) handleGetNoteByID(w http.ResponseWriter, r *http.Request) {
	requestedId := chi.URLParam(r, "id")
	id, err := strconv.Atoi(requestedId)
	if err != nil {
		w.Write([]byte("Requested id is not a number"))
		return
	}

	note := app.queryOneNoteByID(id)
	res, err := json.Marshal(note)
	if err != nil {
		fmt.Println(err.Error())
	}

	w.Write(res)
}
