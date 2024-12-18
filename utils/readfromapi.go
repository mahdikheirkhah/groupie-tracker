package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func ReadFromAPI(method string, URL string, toSaveResult any, w http.ResponseWriter) bool {
	fmt.Println(URL)
	client := &http.Client{}
	req, err := http.NewRequest(method, URL, nil)
	if err != nil {
		log.Println("Error Reading From API1:", err)
		Error(w, "Internal Server Error", "internalServer.html", http.StatusInternalServerError)
		return false
	}
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error Reading From API2:", err)
		Error(w, "Internal Server Error", "internalServer.html", http.StatusInternalServerError)
		return false
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error Reading From API3:", err)
		Error(w, "Internal Server Error", "internalServer.html", http.StatusInternalServerError)
		return false
	}
	err = json.Unmarshal(body, &toSaveResult)
	if err != nil {
		log.Println("Error Reading From API4:", err)
		Error(w, "Internal Server Error", "internalServer.html", http.StatusInternalServerError)
		return false
	}
	return true
}
