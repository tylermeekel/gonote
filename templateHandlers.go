package main

import (
	"fmt"
	"html/template"
	"net/http"
)

//Pages
func indexHandler(w http.ResponseWriter, r *http.Request) {
	indexTemplate, err := template.ParseFiles("templates/pages/index.html")
	if err != nil{
		fmt.Println("Error processing index template")
		w.Write([]byte("Server error"))
	}
	indexTemplate.Execute(w, nil)
}

//Components
