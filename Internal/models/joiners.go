package models

import (
	"time"

	"github.com/google/uuid"
)

type SessionJoiners struct {
	ID        uuid.UUID     `gorm:"primaryKey"`
	SessionID uuid.UUID     `gorm:"not null"`                           // Foreign key for the session
	Session   StudioSession `gorm:"foreignKey:SessionID;references:ID"` // This defines the relationship
	UserID    uuid.UUID     `gorm:"not null"`                           // Foreign key for the user
	User      Users         `gorm:"foreignKey:UserID;references:ID"`    // This defines the relationship
	CreatedAt time.Time     `gorm:"autoCreateTime"`
	UpdatedAt time.Time     `gorm:"autoUpdateTime"`
}
