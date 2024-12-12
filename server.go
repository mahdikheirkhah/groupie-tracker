package main

import (
	"grpt/utils"
	"log"
	"net/http"
)

// main starts the HTTP server and registers routes
func main() {
	http.Handle("/static/", http.FileServer(http.Dir(".")))
	// Register route handlers
	http.HandleFunc("/", utils.MainPageHandler) // Define the route and its handler function
	http.HandleFunc("/artist/", utils.ConcertsHandler)
	//start the server on port 8080
	log.Println("Starting server on: http://localhost:8080")
	log.Println("Status ok: ", http.StatusOK)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
