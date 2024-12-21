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
	"strings"
	"testing"
)

type atrist_name struct {
	Name string
}

func TestMainPageHandler(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		method         string
		body           string
		expectedStatus int
		expectedName   []atrist_name
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
				divVal, err := extractValueByID(responseHolder.Body.String(), "find")
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

func extractValueByID(html, ID string) ([]string, error) {
	startTag := `<h2 id="` + ID + `">`
	endTag := `</h2>`
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
