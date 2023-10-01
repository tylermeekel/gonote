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
	Title   string
	Content string
}

//Functions

// getAllNotes gets all notes from the db connection and returns them as a list of notes
func (app App) getAllNotes() []Note {
	var notes []Note
	rows, err := app.db.Query("SELECT title, content FROM notes;")
	if err != nil {
		fmt.Println("error querying notes")
	}
	for rows.Next() {
		var note Note
		rows.Scan(&note.Title, &note.Content)
		notes = append(notes, note)
	}

	return notes
}

func (app App) getNote(id int) Note {
	var note Note
	query := fmt.Sprintf("SELECT title, content FROM notes WHERE id = %d;", id)
	row, err := app.db.Query(query)
	if err != nil{
		fmt.Println("query error: " + query)
		fmt.Println(err.Error())
	}
	row.Next()
	row.Scan(&note.Title, &note.Content)

	return note
}

//Routes

func (app App) noteRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", app.getAllNotesHandler)
	r.Get("/{id}", app.getNoteHandler)

	return r
}

// getAllNotesHandler gets all notes from the getAllNotes function and responds with them in an HTML template
func (app App) getAllNotesHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/components/notes.html")
	if err != nil{
		w.Write([]byte("server error"))
		fmt.Println("error parsing file")
	}

	notes := app.getAllNotes()
	t.Execute(w, notes)
	
}

func (app App) getNoteHandler(w http.ResponseWriter, r *http.Request)  {
	requestedId := chi.URLParam(r, "id")
	id, err := strconv.Atoi(requestedId)
	if err != nil{
		w.Write([]byte("Requested id is not a number"))
		return
	}
	
	note := app.getNote(id)
	res, err := json.Marshal(note)
	if err != nil{
		fmt.Println(err.Error())
	}

	w.Write(res)
}
