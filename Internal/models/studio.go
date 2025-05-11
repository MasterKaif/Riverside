package models

import (
	"time"

	"github.com/google/uuid"
)

type StudioSession struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	Name      string    `gorm:"not null"`
	HostID		uuid.UUID `gorm:"not null"` // Foreign key for the host
	Host      Users     `gorm:"foreignKey:HostID;references:ID"` // This defines the relationship
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
