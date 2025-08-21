# --- build ---
FROM golang:1.22 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o urlshortener ./cmd/api
# --- run ---
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/urlshortener /usr/local/bin/urlshortener
COPY ./.env /app/.env
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/urlshortener"]