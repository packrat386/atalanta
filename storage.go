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
	WriteArticle(title string, content []byte) error
	ReadArticle(title string) ([]byte, error)
	ReadArticleVersion(title, version string) ([]byte, error)
	ListArticleVersions(title string) ([]string, error)
	ListArticles() ([]string, error)
}

type localStorage struct {
	baseDirectory string
}

func (l *localStorage) WriteArticle(title string, content []byte) error {
	if !l.exists(title) {
		err := os.Mkdir(l.relpath(title), 0755)
		if err != nil {
			return fmt.Errorf("could not make directory: %w", err)
		}
	}

	fname := l.relpath(title, ts())
	err := os.WriteFile(fname, content, 0644)
	if err != nil {
		return fmt.Errorf("could not write file: %w", err)
	}

	newsym := fname + "_ptr"
	err = os.Symlink(fname, newsym)
	if err != nil {
		return fmt.Errorf("could not symlink: %w", err)
	}

	// is this atomic?
	err = os.Rename(newsym, l.relpath(title, "current"))
	if err != nil {
		return fmt.Errorf("could not make symlink current %w", err)
	}

	// error here doesn't matter?
	os.Remove(newsym)

	return nil
}

func (l *localStorage) ReadArticle(title string) ([]byte, error) {
	return l.ReadArticleVersion(title, "current")
}

func (l *localStorage) ReadArticleVersion(title, version string) ([]byte, error) {
	if !l.exists(title) {
		return nil, errArticleDNE
	}

	data, err := os.ReadFile(l.relpath(title, version))
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}

	return data, nil
}

func (l *localStorage) ListArticleVersions(title string) ([]string, error) {
	if !l.exists(title) {
		return nil, errArticleDNE
	}

	entries, err := os.ReadDir(l.relpath(title))
	if err != nil {
		return nil, fmt.Errorf("could not read directory: %w", err)
	}

	size := len(entries)
	versions := make([]string, size)

	for i, e := range entries {
		if !e.IsDir() {
			versions[(size-1)-i] = e.Name()
		}
	}

	return versions, nil
}

func (l *localStorage) ListArticles() ([]string, error) {
	entries, err := os.ReadDir(l.baseDirectory)
	if err != nil {
		return nil, fmt.Errorf("could not read directory: %w", err)
	}

	titles := []string{}
	for _, e := range entries {
		if e.IsDir() {
			titles = append(titles, e.Name())
		}
	}

	return titles, nil
}

func ts() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

func (l *localStorage) relpath(elem ...string) string {
	return filepath.Join(append([]string{l.baseDirectory}, elem...)...)
}

func (l *localStorage) exists(title string) bool {
	if _, err := os.Stat(l.relpath(title, "current")); err != nil {
		return false
	} else {
		return true
	}
}
