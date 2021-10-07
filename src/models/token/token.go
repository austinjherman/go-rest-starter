package token

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Token todo
type Token struct {

	// ID is the primary token identifier
	ID uuid.UUID `gorm:"primaryKey;type:uuid;not null"`

	// UserID is the foreign key that ties a token to a user. This association allows
	// us to create a token whitelist, which subsequently allows us to revoke all tokens
	// or revoke tokens on behalf of a user.
	UserID uuid.UUID `gorm:"type:uuid;not null"`

	SessionID string `gorm:"not null"`

	// Token is the string representation of the JWT token.
	Token string

	// Type, e.g. "access" or "refresh"
	// Acceptable type values can be found in enums.
	Type string
}

// BeforeCreate is a GORM lifecycle method that will run before a token
// is created in the database.
func (t *Token) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == uuid.Nil {
  	t.ID = uuid.New()
	}
  return nil
}

// TableName provides GORM with a customized table name.
func (Token) TableName() string {
	return "token_whitelist"
}


// JWTClaims is a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time.
type JWTClaims struct {
	ID uuid.UUID `json:"id"`
	SessionID string `json:"session_id"`
	TokenType string `json:"token_type"`
	jwt.StandardClaims
}
