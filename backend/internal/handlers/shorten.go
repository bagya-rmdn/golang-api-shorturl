package handlers

import (
"net/http"

"github.com/gofiber/fiber/v2"
"gorm.io/gorm"

"urlshortener/internal/models"
"urlshortener/internal/utils"
)

type ShortenRequest struct {
URL string `json:"url"`
}

type ShortenResponse struct {
Token string `json:"token"`
}

func Shorten(db *gorm.DB) fiber.Handler {
return func(c *fiber.Ctx) error {
var req ShortenRequest
if err := c.BodyParser(&req); err != nil {
return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON"})
}
if req.URL == "" {
return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "url is required"})
}
norm, err := utils.NormalizeURL(req.URL)
if err != nil {
return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"error": "invalid url format"})
}
if norm == "" {
return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "url is required"})
}

token := utils.TokenFromURL(norm)
m := models.URLMapping{LongURL: norm}
// Idempotent create: same LongURL returns same token
if err := db.Where(&models.URLMapping{LongURL: norm}).Attrs(models.URLMapping{Token: token}).FirstOrCreate(&m).Error; err != nil {
return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "db error"})
}
return c.Status(http.StatusCreated).JSON(ShortenResponse{Token: m.Token})
}
}