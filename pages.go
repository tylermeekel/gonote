package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type headerData struct {
	Title    string
	HideHeader bool
}

func (app *App) handleIndex(w http.ResponseWriter, r *http.Request) {
	var data struct {
		HeaderData headerData
	}

	data.HeaderData.Title = "eGoNote"
	data.HeaderData.HideHeader = true

	app.templates.ExecuteTemplate(w, "index", data)
}

func (app *App) handleLoginPage(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID != 0 {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	var data struct {
		HeaderData headerData
	}

	data.HeaderData.Title = "Login"
	data.HeaderData.HideHeader = true

	app.templates.ExecuteTemplate(w, "login", data)
}

func (app *App) handleRegisterPage(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID != 0 {
		http.Redirect(w, r, "/notes", http.StatusUnauthorized)
		return
	}
	var data struct {
		HeaderData headerData
	}

	data.HeaderData.Title = "Register"
	data.HeaderData.HideHeader = true

	app.templates.ExecuteTemplate(w, "register", data)
}

func (app *App) handleNotesPage(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	var data struct {
		HeaderData headerData
	}

	data.HeaderData.Title = "Notes"

	app.templates.ExecuteTemplate(w, "notes_page", data)
}

func (app *App) handleIndividualNotePage(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0{
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return
	}

	idString := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idString)
	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
		return
	}

	var data struct{
		HeaderData headerData
		NoteID int
	}

	data.HeaderData.Title = fmt.Sprintf("Note #%d", id)
	data.NoteID = id

	app.templates.ExecuteTemplate(w, "individual_note_page", data)
}
