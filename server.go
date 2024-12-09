package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Artists struct {
	Id           int
	Image        string
	Name         string
	Members      []string
	CreationDate int
	FirstAlbum   string
	Locations    string
	ConcertDates string
	Relations    string
}

type Locations struct {
	Id        int
	Locations []string
	Dates     string
}

type Dates struct {
	Id    int
	Dates []string
}

type Relations struct {
	Id             int
	DatesLocations map[string][]string
}

type InformationPage struct {
	Relations Relations
	Locations Locations
	Dates     Dates
	Artist    Artists
}

// main starts the HTTP server and registers routes
func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	// Register route handlers
	http.HandleFunc("/", MainPageHandler) // Define the route and its handler function
	http.HandleFunc("/artist/", ConcertsHandler)
	//start the server on port 8080
	log.Println("Starting server on: http://localhost:8080")
	log.Println("Status ok: ", http.StatusOK)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

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

		for i := 0; i < len(artists); i++ {
			for j := 0; j < len(artists[i].Members); j++ {
				if j != len(artists[i].Members)-1 {
					artists[i].Members[j] += ","
				}
			}
		}

		if len(artists) == 0 {
			log.Println("No artists data available")
			Error(w, "Internal Server Error", "internalServer.html", http.StatusInternalServerError)
			return
		}

		err = safeRenderTemplate(w, tmpl, "index.html", http.StatusOK, artists)
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
	tmpl, err := template.ParseFiles("templates/concerts.html")
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
	fmt.Println(information.Artist.ConcertDates)
	result = ReadFromAPI(http.MethodGet, information.Artist.ConcertDates, &information.Dates, w)
	for i := 0; i < len(information.Dates.Dates); i++ {
		information.Dates.Dates[i] = strings.TrimPrefix(information.Dates.Dates[i], "*")
	}
	if !result {
		fmt.Println("hdhhdhd")
		return
	}
	result = ReadFromAPI(http.MethodGet, information.Artist.Locations, &information.Locations, w)
	if !result {
		return
	}
	err = safeRenderTemplate(w, tmpl, "concerts.html", http.StatusOK, information)
	if err != nil {
		log.Println("Error executing template:", err)
		Error(w, "Internal Server Error", "internalServer.html", http.StatusInternalServerError)
		return
	}
}

func Error(w http.ResponseWriter, errorMessage string, htmlFileName string, statusCode int) {
	log.Printf("Response Status: %d\n", statusCode)
	htmlFileAddress := "templates/" + htmlFileName
	tmpl, err := template.ParseFiles(htmlFileAddress)
	if err != nil {
		http.Error(w, errorMessage, statusCode)
		return
	}
	err = safeRenderTemplate(w, tmpl, htmlFileName, statusCode, nil)
	if err != nil {
		if statusCode == http.StatusInternalServerError {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		} else {
			Error(w, "internal Server Error", "internalServer.html", http.StatusInternalServerError)
		}
	}
}

// safeRenderTemplate renders a template safely and writes to the response
func safeRenderTemplate(w http.ResponseWriter, tmpl *template.Template, htmlFileName string, status int, data any) error {
	var buffer bytes.Buffer
	err := tmpl.ExecuteTemplate(&buffer, htmlFileName, data)
	if err != nil {
		return err
	}
	w.WriteHeader(status)
	buffer.WriteTo(w)
	return nil
}

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
		fmt.Println()
		log.Println("Error Reading From API4:", err)
		Error(w, "Internal Server Error", "internalServer.html", http.StatusInternalServerError)
		return false
	}
	return true
}
