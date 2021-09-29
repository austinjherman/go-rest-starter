package middleware

import (
	"aherman/src/enums"
	"aherman/src/http/response"
	userModels "aherman/src/models/user"
	"aherman/src/types/container"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type authorizationHeader struct {
	Authorization string `header:"Authorization"`
}

// Token is a middleware that requires a JWT web token to proceed
// with the request
func Token(app container.Container) gin.HandlerFunc {
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
			res := response.ErrUnauthorized
			c.AbortWithStatusJSON(res.Status, res)
			return
		}

		// validate token
		token, err := jwt.ParseWithClaims(bearerHeader[1], claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(enums.AppSecret), nil
		})
		if err != nil {
			c.Error(err)
			res := response.ErrUnauthorized
			c.AbortWithStatusJSON(res.Status, res)
			return
		}

		if !token.Valid {
			res := response.ErrUnauthorized
			c.AbortWithStatusJSON(res.Status, res)
			return
		}

		// find the user
		user, err := app.User.FindByID(claims.ID)
		if err != nil {
			c.Error(err)
			res := response.ErrUnauthorized
			c.AbortWithStatusJSON(res.Status, res)
			return
		}

		app.User.Current = user
		c.Next()
	}
}
