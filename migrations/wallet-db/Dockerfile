FROM golang:1.25-alpine AS builder

RUN apk add --no-cache ca-certificates git
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

WORKDIR /migrations

ENTRYPOINT ["goose"]
