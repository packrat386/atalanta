package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func newArticlesHandler(s storage, tmpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		articlesHandler(w, r, s, tmpl)
	})
}

func articlesHandler(w http.ResponseWriter, r *http.Request, s storage, tmpl *template.Template) {
	titles, err := s.ListArticles()
	if err != nil {
		renderError(w, tmpl, fmt.Errorf("could not list articles: %w", err))
		return
	}

	render(w, tmpl, "list_articles.tmpl", titles)
}
