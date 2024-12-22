package utils

import (
	"html/template"
	"log"
	"net/http"
)

// Error renders an error page with the specified template, error message, and status code.
func Error(w http.ResponseWriter, errorMessage string, htmlFileName string, statusCode int) {
	log.Printf("Response Status: %d\n", statusCode)
	htmlFileAddress := "templates/" + htmlFileName
	tmpl, err := template.ParseFiles(htmlFileAddress, "templates/header.html", "templates/navBar.html", "templates/goBackButton.html")
	if err != nil {
		log.Println(err)
		http.Error(w, errorMessage, statusCode)
		return
	}
	err = SafeRenderTemplate(w, tmpl, htmlFileName, statusCode, nil)
	if err != nil {
		if statusCode == http.StatusInternalServerError {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		} else {
			Error(w, "internal Server Error", "internalServer.html", http.StatusInternalServerError)
		}
	}
}
