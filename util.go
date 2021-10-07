package main

import (
	"fmt"
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

func md2html(input []byte) (template.HTML, error) {
	html, err := markdown.GenerateHTML(input)
	if err != nil {
		return template.HTML(""), fmt.Errorf("error generating html from markdown: %w", err)
	}

	return template.HTML(string(html)), nil
}

func checkmd(input []byte) error {
	_, err := markdown.GenerateHTML(input)
	if err != nil {
		return fmt.Errorf("error parsing markdown: %w", err)
	}

	return nil
}

type loggingResponseWriter struct {
	http.ResponseWriter
	code int
}

func (l *loggingResponseWriter) WriteHeader(code int) {
	l.code = code
	l.ResponseWriter.WriteHeader(code)
}

func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := &loggingResponseWriter{w, 200}

		next.ServeHTTP(lrw, r)

		log.Printf("%s [%d] %s", r.Method, lrw.code, r.URL.String())
	})
}
