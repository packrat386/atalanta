package main

import (
	"bytes"
	"html/template"
	"io/fs"
	"net/http"
)

func newPublicHandler(title, blurb string, public fs.FS, tmpl *template.Template) http.Handler {
	homebuf := new(bytes.Buffer)

	render(
		homebuf,
		tmpl,
		"home.tmpl",
		struct {
			Title string
			Blurb string
		}{
			Title: title,
			Blurb: blurb,
		},
	)

	fileserver := http.FileServer(http.FS(public))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.Write(homebuf.Bytes())
		} else {
			fileserver.ServeHTTP(w, r)
		}
	})
}
