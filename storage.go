package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type storage interface {
	Exists(name string) (bool, error)
	Create(name string) error
	Read(name string) (string, error)
	Write(name, content string) error
}

type localStorage struct {
	basedir string
}

func (l *localStorage) relpath(dir, file string) string {
	return filepath.Join(l.basedir, dir, file)
}

func (l *localStorage) Exists(name string) (bool, error) {
	if _, err := os.Open(l.relpath(name, "current")); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		} else {
			return false, fmt.Errorf("couldn't check file existence: %w", err)
		}
	} else {
		return true, nil
	}
}

func (l *localStorage) Create(name string) error {
	log.Println("here we go?")
	if ok, _ := l.Exists(name); ok {
		return nil
	}

	log.Println("making: ", filepath.Join(l.basedir, name))
	err := os.Mkdir(filepath.Join(l.basedir, name), 0755)
	if err != nil {
		return fmt.Errorf("could not make directory: %w", err)
	}

	fname := l.relpath(name, ts())
	err = os.WriteFile(fname, []byte(""), 0644)
	if err != nil {
		return fmt.Errorf("could not write file: %w", err)
	}

	err = os.Symlink(fname, l.relpath(name, "current"))
	if err != nil {
		return fmt.Errorf("could not symlink: %w", err)
	}

	return nil
}

func (l *localStorage) Read(name string) (string, error) {
	data, err := os.ReadFile(l.relpath(name, "current"))
	if err != nil {
		return "", fmt.Errorf("could not read file: %w", err)
	}

	return string(data), nil
}

func (l *localStorage) Write(name, content string) error {
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

func ts() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}
