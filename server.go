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

// PageVariables holds data to be passed to HTML templates
type PageVariables struct {
	Response       string
	Input          string
	SelectedBanner string
	SpecialTrigger bool
}

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
	Dates string
}

type Relations struct {
	Id             int
	DatesLocations map[string][]string
}

type Indexs struct {
	Location  []Locations
	Dates     []Dates
	Relations []Relations
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

// PostHandler handles POST requests to generate ASCII art
func MainPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		Error(w, "Not Found Error", "notFound.html", 404)
		return
	}
	if r.Method != http.MethodGet {
		Error(w, "Bad Request Error", "badRequest.html", 400)
		return
	}
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Println("Error parsing template:", err)
		Error(w, "internal Server Error", "internalServer.html", 500)
		return
	}
	var body []byte
	body, err = ReadFromAPI(http.MethodGet, "https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		log.Println(err.Error())
		Error(w, "internal Server Error", "internalServer.html", 500)
		return
	}
	// Unmarshal JSON into Go struct
	var artists []Artists
	// var relations []Relations
	// var relation Relations
	err = json.Unmarshal(body, &artists)
	if err != nil {
		log.Println("Error unmarshalling JSON:", err)
		Error(w, "internal Server Error", "internalServer.html", 500)
		return
	}
	for i := 0; i < len(artists); i++ {
		for j := 0; j < len(artists[i].Members); j++ {
			if j != len(artists[i].Members)-1 {
				artists[i].Members[j] += ","
			}
		}
	}
	// Check if there is data
	if len(artists) == 0 {
		log.Println("No artists data available")
		http.Error(w, "No data available", http.StatusNotFound)
		return
	}

	// Render the first artist for demonstration
	err = tmpl.ExecuteTemplate(w, "index.html", artists)
	if err != nil {
		log.Println("Error executing template:", err)
		Error(w, "internal Server Error", "internalServer.html", 500)
		return
	}
}

func ConcertsHandler(w http.ResponseWriter, r *http.Request) {
	Artistid := strings.TrimPrefix(r.URL.Path, "/artist/")
	ArtistIntId, err := strconv.Atoi(Artistid)
	if err != nil || ArtistIntId < 1 || ArtistIntId > 52 {
		Error(w, "Not Found Error", "notFound.html", 404)
		return
	}
	if r.URL.Path != ("/artist/" + Artistid) {
		Error(w, "Not Found Error", "notFound.html", 404)
		return
	}
	if r.Method != http.MethodGet {
		Error(w, "Bad Request Error", "badRequest.html", 400)
		return
	}
	URL := "https://groupietrackers.herokuapp.com/api/relation/" + Artistid
	fmt.Println(URL)
	tmpl, err := template.ParseFiles("templates/concerts.html")
	if err != nil {
		log.Println("Error parsing template:", err)
		Error(w, "internal Server Error", "internalServer.html", 500)
		return
	}
	var body []byte
	var relation Relations
	body, err = ReadFromAPI(http.MethodGet, URL)
	if err != nil {
		log.Println("Error Reading From API:", err)
		Error(w, "internal Server Error", "internalServer.html", 500)
		return
	}
	err = json.Unmarshal(body, &relation)
	if err != nil {
		log.Println("Error unmarshalling JSON:", err)
		log.Println("we are here")
		Error(w, "internal Server Error", "internalServer.html", 500)
		return
	}
	fmt.Println(relation)
	err = tmpl.ExecuteTemplate(w, "concerts.html", relation.DatesLocations)
	if err != nil {
		log.Println("Error executing template:", err)
		Error(w, "internal Server Error", "internalServer.html", 500)
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
		if statusCode == 500 {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		} else {
			Error(w, "internal Server Error", "internalServer.html", 500)
		}
	}
}

// safeRenderTemplate renders a template safely and writes to the response
func safeRenderTemplate(w http.ResponseWriter, tmpl *template.Template, templateName string, status int, data any) error {
	var buffer bytes.Buffer
	err := tmpl.ExecuteTemplate(&buffer, templateName, data)
	if err != nil {
		return err
	}
	w.WriteHeader(status)
	buffer.WriteTo(w)
	return nil
}

func ReadFromAPI(method string, URL string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, URL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
