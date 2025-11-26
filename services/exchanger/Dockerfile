FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/services/exchanger
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/exchanger ./cmd

FROM alpine:3.20

RUN adduser -D appuser
USER appuser

WORKDIR /app

COPY --from=builder /bin/exchanger /app/exchanger

EXPOSE 50051

ENTRYPOINT ["/app/exchanger"]
