package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"urlshortener/internal/models"
)

func Stats(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Params("token")
		var m models.URLMapping
		if err := db.Where("token = ?", token).First(&m).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "not found"})
			}
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "db error"})
		}
		return c.Status(http.StatusOK).JSON(m)
	}
}
