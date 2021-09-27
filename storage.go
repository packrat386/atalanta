package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var (
	errArticleDNE = errors.New("error: article does not exist")
)

type storage interface {
	WriteArticle(name, content string) error
	ReadArticle(name string) (string, error)
	ReadArticleRevision(name, rev string) (string, error)
	ListArticleRevisions(name string) ([]string, error)
	ListArticles() ([]string, error)
}

type localStorage struct {
	baseDirectory string
}

func (l *localStorage) WriteArticle(name, content string) error {
	if !l.exists(name) {
		err := os.Mkdir(l.relpath(name), 0755)
		if err != nil {
			return fmt.Errorf("could not make directory: %w", err)
		}
	}

	fname := l.relpath(name, ts())
	err := os.WriteFile(fname, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("could not write file: %w", err)
	}

	newsym := fname + "_ptr"
	err = os.Symlink(fname, newsym)
	if err != nil {
		return fmt.Errorf("could not symlink: %w", err)
	}

	// is this atomic?
	err = os.Rename(newsym, l.relpath(name, "current"))
	if err != nil {
		return fmt.Errorf("could not make symlink current %w", err)
	}

	// error here doesn't matter?
	os.Remove(newsym)

	return nil
}

func (l *localStorage) ReadArticle(name string) (string, error) {
	return l.ReadArticleRevision(name, "current")
}

func (l *localStorage) ReadArticleRevision(name, revision string) (string, error) {
	if !l.exists(name) {
		return "", errArticleDNE
	}

	data, err := os.ReadFile(l.relpath(name, revision))
	if err != nil {
		return "", fmt.Errorf("could not read file: %w", err)
	}

	return string(data), nil
}

func (l *localStorage) ListArticleRevisions(name string) ([]string, error) {
	if !l.exists(name) {
		return nil, errArticleDNE
	}

	entries, err := os.ReadDir(l.relpath(name))
	if err != nil {
		return nil, fmt.Errorf("could not read directory: %w", err)
	}

	names := []string{}
	for _, e := range entries {
		if !e.IsDir() {
			names = append(names, e.Name())
		}
	}

	return names, nil
}

func (l *localStorage) ListArticles() ([]string, error) {
	entries, err := os.ReadDir(l.baseDirectory)
	if err != nil {
		return nil, fmt.Errorf("could not read directory: %w", err)
	}

	names := []string{}
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}

	return names, nil
}

func ts() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

func (l *localStorage) relpath(elem ...string) string {
	return filepath.Join(append([]string{l.baseDirectory}, elem...)...)
}

func (l *localStorage) exists(name string) bool {
	if _, err := os.Stat(l.relpath(name, "current")); err != nil {
		return false
	} else {
		return true
	}
}
