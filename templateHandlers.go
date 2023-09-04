package main

import (
	"fmt"
	"html/template"
	"net/http"
)

// Pages
func (app App) indexHandler(w http.ResponseWriter, r *http.Request) {
	notes := app.getAllNotes()

	indexTemplate, err := template.ParseFiles("templates/pages/index.html")
	if err != nil {
		fmt.Println("Error processing index template")
		w.Write([]byte("Server error"))
	}
	indexTemplate.Execute(w, notes)
}

//Components
