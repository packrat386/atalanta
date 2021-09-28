package main

import (
	"fmt"
	"net/http"
)

func newGotoHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		title := r.URL.Query().Get("title")

		if title == "" {
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			http.Redirect(w, r, fmt.Sprintf("/articles/%s", title), http.StatusFound)
		}
	})
}
