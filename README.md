# URL Shortener (Go + PostgreSQL + React)

## Setup

- Copy `backend/.env.example` to `.env`
- `docker compose up --build`

## Endpoints

- POST /shorten
- GET /{token}
- GET /stats/{token}

## Tests

- `cd backend && go test ./...`

## Design

- Deterministic tokens from normalized URL (base62(sha256)[:8])- GORM with Postgres; AutoMigrate- Fiber for routing & JSO
