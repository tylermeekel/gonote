package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux, port := initializeApp()

	//Start the App
	fmt.Println("Listening on port " + port)
	log.Fatalln(http.ListenAndServe(":"+port, mux))
}
