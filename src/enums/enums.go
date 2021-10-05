package enums

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"

	// Import OS vars so we can use them in this file
	_ "github.com/joho/godotenv/autoload"
)

var (

	// AppSecret is our app secret
	AppSecret string = os.Getenv("APP_SECRET")

	// InvitationCode is the invitation code needed to use this application
	InvitationCode string = os.Getenv("INVITATION_CODE")

	// JWTAccessTokenExpiry is the global expiry time we'll use for JWT Access tokens.
	JWTAccessTokenExpiry time.Time = time.Now().Add(24 * time.Hour) // one day

	// JWTRefreshTokenExpiry is the global expiry time we'll use for JWT Refresh tokens.
	JWTRefreshTokenExpiry time.Time = time.Now().Add(24 * 7 * time.Hour * 30) // one month

	// JWTSigningMethod is the global signing method we'll use for JWT tokens
	JWTSigningMethod jwt.SigningMethod = jwt.SigningMethodHS256

	// JWTTokenTypeBearerOutgoing is the value that should be used for the token_type param
	JWTTokenTypeBearerOutgoing string = "bearer"

	// JWTTokenTypeAccess is the name to be used for access tokens
	JWTTokenTypeAccess string = "access"

	// JWTTokenTypeRefresh is the name to be used for refresh tokens
	JWTTokenTypeRefresh string = "refresh"

	// OAuthGrantTypePassword is the name used for the password grant type
	OAuthGrantTypePassword string = "password_grant"

	// RoleAdmin is the key used for the admin role
	RoleAdmin string = "admin"

	// RoleReader is the key used for the reader role
	RoleReader string = "reader"

	// RoleWriter is the key used for the writer role
	RoleWriter string = "writer"

	// TrustedClientID is the client id needed to use the password grant oauth flow.
	TrustedClientID string = os.Getenv("TRUSTED_CLIENT_ID")

	// TrustedClientSecret is the client secret needed to use the password grant oauth flow.
	TrustedClientSecret string = os.Getenv("TRUSTED_CLIENT_SECRET")

)
