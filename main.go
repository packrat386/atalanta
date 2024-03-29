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

	storageBaseDirectory := os.Getenv("ATALANTA_BASE_DIR")
	if storageBaseDirectory == "" {
		storageBaseDirectory = "."
	}

	storage, err := NewLocalStorage(storageBaseDirectory)
	if err != nil {
		panic(err)
	}

	title := os.Getenv("ATALANTA_WIKI_TITLE")
	if title == "" {
		title = "Atalanta"
	}

	blurb := os.Getenv("ATALANTA_WIKI_BLURB")
	if blurb == "" {
		blurb = "Run free"
	}

	mux := http.NewServeMux()
	mux.Handle("/goto", newGotoHandler())
	mux.Handle("/articles/", newArticleHandler(storage, tmpl))
	mux.Handle("/articles", newArticlesHandler(storage, tmpl))
	mux.Handle("/versions/", newVersionHandler(storage, tmpl))
	mux.Handle("/", newPublicHandler(title, blurb, public, tmpl))

	srv := http.Server{
		Addr:    os.Getenv("ATALANTA_ADDR"),
		Handler: withLogging(mux),
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
