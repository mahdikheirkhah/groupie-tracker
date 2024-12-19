package main

import (
	"grpt/utils"
	"log"
	"net/http"
)

// main initializes the HTTP server, registers routes, and starts listening for incoming requests.
func main() {
	http.Handle("/static/", http.FileServer(http.Dir(".")))
	// Register route handlers
	http.HandleFunc("/", utils.MainPageHandler)
	http.HandleFunc("/artist/", utils.MoreInfoHandler)
	//start the server on port 8080
	log.Println("Starting server on: http://localhost:8080")
	log.Println("Status ok: ", http.StatusOK)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
