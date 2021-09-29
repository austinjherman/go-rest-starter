package base

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Model TODO
type Model struct {
	ID        uuid.UUID      `json:"-" gorm:"primaryKey;type:uuid;not null"`
	CreatedAt time.Time      `json:"-" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"-" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}