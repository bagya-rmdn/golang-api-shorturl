package db
import (
"log"
"gorm.io/driver/postgres"
"gorm.io/gorm"
"urlshortener/internal/models"
)
func Connect(dsn string) (*gorm.DB, error) {
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
if err != nil {
return nil, err
}
if err := db.AutoMigrate(&models.URLMapping{}); err != nil {
return nil, err
}
log.Println("DB connected & migrated")
return db, nil
}