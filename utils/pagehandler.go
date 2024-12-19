package utils

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// MainPageHandler serves the main page and handles artist selection redirection
func MainPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		Error(w, "Not Found Error", "notFound.html", http.StatusNotFound)
		return
	}
	if r.Method == http.MethodGet {
		artistId := r.URL.Query().Get("artistId")
		if artistId != "" {
			// Redirect to the appropriate artist URL
			http.Redirect(w, r, "/artist/"+artistId, http.StatusSeeOther)
			return
		}

		// Render the main page
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

// MorInfoHandler handles individual artist pages
func MorInfoHandler(w http.ResponseWriter, r *http.Request) {
	Artistid := strings.TrimPrefix(r.URL.Path, "/artist/")
	ArtistIntId, err := strconv.Atoi(Artistid)
	if err != nil || ArtistIntId < 1 || ArtistIntId > 52 {
		Error(w, "Not Found Error", "notFound.html", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		Error(w, "Bad Request Error", "badRequest.html", http.StatusBadRequest)
		return
	}

	URL := "https://groupietrackers.herokuapp.com/api/artists/" + Artistid
	tmpl, err := template.ParseFiles("templates/MoreInformationPage.html")
	if err != nil {
		log.Println("Error parsing template:", err)
		Error(w, "Internal Server Error", "internalServer.html", http.StatusInternalServerError)
		return
	}
	var information InformationPage
	result := ReadFromAPI(http.MethodGet, URL, &information.Artist, w)
	if !result {
		return
	}
	result = ReadFromAPI(http.MethodGet, information.Artist.Relations, &information.Relations, w)
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
	err = SafeRenderTemplate(w, tmpl, "MoreInformationPage.html", http.StatusOK, information)
	if err != nil {
		log.Println("Error executing template:", err)
		Error(w, "Internal Server Error", "internalServer.html", http.StatusInternalServerError)
		return
	}
}
