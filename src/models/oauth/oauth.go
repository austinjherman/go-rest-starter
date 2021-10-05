package oauth

import (
	"time"

	"gorm.io/gorm"
)

// Env for dependency injection
type Env struct {
	DB *gorm.DB
}

// RouterRequest todo
type RouterRequest struct {
	GrantType string `json:"grant_type" binding:"required,oneof='password_grant'"`
}

// PasswordGrantRequest represents what's needed to authenticate a user via an email/password.
type PasswordGrantRequest struct {
	ClientID string `json:"client_id" binding:"required"`
	ClientSecret string `json:"client_secret" binding:"required"`
	Email string `json:"email" binding:"required"`
	GrantType string `json:"grant_type" binding:"required,oneof='password_grant'"`
	Password string `json:"password" binding:"required"`
	Scope string `json:"scope"`
}

// Success represents the response returned for a successful authentication
type Success struct {
	AccessToken string `json:"access_token"`
  ExpiresIn time.Time `json:"expires_in"`
  RefreshToken string `json:"refresh_token"`
  Scope string `json:"scope"`
  TokenType string `json:"token_type"`
}