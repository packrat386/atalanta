package main

import (
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/packrat386/atalanta/internal/markdown"
)

func render(w io.Writer, tmpl *template.Template, name string, data interface{}) {
	err := tmpl.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Printf("error executing template '%s': %s", name, err.Error())
	}
}

type errorView struct {
	ErrorMessage string
}

func renderError(w http.ResponseWriter, tmpl *template.Template, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	render(w, tmpl, "error.tmpl", errorView{ErrorMessage: err.Error()})
}

func md2html(input []byte) template.HTML {
	return template.HTML(string(markdown.GenerateHTML(input)))
}
