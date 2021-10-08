FROM golang:latest as builder

WORKDIR /atalanta
COPY . .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /go/bin/atalanta

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /go/bin/atalanta /usr/local/bin

WORKDIR /atalanta
ENTRYPOINT ["atalanta"]