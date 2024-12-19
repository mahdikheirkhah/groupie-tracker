package utils

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// ReadFromAPI sends an HTTP request to the specified URL, parses the JSON response, and saves the result.
// Returns false if any step fails, logging the error and sending an appropriate HTTP response.
func ReadFromAPI(URL string, toSaveResult any, w http.ResponseWriter) bool {
	log.Println(URL)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		log.Println("Error Reading From API:", err)
		Error(w, "Internal Server Error", "internalServer.html", http.StatusInternalServerError)
		return false
	}
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error Reading From API:", err)
		Error(w, "Internal Server Error", "internalServer.html", http.StatusInternalServerError)
		return false
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error Reading From API:", err)
		Error(w, "Internal Server Error", "internalServer.html", http.StatusInternalServerError)
		return false
	}
	err = json.Unmarshal(body, &toSaveResult)
	if err != nil {
		log.Println("Error Reading From API:", err)
		Error(w, "Internal Server Error", "internalServer.html", http.StatusInternalServerError)
		return false
	}
	return true
}
