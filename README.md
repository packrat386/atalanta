atalanta
--------

A very simple wiki.

## Design

`atalanta` is a simple http service. Users can create and edit articles via their browser. Articles are stored on the local disk along with past versions. Articles are rendered with markdown.

## Installing

Install via `go install`

```
go install github.com/packrat386/atalanta
```

## Running

`atalanta` is configured by environment variables:

* `ATALANTA_BASE_DIR` is the directory to use for storage. Defaults to `.`.
* `ATALANTA_ADDR` is the address to listen on. Defaults to `:http` (port 80).
* `ATALANTA_WIKI_TITLE` is the title for the homepage.
* `ATALANTA_WIKI_BLURB` is the blurb for the homepage.

To run simply run the binary.

```
# for example

ATALANTA_BASE_DIR=~/wikidata ATALANTA_ADDR=':9000' atalanta
```

## Coming Later?

Things I may add one day

* Docker image
* Tests
* Pruning of storage
* Configurable storage
* Users

## Why Make This?

Why not?
