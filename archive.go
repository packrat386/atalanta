package main

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
)

func newArchiveHandler(s storage, tmpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		archiveHandler(w, r, s, tmpl)
	})
}

var archivePathMatcher = regexp.MustCompile(`^/archives/([0-9a-zA-Z_]+)$`)

func archiveName(r *http.Request) (string, error) {
	matches := archivePathMatcher.FindStringSubmatch(r.URL.Path)
	if matches == nil {
		return "", fmt.Errorf("invalid archive URL")
	}

	return matches[1], nil
}

type revisionListView struct {
	Title       string
	RevisionIDs []string
}

type archiveView struct {
	Title      string
	RevisionID string
	Content    string
}

func archiveHandler(w http.ResponseWriter, r *http.Request, s storage, tmpl *template.Template) {
	name, err := archiveName(r)
	if err != nil {
		renderError(w, tmpl, err)
		return
	}

	revisionID := r.URL.Query().Get("revision_id")
	if revisionID == "" {
		names, err := s.ListArticleRevisions(name)
		if err != nil {
			renderError(w, tmpl, fmt.Errorf("could not list article revisions: %w", err))
			return
		}

		render(
			w,
			tmpl,
			"list_article_revisions.tmpl",
			revisionListView{
				Title:       name,
				RevisionIDs: names,
			},
		)

		return
	}

	content, err := s.ReadArticleRevision(name, revisionID)
	if err != nil {
		renderError(w, tmpl, fmt.Errorf("could not get content of article revision: %w", err))
		return
	}

	a := archiveView{
		Title:      name,
		RevisionID: revisionID,
		Content:    content,
	}

	if r.URL.Query().Get("raw") == "true" {
		w.Write([]byte(a.Content))
		return
	}

	render(w, tmpl, "show_article_revision.tmpl", a)
}
