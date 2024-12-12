package utils

import (
	"bytes"
	"html/template"
	"net/http"
)

// safeRenderTemplate renders a template safely and writes to the response
func SafeRenderTemplate(w http.ResponseWriter, tmpl *template.Template, htmlFileName string, status int, data any) error {
	var buffer bytes.Buffer
	err := tmpl.ExecuteTemplate(&buffer, htmlFileName, data)
	if err != nil {
		return err
	}
	w.WriteHeader(status)
	buffer.WriteTo(w)
	return nil
}
