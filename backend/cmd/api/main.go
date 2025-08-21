package main

import (
"log"
"github.com/gofiber/fiber/v2"
"github.com/gofiber/fiber/v2/middleware/cors"

"urlshortener/internal/config"
"urlshortener/internal/db"
"urlshortener/internal/routes"
)

func main() {
cfg := config.Load()

dbConn, err := db.Connect(cfg.DatabaseURL)
if err != nil {
log.Fatalf("DB connect error: %v", err)
}

app := fiber.New()
app.Use(cors.New())

routes.Register(app, dbConn, cfg)

if err := app.Listen(":" + cfg.Port); err != nil {
log.Fatalf("app.Listen error: %v", err)
}
}