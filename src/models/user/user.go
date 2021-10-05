package user

import (
	"aherman/src/models/base"
	tokenModels "aherman/src/models/token"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	// TokensFieldName should match the field name for user tokens
	TokensFieldName string = "Tokens"
)

// User represents a user of the application.
// Should never be made publicly available.
type User struct {
	base.Model
	Name string `json:"-"`
	Email string `json:"-" gorm:"unique_index:user_email_index"`
	Password string `json:"-" gorm:"size:72"`

	// foreign key will be Token.ID
	Tokens []tokenModels.Token `json:"-"`
}

// BeforeCreate is a GORM lifecycle method that will run before a user
// is created in the database.
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
  u.ID = uuid.New()
  return nil
}

// Creatable represents what's needed to successfully
// create a user.
type Creatable struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// Credentials represents what's needed to successfully
// log in to a user account.
type Credentials struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// JWTClaims is a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time.
type JWTClaims struct {
	ID uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
	SessionID string `json:"session_id"`
	TokenType string `json:"token_type"`
	jwt.StandardClaims
}

// Public represents the user object that is publicly
// available
type Public struct {
	ID uuid.UUID `json:"id"`
	Email string `json:"email"`
}

// BindAttributes TODO
func (up *Public) BindAttributes(u *User) {
	up.Email = u.Email
	up.ID = u.ID
}

// Readable represents the fields required to read a user from the database.
type Readable struct {
	ID uuid.UUID `json:"id" binding:"required"`
}

// Updateable represents the fields available for updating
type Updateable struct {
	Email    string `json:"email" binding:"email"`
	Password string `json:"password" binding:"min=8"`
}