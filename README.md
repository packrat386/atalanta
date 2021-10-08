atalanta
--------

A very simple wiki.

## Design

`atalanta` is a simple http service. Users can create and edit articles via their browser. Articles are stored on the local disk along with past versions. Articles are rendered with markdown.

## Installing

### Docker

Docker images are hosted on GitHub Container Registry [here](https://github.com/packrat386/atalanta/pkgs/container/atalanta)

If you want your storage to persist between container runs you'll need to mount a directory. For example:

```
docker run \
  -p 9000:80 \
  --mount type=bind,source=/var/wikidata,target=/wikidata \
  --env ATALANTA_BASE_DIR=/wikidata \
  ghcr.io/packrat386/atalanta:latest
```

### From Source

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

Logs are sent to standard outut.

## Coming Later?

Things I may add one day

* Tests
* Pruning of storage
* Configurable storage
* Users

## Why Make This?

Why not?
