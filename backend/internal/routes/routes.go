package routes

import (
"github.com/gofiber/fiber/v2"
"gorm.io/gorm"

"urlshortener/internal/config"
"urlshortener/internal/handlers"
)

func Register(app *fiber.App, db *gorm.DB, cfg *config.Config) {
api := app.Group("/")

api.Post("/shorten", handlers.Shorten(db))
api.Get("/:token", handlers.Redirect(db))
api.Get("/stats/:token", handlers.Stats(db))

app.Get("/health", func(c *fiber.Ctx) error { return c.SendString("ok") })
}