package main

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"regexp"
)

func newArticleHandler(s storage, tmpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleArticle(w, r, s, tmpl)
	})
}

func handleArticle(w http.ResponseWriter, r *http.Request, s storage, tmpl *template.Template) {
	if r.Method == "POST" {
		postArticle(w, r, s, tmpl)
	} else if r.Method == "GET" {
		getArticle(w, r, s, tmpl)
	} else {
		renderError(w, tmpl, fmt.Errorf("method must be POST or GET"))
	}
}

func postArticle(w http.ResponseWriter, r *http.Request, s storage, tmpl *template.Template) {
	title, err := articleTitle(r)
	if err != nil {
		renderError(w, tmpl, err)
		return
	}

	err = r.ParseForm()
	if err != nil {
		renderError(w, tmpl, fmt.Errorf("could not parse form values: %w", err))
		return
	}

	content := r.Form.Get("content")
	err = s.WriteArticle(title, []byte(content))
	if err != nil {
		renderError(w, tmpl, fmt.Errorf("could not write article content:  %w", err))
		return
	}

	http.Redirect(w, r, r.URL.Path, http.StatusFound)
}

type articleView struct {
	Title   string
	Content string
}

func getArticle(w http.ResponseWriter, r *http.Request, s storage, tmpl *template.Template) {
	title, err := articleTitle(r)
	if err != nil {
		renderError(w, tmpl, err)
		return
	}

	content, err := s.ReadArticle(title)
	if errors.Is(err, errArticleDNE) {
		render(w, tmpl, "article_dne.tmpl", articleView{Title: title})
		return
	} else if err != nil {
		renderError(w, tmpl, fmt.Errorf("could not read article: %w", err))
		return
	}

	a := articleView{
		Title:   title,
		Content: string(content),
	}

	if r.URL.Query().Get("raw") == "true" {
		w.Write([]byte(a.Content))
		return
	}

	if r.URL.Query().Get("edit") == "true" {
		render(w, tmpl, "edit_article.tmpl", a)
		return
	}

	render(w, tmpl, "show_article.tmpl", a)
}

var articlePathMatcher = regexp.MustCompile(`^/articles/([0-9a-zA-Z_]+)$`)

func articleTitle(r *http.Request) (string, error) {
	matches := articlePathMatcher.FindStringSubmatch(r.URL.Path)
	if matches == nil {
		return "", fmt.Errorf("invalid article URL")
	}

	return matches[1], nil
}

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
