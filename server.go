package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
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

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Println("Error parsing template:", err)
		internalServerError(w)
		return
	}
	var body []byte
	body, err = ReadFromAPI(http.MethodGet, "https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		log.Println(err.Error())
		internalServerError(w)
		return
	}
	// Unmarshal JSON into Go struct
	var artists []Artists
	// var relations []Relations
	// var relation Relations
	err = json.Unmarshal(body, &artists)
	if err != nil {
		log.Println("Error unmarshalling JSON:", err)
		internalServerError(w)
		return
	}
	for i := 0; i < len(artists); i++ {
		for j := 0; j < len(artists[i].Members); j++ {
			if j != len(artists[i].Members)-1 {
				artists[i].Members[j] += ","
			}
		}
		// body, err = ReadFromAPI(http.MethodGet, artists[i].Relations)
		// err = json.Unmarshal(body, &relation)
		// if err != nil {
		// 	log.Println("Error unmarshalling JSON:", err)
		// 	log.Println("we are here")
		// 	internalServerError(w)
		// 	return
		// }
		// relations = append(relations, relation)
	}
	// for i := 0; i < len(relations); i++ {
	// 	fmt.Println(relations[i])
	// }
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
		internalServerError(w)
		return
	}
}

func ConcertsHandler(w http.ResponseWriter, r *http.Request) {
	Artistid := strings.TrimPrefix(r.URL.Path, "/artist/")
	URL := "https://groupietrackers.herokuapp.com/api/relation/" + Artistid
	fmt.Println(URL)
	tmpl, err := template.ParseFiles("templates/concerts.html")
	if err != nil {
		log.Println("Error parsing template:", err)
		internalServerError(w)
		return
	}
	var body []byte
	var relation Relations
	body, err = ReadFromAPI(http.MethodGet, URL)
	err = json.Unmarshal(body, &relation)
	if err != nil {
		log.Println("Error unmarshalling JSON:", err)
		log.Println("we are here")
		internalServerError(w)
		return
	}
	fmt.Println(relation)
	err = tmpl.ExecuteTemplate(w, "concerts.html", relation.DatesLocations)
	if err != nil {
		log.Println("Error executing template:", err)
		internalServerError(w)
		return
	}
}

// badRequestError serves a 400 error page
func badRequestError(w http.ResponseWriter) {
	log.Printf("Response Status: %d\n", http.StatusBadRequest)
	tmpl, err := template.ParseFiles("templates/badRequest.html")
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	err = safeRenderTemplate(w, tmpl, "badRequest.html", http.StatusBadRequest, nil)
	if err != nil {
		internalServerError(w)
		return
	}
}

// notFoundError serves a 404 error page
func notFoundError(w http.ResponseWriter) {
	log.Printf("Response Status: %d\n", http.StatusNotFound)
	tmpl, err := template.ParseFiles("templates/notFound.html")
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	err = safeRenderTemplate(w, tmpl, "notFound.html", http.StatusNotFound, nil)
	if err != nil {
		internalServerError(w)
		return
	}
}

// internalServerError serves a 500 error page
func internalServerError(w http.ResponseWriter) {
	log.Printf("Response Status: %d\n", http.StatusInternalServerError)
	tmpl, err := template.ParseFiles("templates/internalServer.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = safeRenderTemplate(w, tmpl, "internalServer.html", http.StatusInternalServerError, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
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
