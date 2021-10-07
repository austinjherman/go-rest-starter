package middleware

import (
	"aherman/src/container"
	"aherman/src/enums"
	"aherman/src/http/response"
	tokenModels "aherman/src/models/token"
	userModels "aherman/src/models/user"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/golang-jwt/jwt"
)

type authorizationHeader struct {
	Authorization string `header:"Authorization"`
}

// Token is a middleware that requires a valid JWT 
// to proceed with the request
func Token(app *container.Container) gin.HandlerFunc {
	return func(c *gin.Context) {

		var (
			bearer *authorizationHeader = &authorizationHeader{}
			claims *tokenModels.JWTClaims = &tokenModels.JWTClaims{}
			jwtToken *jwt.Token = &jwt.Token{}
			token *tokenModels.Token = &tokenModels.Token{}
			user *userModels.User = &userModels.User{}
		)

		// validate authorization header
		err := c.ShouldBindWith(bearer, binding.Header)
		ok, httpResponse := app.Facades.Error.ShouldContinue(err, &response.ErrValidation)
		if !ok {
			c.AbortWithStatusJSON(httpResponse.Status, httpResponse)
			return
		}

		// split authorization header value, i.e. from Bearer xxx to ["Bearer", "xxx"]
		bearerHeader := strings.Split(bearer.Authorization, " ")
		if len(bearerHeader) != 2 || bearerHeader[0] != "Bearer" {
			httpResponse := response.ErrTokenInvalid
			c.AbortWithStatusJSON(httpResponse.Status, httpResponse)
			return
		}

		// validate the token
		jwtToken, err = app.Facades.Token.ParseWithClaims(bearerHeader[1], claims)
		ok, httpResponse = app.Facades.Error.ShouldContinue(err, &response.ErrTokenInvalid)
		if !ok {
			c.AbortWithStatusJSON(httpResponse.Status, httpResponse)
			return
		}

		// find the user
		err = app.Facades.User.BindByID(user, claims.Subject)
		ok, httpResponse = app.Facades.Error.ShouldContinue(err, &response.ErrUserNotFound)
		if !ok {
			c.AbortWithStatusJSON(httpResponse.Status, httpResponse)
			return
		}

		// explicitly check for expiry for both refresh and access types
		// remove if expired
		now := time.Now()
		tokenExpiry := time.Unix(claims.ExpiresAt, 0)

		// if expiry is before now
		if tokenExpiry.Before(now) {
			// revoke the token
			app.Facades.User.RevokeTokenByID(claims.ID.String())

			// respond
			httpResponse := response.ErrTokenExpired
			// serverErr := logging.NewServerError(err).BindClientErr(httpResponse)
			// c.Error(serverErr)
			c.AbortWithStatusJSON(httpResponse.Status, httpResponse)
			return
		}

		// We'll only allow the access token flow in this middleware. If the user has
		// a refresh token, they should go through a different flow. For example, exchanging
		// their refresh token for a new access token.
		if claims.TokenType != enums.JWTTokenTypeAccess {
			httpResponse := response.ErrTokenTypeInvalid
			// serverErr := logging.NewServerError(nil).BindClientErr(httpResponse)
			// c.Error(serverErr)
			c.AbortWithStatusJSON(httpResponse.Status, httpResponse)
			return
		}

		// We'll disallow invalid tokens, too.
		if !jwtToken.Valid {
			fmt.Println("here it is")
			httpResponse := response.ErrTokenInvalid
			c.AbortWithStatusJSON(httpResponse.Status, httpResponse)
			return
		}

		// lastly we'll check to see if the token is in the whitelist
		err = app.Facades.Token.BindByID(token, claims.ID)
		if err != nil {
			httpResponse := response.ErrTokenNotFound
			// serverErr := logging.NewServerError(err).BindClientErr(httpResponse)
			// c.Error(serverErr)
			c.AbortWithStatusJSON(httpResponse.Status, httpResponse)
			return
		}

		// store user in the handler dependencies
		app.Current.User = user

		// store the token so we can revoke tokens related to the session
		app.Current.Token = token

		c.Next()
	}
}
