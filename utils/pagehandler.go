package utils

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

// MainPageHandler serves the main page and handles artist selection redirection
func MainPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		Error(w, "Not Found Error", "notFound.html", http.StatusNotFound)
		return
	}
	if r.Method == http.MethodGet {
		artistName := r.URL.Query().Get("artistName")
		if artistName != "" {
			// Redirect to the appropriate artist URL
			http.Redirect(w, r, "/artist/"+artistName, http.StatusSeeOther)
			return
		}
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			log.Println("Error parsing template:", err)
			Error(w, "Internal Server Error", "internalServer.html", http.StatusInternalServerError)
			return
		}
		var artists []Artists
		ReadFromAPI(http.MethodGet, "https://groupietrackers.herokuapp.com/api/artists", &artists, w)
		if len(artists) == 0 {
			log.Println("No artists data available")
			Error(w, "Internal Server Error", "internalServer.html", http.StatusInternalServerError)
			return
		}

		err = SafeRenderTemplate(w, tmpl, "index.html", http.StatusOK, artists)
		if err != nil {
			log.Println("Error executing template:", err)
			Error(w, "Internal Server Error", "internalServer.html", http.StatusInternalServerError)
			return
		}
	} else {
		Error(w, "Bad Request Error", "badRequest.html", http.StatusBadRequest)
	}
}

// ConcertsHandler handles individual artist pages
func ConcertsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		Error(w, "Bad Request Error", "badRequest.html", http.StatusBadRequest)
		return
	}
	artistName := strings.TrimPrefix(r.URL.Path, "/artist/")
	URL := "https://groupietrackers.herokuapp.com/api/artists"
	var artists []Artists
	ReadFromAPI(http.MethodGet, URL, &artists, w)
	IsContain, artistId := Contains(artists, artistName)
	if !IsContain {
		Error(w, "Not Found Error", "notFound.html", http.StatusNotFound)
		return
	}
	var information InformationPage
	information.Artist = artists[artistId-1]
	tmpl, err := template.ParseFiles("templates/concerts.html")
	if err != nil {
		log.Println("Error parsing template:", err)
		Error(w, "Internal Server Error", "internalServer.html", http.StatusInternalServerError)
		return
	}
	result := ReadFromAPI(http.MethodGet, information.Artist.Relations, &information.Relations, w)
	if !result {
		return
	}
	result = ReadFromAPI(http.MethodGet, information.Artist.ConcertDates, &information.Dates, w)
	for i := 0; i < len(information.Dates.Dates); i++ {
		information.Dates.Dates[i] = strings.TrimPrefix(information.Dates.Dates[i], "*")
	}
	if !result {
		return
	}
	result = ReadFromAPI(http.MethodGet, information.Artist.Locations, &information.Locations, w)
	if !result {
		return
	}
	err = SafeRenderTemplate(w, tmpl, "concerts.html", http.StatusOK, information)
	if err != nil {
		log.Println("Error executing template:", err)
		Error(w, "Internal Server Error", "internalServer.html", http.StatusInternalServerError)
		return
	}
}

func Contains(slice []Artists, str string) (bool, int) {
	for _, item := range slice {
		if item.Name == str {
			return true, item.Id
		}
	}
	return false, -1
}
