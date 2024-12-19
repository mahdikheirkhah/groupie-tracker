package utils

import (
	"bytes"
	"html/template"
	"net/http"
)

// SafeRenderTemplate safely renders a template using bytes.Buffer and writes the output to the response.
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
