package enums

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

var (

	// AppSecret is our app secret
	AppSecret string = os.Getenv("APP_SECRET")

	// InvitationCodeEnvironmentKey is the invitation code key that can
	// be used in a .env file
	InvitationCodeEnvironmentKey string = "INVITATION_CODE"

	// JWTAccessTokenExpiry is the global expiry time we'll use for JWT Acess tokens.
	JWTAccessTokenExpiry time.Time = time.Now().Add(24 * 7 * time.Hour) // one week

	// JWTSigningMethod is the global signing method we'll use for JWT tokens
	JWTSigningMethod jwt.SigningMethod = jwt.SigningMethodHS256


)
