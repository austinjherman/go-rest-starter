package middleware

import (
	"aherman/src/container"
	"aherman/src/enums"
	"aherman/src/http/response"
	userModels "aherman/src/models/user"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type authorizationHeader struct {
	Authorization string `header:"Authorization"`
}

// Token is a middleware that requires a JWT web token to proceed
// with the request
func Token(app *container.Container) gin.HandlerFunc {
	return func(c *gin.Context) {

		var (
			bearer *authorizationHeader = &authorizationHeader{}
			claims *userModels.JWTClaims = &userModels.JWTClaims{}
		)

		// validate header
		err := c.ShouldBindHeader(bearer)
		if err != nil {
			c.Error(err)
			res := response.ErrUnknown
			c.AbortWithStatusJSON(res.Status, res)
			return
		}

		// split bearer header
		bearerHeader := strings.Split(bearer.Authorization, "Bearer ")
		if len(bearerHeader) != 2 {
			res := response.ErrAccessTokenInvalid
			c.AbortWithStatusJSON(res.Status, res)
			return
		}

		// validate token

		// for debugging
		// jwt.TimeFunc = func () time.Time {return time.Now().Add(24 * time.Hour)}
		
		token, err := jwt.ParseWithClaims(bearerHeader[1], claims, func(t *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				res := response.ErrAccessTokenInvalid
				res.Description = fmt.Sprintf("Unexpected signing method: %v", t.Header["alg"])
				return nil, res
			}
			return []byte(enums.AppSecret), nil
		})
		if err != nil {
			c.Error(err)
			res := response.ErrAccessTokenParse
			c.AbortWithStatusJSON(res.Status, res)
			return
		}

		// find the user
		user, err := app.User.FindByID(claims.UserID)
		if err != nil {
			c.Error(err)
			res := response.ErrAccessTokenInvalid
			c.AbortWithStatusJSON(res.Status, res)
			return
		}

		// explicitly check for expiry
		now := time.Now()
		tokenExpiry := time.Unix(claims.ExpiresAt, 0)
		if tokenExpiry.Before(now) {
			res := response.ErrAccessTokenExpired
			c.Error(res)
			c.AbortWithStatusJSON(res.Status, res)
			return
		}

		// explicitly check for refresh token type
		// we'll disallow on this route
		// also all other errors
		if !token.Valid || claims.TokenType != enums.JWTTokenTypeAccess {
			res := response.ErrAccessTokenInvalid
			c.AbortWithStatusJSON(res.Status, res)
			return
		}

		// lastly we'll check to see if the token is in the whitelist
		ok, err := app.Token.InWhitelist(claims.ID)
		if err != nil {
			res := response.ErrAccessTokenInvalid
			c.AbortWithStatusJSON(res.Status, res)
			return
		}
		if !ok {
			res := response.ErrAccessTokenInvalid
			c.AbortWithStatusJSON(res.Status, res)
			return
		}

		// store user in the handler dependencies
		app.Current.User = user

		// store the session id so we can revoke tokens related to the session
		app.Current.Token.SessionID = claims.SessionID

		c.Next()
	}
}
