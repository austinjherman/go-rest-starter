package oauth

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Env for dependency injection
type Env struct {
	DB *gorm.DB
}

// PasswordGrantCredentials todo
type PasswordGrantCredentials struct {
	GrantType string `json:"grant_type" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	ClientID uuid.UUID `json:"client_id"`
	ClientSecret uuid.UUID `json:"client_secret"`
}