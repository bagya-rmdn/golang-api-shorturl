package handlers

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"urlshortener/internal/models"
)

func Redirect(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Params("token")
		var m models.URLMapping
		if err := db.Where("token = ?", token).First(&m).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "not found"})
			}
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "db error"})
		}
		// increment analytics
		now := time.Now()
		db.Model(&m).Updates(map[string]any{"clicks": m.Clicks + 1, "last_accessed": &now})

		return c.Redirect(m.LongURL, http.StatusFound)
	}
}
