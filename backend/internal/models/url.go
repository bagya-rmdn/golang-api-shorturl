package models

import "time"

type URLMapping struct {
ID uint `gorm:"primaryKey" json:"id"`
Token string `gorm:"uniqueIndex;size:12" json:"token"`
LongURL string `gorm:"uniqueIndex;size:2048" json:"long_url"`
Clicks int64 `json:"clicks"`
CreatedAt time.Time `json:"created_at"`
UpdatedAt time.Time `json:"updated_at"`
LastAccessed *time.Time `json:"last_accessed"`
}