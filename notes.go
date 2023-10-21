package main

import (
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
func (app *App) noteRouter() http.Handler {
	router := chi.NewRouter()

	//Routes
	router.Get("/", app.handleGetAllNotes)
	router.Get("/{id}", app.handleGetNoteByID)
	router.Post("/", app.handleNewNote)
	router.Post("/{id}", app.handleUpdateNote)

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
func (app *App) getNoteByID(id, userID int) Note {
	var note Note
	row := app.db.QueryRow("SELECT * FROM notes WHERE id = $1 AND user_id = $2", id, userID)
	err := row.Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &note.CreatedAt)
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

func (app *App) updateNote(id, userID int, title, content string) {
	row := app.db.QueryRow("UPDATE notes SET title = $1, content = $2 WHERE id = $3 AND user_id = $4", title, content, id, userID)
	err := row.Scan()
	if err != nil {
		fmt.Println(err.Error())
	}
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

	note := app.getNoteByID(id, userID)

	if note.UserID != userID {
		app.templates.ExecuteTemplate(w, "error_toast", "Note not found")
		return
	}

	if r.URL.Query().Get("edit") == "true" {
		app.templates.ExecuteTemplate(w, "edit_note", note)
	} else {
		safeHTML := mdToHTML(note.Content)
		safeHTMLString := string(safeHTML)
		data := struct{
			Note Note
			NoteTemplate template.HTML
		}{
			Note: note,
			NoteTemplate: template.HTML(safeHTMLString),
		}
		app.templates.ExecuteTemplate(w, "individual_note", data)
	}
}

// handleNewNote calls postNote with a default title and content.
// It then redirects the user to the page to edit the new note
func (app *App) handleNewNote(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	title := "New Note"
	content := "Lorem ipsum..."

	note := app.postNote(userID, title, content)
	redirectURL := fmt.Sprintf("/notes/%d", note.ID)
	w.Header().Add("HX-Redirect", redirectURL)
	w.WriteHeader(http.StatusOK)
}

// handleUpdateNote gathers the title and content fields from the request form data and calls updateNote with them
// It then redirects the user to the notes page
func (app *App) handleUpdateNote(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := getUserIDFromContext(r)
	title := r.FormValue("title")
	content := r.FormValue("content")

	app.updateNote(id, userID, title, content)
	w.Header().Add("HX-Redirect", fmt.Sprintf("/notes/%d", id))
	w.WriteHeader(http.StatusOK)
}
