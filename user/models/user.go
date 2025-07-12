package models

import (
	"time"
)

type User struct {
	ID           string    `gorm:"primaryKey;autoIncrement"`
	Name         string    `gorm:"type:text;not null"`
	Email        string    `gorm:"uniqueIndex;not null"`
	PasswordHash string    `gorm:"not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}
