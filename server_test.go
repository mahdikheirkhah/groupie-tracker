package main

import (
	"encoding/json"
	"errors"
	"grpt/utils"
	"html"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestMainPageHandler(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		method         string
		body           string
		expectedStatus int
		expectedName   []utils.Artists
	}{
		{
			name:           "Valid Root Path",
			url:            "/",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid answer",
			url:            "/",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Non-existent Path",
			url:            "/non-existent",
			method:         http.MethodGet,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "bad Request",
			url:            "/",
			method:         http.MethodPost,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Internal Server Error",
			url:            "/",
			method:         http.MethodGet,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(test.method, test.url, nil)
			if err != nil {
				t.Fatalf("Could not create request: %v", err)
			}
			if test.name == "Internal Server Error" { //cahnge the name of the file to have internal server Error
				sourceName := "./templates/index.html"
				destinationName := "./templates/home.html"
				err := os.Rename(sourceName, destinationName)
				if err != nil {
					t.Errorf("Error renaming file: %v\n", err)
					return
				}
			} else if test.name == "Valid answer" {
				apiRes, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
				if err != nil {
					t.Errorf("Error reading from api of test file: %v\n", err)
					return
				}
				body, err := io.ReadAll(apiRes.Body)
				if err != nil {
					t.Errorf("Error reading from ReadAll of test file: %v\n", err)
					return
				}

				err = json.Unmarshal(body, &test.expectedName)
				if err != nil {
					t.Errorf("Error reading from Unmarshal of test file: %v\n", err)
					return
				}
			}
			responseHolder := httptest.NewRecorder()
			handler := http.HandlerFunc(utils.MainPageHandler)
			handler.ServeHTTP(responseHolder, req)

			if responseHolder.Code != test.expectedStatus {
				t.Errorf("Expected status code %d, but got %d", test.expectedStatus, responseHolder.Code)
			}
			if strings.HasPrefix(test.name, "Valid answer") { //only for the test cases with name <<valid input>> check for the output
				divVal, err := extractValueByID(responseHolder.Body.String(), "find", "h2", "</h2>")
				if err != nil {
					t.Errorf("Error extracting div: %v", err)
				}
				if len(divVal) != len(test.expectedName) {
					t.Errorf("Expected to have %d artist, but got %d artist", len(test.expectedName), len(divVal))
					return
				}
				for i, name := range test.expectedName {
					divVal[i] = html.UnescapeString(divVal[i])
					if name.Name != divVal[i] {
						t.Errorf("Expected h2 #find to have value %s, but got %s", name.Name, divVal[i])
					}
				}
			}
			if test.name == "Internal Server Error" { //recahnge the name of the file
				sourceName := "./templates/home.html"
				destinationName := "./templates/index.html"
				err := os.Rename(sourceName, destinationName)
				if err != nil {
					t.Errorf("Error renaming file: %v\n", err)
					return
				}
			}
		})
	}
}

func TestMoreInfoHandler(t *testing.T) {
	tests := []struct {
		name               string
		url                string
		method             string
		body               string
		expectedStatus     int
		expectedArtistInfo utils.InformationPage
	}{
		{
			name:           "Valid Root Path",
			url:            "/artist/queen",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid answer m",
			url:            "/artist/Queen",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid answer f",
			url:            "/artist/Gorillaz",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid answer l",
			url:            "/artist/Travis Scott",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid answer m",
			url:            "/artist/Foo Fighters",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Non-existent Path",
			url:            "/artist/non-existent",
			method:         http.MethodGet,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "bad Request",
			url:            "/artist/queen",
			method:         http.MethodPost,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Internal Server Error",
			url:            "/artist/queen",
			method:         http.MethodGet,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(test.method, test.url, nil)
			if err != nil {
				t.Fatalf("Could not create request: %v", err)
			}
			if test.name == "Internal Server Error" { //cahnge the name of the file to have internal server Error
				sourceName := "./templates/MoreInformationPage.html"
				destinationName := "./templates/home.html"
				err := os.Rename(sourceName, destinationName)
				if err != nil {
					t.Errorf("Error renaming file: %v\n", err)
					return
				}
			} else if strings.HasPrefix(test.name, "Valid answer") {
				var artists []utils.Artists

				result := getInfo("https://groupietrackers.herokuapp.com/api/artists", &artists)
				if !result {
					t.Errorf("Error reading from api of test file: %v\n", err)
					return
				}

				artistName := strings.TrimPrefix(test.url, "/artist/")
				IsContain, artistId := utils.Contains(artists, artistName)
				if !IsContain {
					t.Errorf("there is no valid artist")
					return
				}
				test.expectedArtistInfo.Artist = artists[artistId-1]

				switch strings.TrimPrefix(test.name, "Valid answer ") {
				case "d":
					result := getInfo(test.expectedArtistInfo.Artist.ConcertDates, &test.expectedArtistInfo.Dates)
					for i := 0; i < len(test.expectedArtistInfo.Dates.Dates); i++ {
						test.expectedArtistInfo.Dates.Dates[i] = strings.TrimPrefix(test.expectedArtistInfo.Dates.Dates[i], "*")
					}
					if !result {
						t.Errorf("Error reading from api of test file: %v\n", err)
						return
					}
				case "l":
					result := getInfo(test.expectedArtistInfo.Artist.Locations, &test.expectedArtistInfo.Locations)
					if !result {
						t.Errorf("Error reading from api of test file: %v\n", err)
						return
					}
				case "r":
					result := getInfo(test.expectedArtistInfo.Artist.Relations, &test.expectedArtistInfo.Relations)
					if !result {
						t.Errorf("Error reading from api of test file: %v\n", err)
						return
					}
				}
			}
			responseHolder := httptest.NewRecorder()
			handler := http.HandlerFunc(utils.MoreInfoHandler)
			handler.ServeHTTP(responseHolder, req)

			if responseHolder.Code != test.expectedStatus {
				t.Errorf("Expected status code %d, but got %d", test.expectedStatus, responseHolder.Code)
			}
			if strings.HasPrefix(test.name, "Valid answer") { //only for the test cases with name <<valid input>> check for the output
				switch strings.TrimPrefix(test.name, "Valid answer ") {
				case "d":
					divVal, err := extractValueByID(responseHolder.Body.String(), "dates", "li", "</li>")
					if err != nil {
						t.Errorf("Error extracting div: %v", err)
					}
					if len(divVal) != len(test.expectedArtistInfo.Dates.Dates) {
						t.Errorf("Expected to have %d artist, but got %d artist", len(test.expectedArtistInfo.Dates.Dates), len(divVal))
						return
					}
					for i, dates := range test.expectedArtistInfo.Dates.Dates {
						divVal[i] = html.UnescapeString(divVal[i])
						if dates != divVal[i] {
							t.Errorf("Expected h2 #find to have value %s, but got %s", dates, divVal[i])
						}
					}
				case "l":
					divVal, err := extractValueByID(responseHolder.Body.String(), "locations", "li", "</li>")
					if err != nil {
						t.Errorf("Error extracting div: %v", err)
					}
					if len(divVal) != len(test.expectedArtistInfo.Locations.Locations) {
						t.Errorf("Expected to have %d artist, but got %d artist", len(test.expectedArtistInfo.Dates.Dates), len(divVal))
						return
					}
					for i, locations := range test.expectedArtistInfo.Locations.Locations {
						divVal[i] = html.UnescapeString(divVal[i])
						if locations != divVal[i] {
							t.Errorf("Expected to have value %s, but got %s", locations, divVal[i])
						}
					}
				case "m":
					divVal, err := extractValueByID(responseHolder.Body.String(), "members", "li", "</li>")
					if err != nil {
						t.Errorf("Error extracting div: %v", err)
					}
					if len(divVal) != len(test.expectedArtistInfo.Artist.Members) {
						t.Errorf("Expected to have %d artist, but got %d artist", len(test.expectedArtistInfo.Dates.Dates), len(divVal))
						return
					}
					for i, member := range test.expectedArtistInfo.Artist.Members {
						divVal[i] = html.UnescapeString(divVal[i])
						if member != divVal[i] {
							t.Errorf("Expected to have value %s, but got %s", member, divVal[i])
						}
					}
				case "f":
					divVal, err := extractValueByID(responseHolder.Body.String(), "firstAlbum", "li", "</li>")
					if err != nil {
						t.Errorf("Error extracting div: %v", err)
					}
					if len(divVal) != 1 {
						t.Errorf("Expected to have one date for first album, but got %d", len(divVal))
						return
					}
					divVal[0] = html.UnescapeString(divVal[0])
					if test.expectedArtistInfo.Artist.FirstAlbum != divVal[0] {
						t.Errorf("Expected to have value %s, but got %s", test.expectedArtistInfo.Artist.FirstAlbum, divVal[0])
					}
				case "c":
					divVal, err := extractValueByID(responseHolder.Body.String(), "creationDate", "li", "</li>")
					if err != nil {
						t.Errorf("Error extracting div: %v", err)
					}
					if len(divVal) != 1 {
						t.Errorf("Expected to have one date for first album, but got %d", len(divVal))
						return
					}
					divVal[0] = html.UnescapeString(divVal[0])
					number, err := strconv.Atoi(divVal[0])
					if err != nil {
						t.Errorf("Error: %v", err)
					}
					if test.expectedArtistInfo.Artist.CreationDate != number {
						t.Errorf("Expected to have value %s, but got %s", test.expectedArtistInfo.Artist.FirstAlbum, divVal[0])
					}
				}
			}
			if test.name == "Internal Server Error" { //recahnge the name of the file
				sourceName := "./templates/home.html"
				destinationName := "./templates/MoreInformationPage.html"
				err := os.Rename(sourceName, destinationName)
				if err != nil {
					t.Errorf("Error renaming file: %v\n", err)
					return
				}
			}
		})
	}
}

func extractValueByID(html, ID, startTagname, endTag string) ([]string, error) {
	startTag := `<` + startTagname + ` id="` + ID + `">`
	var err error
	notFind := false
	var names []string
	// Find the start index of the desired <h2>
	for i := 0; ; {
		startIndex := strings.Index(html, startTag)
		if startIndex == -1 {
			if i == 0 {
				err = errors.New("tag h2 with id '" + ID + "' not found")
				notFind = true
			}
			break
		}
		startIndex += len(startTag)

		// Find the end index of the </h2>
		endIndex := strings.Index(html[startIndex:], endTag)
		if endIndex == -1 {
			if i == 0 {
				err = errors.New("tag h2 with id '" + ID + "' not found")
				notFind = true
			}
			break
		}
		// Extract and return the content
		names = append(names, html[startIndex:startIndex+endIndex])
		html = html[startIndex+endIndex+1:]
		i++
	}
	if notFind {
		return nil, err
	}
	return names, nil
}

func getInfo(URL string, toSaveResult any) bool {
	apiRes, err := http.Get(URL)
	if err != nil {
		return false
	}

	body, err := io.ReadAll(apiRes.Body)
	if err != nil {
		return false
	}
	err = json.Unmarshal(body, &toSaveResult)
	return err == nil
}
