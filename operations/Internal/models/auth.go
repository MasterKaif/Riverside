package models

import (
	"time"

	"github.com/google/uuid"
)

type Users struct {
	ID             uuid.UUID       `gorm:"primaryKey"`
	Username       string          `gorm:"not null"`
	Password       string          `gorm:"not null"`
	Email          string          `gorm:"unique;not null"`
	StudioSessions []StudioSession `gorm:"foreignKey:HostID"` // This defines the relationship
	CreatedAt      time.Time       `gorm:"autoCreateTime"`
	UpdatedAt      time.Time       `gorm:"autoUpdateTime"`
}
