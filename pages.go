package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type headerData struct {
	Title      string
	HideHeader bool
}

// handleIndex is a http.HandlerFunc that renders the index page to the ResponseWriter
func (app *App) handleIndex(w http.ResponseWriter, r *http.Request) {
	var data struct {
		HeaderData headerData
	}

	data.HeaderData.Title = "eGoNote"
	data.HeaderData.HideHeader = true

	app.templates.ExecuteTemplate(w, "index", data)
}

// handleLoginPage is a http.HandlerFunc that renders the login page to the ResponseWriter, it will redirect the request if the user is logged in
func (app *App) handleLoginPage(w http.ResponseWriter, r *http.Request) {
	// check if user is logged in
	userID := getUserIDFromContext(r)
	if userID != 0 {
		http.Redirect(w, r, "/notes", http.StatusSeeOther)
		return
	}
	var data struct {
		HeaderData headerData
	}

	data.HeaderData.Title = "Login"
	data.HeaderData.HideHeader = true

	app.templates.ExecuteTemplate(w, "login", data)
}

// handleRegisterPage is a http.HandlerFunc that renders the register page to the ResponseWriter, it will redirect the request if the user is logged in
func (app *App) handleRegisterPage(w http.ResponseWriter, r *http.Request) {
	// check if user is logged in
	userID := getUserIDFromContext(r)
	if userID != 0 {
		http.Redirect(w, r, "/notes", http.StatusSeeOther)
		return
	}
	var data struct {
		HeaderData headerData
	}

	data.HeaderData.Title = "Register"
	data.HeaderData.HideHeader = true

	app.templates.ExecuteTemplate(w, "register", data)
}

// handleNotesPage is a http.Handler that renders the notes page to the ResponseWriter, it will redirect the request if the user is not logged in
func (app *App) handleNotesPage(w http.ResponseWriter, r *http.Request) {
	// confirm that user is logged in
	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	var data struct {
		HeaderData headerData
	}

	data.HeaderData.Title = "Notes"

	app.templates.ExecuteTemplate(w, "notes_page", data)
}

// handleIndividualNotesPage is a http.Handler that renders a specific note's page to the ResponseWriter, it will redirect the request if the user is not logged in
func (app *App) handleIndividualNotePage(w http.ResponseWriter, r *http.Request) {
	// check if user is logged in
	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// get ID value to display correct note
	idString := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
		return
	}

	var data struct {
		HeaderData headerData
		NoteID     int
	}

	data.NoteID = id

	if r.URL.Query().Get("edit") == "true" {
		data.HeaderData.Title = "Editing Note"
		app.templates.ExecuteTemplate(w, "edit_note_page", data)
	} else {
		data.HeaderData.Title = "Viewing Note"
		app.templates.ExecuteTemplate(w, "individual_note_page", data)
	}
}
