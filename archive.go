package main

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
)

func newVersionHandler(s storage, tmpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		versionHandler(w, r, s, tmpl)
	})
}

type versionListView struct {
	Title      string
	VersionIDs []string
}

type versionView struct {
	Title     string
	VersionID string
	Content   string
}

func versionHandler(w http.ResponseWriter, r *http.Request, s storage, tmpl *template.Template) {
	title, err := versionTitle(r)
	if err != nil {
		renderError(w, tmpl, err)
		return
	}

	versionID := r.URL.Query().Get("version_id")
	if versionID == "" {
		versions, err := s.ListArticleVersions(title)
		if err != nil {
			renderError(w, tmpl, fmt.Errorf("could not list article versions: %w", err))
			return
		}

		render(
			w,
			tmpl,
			"list_article_versions.tmpl",
			versionListView{
				Title:      title,
				VersionIDs: versions,
			},
		)

		return
	}

	content, err := s.ReadArticleVersion(title, versionID)
	if err != nil {
		renderError(w, tmpl, fmt.Errorf("could not get content of article version: %w", err))
		return
	}

	a := versionView{
		Title:     title,
		VersionID: versionID,
		Content:   content,
	}

	if r.URL.Query().Get("raw") == "true" {
		w.Write([]byte(a.Content))
		return
	}

	render(w, tmpl, "show_article_version.tmpl", a)
}

var versionPathMatcher = regexp.MustCompile(`^/versions/([0-9a-zA-Z_]+)$`)

func versionTitle(r *http.Request) (string, error) {
	matches := versionPathMatcher.FindStringSubmatch(r.URL.Path)
	if matches == nil {
		return "", fmt.Errorf("invalid version URL")
	}

	return matches[1], nil
}
