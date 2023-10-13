package main

import "net/http"

type headerData struct {
	LoggedIn bool
	Title    string
}

func (app *App) handleIndex(w http.ResponseWriter, r *http.Request) {
	var data struct{
		HeaderData headerData
	}

	app.templates.ExecuteTemplate(w, "index", data)
}
