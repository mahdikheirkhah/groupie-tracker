package utils

import (
	"fmt"
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
		tmpl, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/navBar.html")
		if err != nil {
			log.Println("Error parsing template:", err)
			Error(w, "Internal Server Error", "internalServer.html", http.StatusInternalServerError)
			return
		}
		var artists []Artists
		ReadFromAPI("https://groupietrackers.herokuapp.com/api/artists", &artists, w)
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
		log.Printf("Response Status: %d\n", http.StatusOK)
	} else {
		Error(w, "Bad Request Error", "badRequest.html", http.StatusBadRequest)
	}
}

// MoreInfoHandler handles individual artist pages
func MoreInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		Error(w, "Bad Request Error", "badRequest.html", http.StatusBadRequest)
		return
	}
	artistName := strings.TrimPrefix(r.URL.Path, "/artist/")
	URL := "https://groupietrackers.herokuapp.com/api"
	var api API
	ReadFromAPI(URL, &api, w)
	var artists []Artists
	ReadFromAPI(api.Artists, &artists, w)
	IsContain, artistId := Contains(artists, artistName)
	if !IsContain {
		Error(w, "Not Found Error", "notFound.html", http.StatusNotFound)
		return
	}
	var information InformationPage
	information.Artist = artists[artistId-1]
	fmt.Println(information.Artist.CreationDate)
	tmpl, err := template.ParseFiles("templates/MoreInformationPage.html", "templates/header.html", "templates/navBar.html", "templates/goBackButton.html")
	if err != nil {
		log.Println("Error parsing template:", err)
		Error(w, "Internal Server Error", "internalServer.html", http.StatusInternalServerError)
		return
	}
	result := ReadFromAPI(information.Artist.Relations, &information.Relations, w)
	if !result {
		return
	}
	result = ReadFromAPI(information.Artist.ConcertDates, &information.Dates, w)
	for i := 0; i < len(information.Dates.Dates); i++ {
		information.Dates.Dates[i] = strings.TrimPrefix(information.Dates.Dates[i], "*")
	}
	if !result {
		return
	}
	result = ReadFromAPI(information.Artist.Locations, &information.Locations, w)
	if !result {
		return
	}
	err = SafeRenderTemplate(w, tmpl, "MoreInformationPage.html", http.StatusOK, information)
	if err != nil {
		log.Println("Error executing template:", err)
		Error(w, "Internal Server Error", "internalServer.html", http.StatusInternalServerError)
		return
	}
	log.Printf("Response Status: %d\n", http.StatusOK)
}

// Contains checks if an artist name exists in the slice and returns true with the ID, or false and -1 if not found.
func Contains(slice []Artists, lookingName string) (bool, int) {
	lookingName = strings.ToLower(lookingName)
	var lowerCaseArtistsName string
	for _, item := range slice {
		lowerCaseArtistsName = strings.ToLower(item.Name)
		if lowerCaseArtistsName == lookingName {
			return true, item.Id
		}
	}
	return false, -1
}
