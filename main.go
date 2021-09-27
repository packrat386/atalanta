package main

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

//go:embed templates/*
var templatesFS embed.FS

//go:embed public/*
var publicFS embed.FS

func main() {
	tmpl, err := template.ParseFS(templatesFS, "templates/*.tmpl")
	if err != nil {
		panic(fmt.Errorf("could not parse templates: %w", err))
	}

	public, err := fs.Sub(publicFS, "public")
	if err != nil {
		panic(fmt.Errorf("could not subsystem public assets: %w", err))
	}

	mux := http.NewServeMux()
	mux.Handle("/articles/", newArticleHandler(&localStorage{"/home/acoyle/data"}, tmpl))
	mux.Handle("/", http.FileServer(http.FS(public)))

	srv := http.Server{
		Addr:    os.Getenv("ADDR"),
		Handler: mux,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigterm := make(chan os.Signal, 1)
		signal.Notify(sigterm, syscall.SIGTERM)
		<-sigterm

		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}

		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}

func newArticleHandler(s storage, tmpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleArticle(w, r, s, tmpl)
	})
}

type contentView struct {
	Content string
}

func handleArticle(w http.ResponseWriter, r *http.Request, s storage, tmpl *template.Template) {
	err := s.Create("test")
	if err != nil {
		http.Error(w, fmt.Sprintln("couldn't create: ", err), http.StatusInternalServerError)
	} else {
		log.Println("worked?")
	}

	if r.Method == "POST" {
		log.Println("We got a POST")
		err := r.ParseForm()
		if err != nil {
			panic(err)
		}

		content := r.Form.Get("content")

		fmt.Println(r.Form)
		log.Println(content)

		err = s.Write("test", content)
		if err != nil {
			http.Error(w, fmt.Sprintln("couldn't write: ", err), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/articles/test", http.StatusFound)
	}

	content, err := s.Read("test")
	if err != nil {
		http.Error(w, fmt.Sprintln("couldn't read: ", err), http.StatusInternalServerError)
		return
	}

	c := contentView{Content: content}

	edit := r.URL.Query().Get("edit")

	if edit == "true" {
		err := tmpl.ExecuteTemplate(w, "edit_article.tmpl", c)
		if err != nil {
			log.Println("couldn't exec template: ", err)
		}
	} else {
		err := tmpl.ExecuteTemplate(w, "show_article.tmpl", c)
		if err != nil {
			log.Println("couldn't exec template: ", err)
		}
	}
}
