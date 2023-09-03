package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func main() {
	mux := chi.NewMux()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	mux.Get("/", indexHandler)

	fmt.Println("Listening on port " + port)
	log.Fatalln(http.ListenAndServe(":"+port, mux))
}
